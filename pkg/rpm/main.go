package rpm

import (
	"errors"
	"fmt"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"golang.org/x/xerrors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ParsePkgInfo(tmpPath string, rc io.Reader) ([]Package, []string, error) {
	filePath, err := writeToTempFile(tmpPath, rc)
	if err != nil {
		return nil, nil, xerrors.Errorf("temp file error: %w", err)
	}
	defer os.RemoveAll(filepath.Dir(filePath))

	db, err := rpmdb.Open(filePath)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to open RPM DB: %w", err)
	}

	pkgList, err := db.ListPackages()
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to list packages: %w", err)
	}

	var pkgs []Package
	var installedFiles []string
	provides := map[string]string{}
	for _, pkg := range pkgList {
		arch := pkg.Arch
		if arch == "" {
			arch = "None"
		}

		var srcName, srcVer, srcRel string
		if pkg.SourceRpm != "(none)" && pkg.SourceRpm != "" {
			srcName, srcVer, srcRel, err = splitFileName(pkg.SourceRpm)
			if err != nil {
				log.Printf("Invalid Source RPM Found: %s", pkg.SourceRpm)
			}
		}

		var files []string
		files, err = pkg.InstalledFileNames()
		if err != nil {
			return nil, nil, xerrors.Errorf("unable to get installed files: %w", err)
		}

		p := Package{
			ID:         fmt.Sprintf("%s@%s-%s.%s", pkg.Name, pkg.Version, pkg.Release, pkg.Arch),
			Name:       pkg.Name,
			Epoch:      pkg.EpochNum(),
			Version:    pkg.Version,
			Release:    pkg.Release,
			Arch:       arch,
			SrcName:    srcName,
			SrcEpoch:   pkg.EpochNum(),
			SrcVersion: srcVer,
			SrcRelease: srcRel,
			Licenses:   []string{pkg.License},
			DependsOn:  pkg.Requires,
			Maintainer: pkg.Vendor,
		}
		pkgs = append(pkgs, p)
		installedFiles = append(installedFiles, files...)

		for _, provide := range pkg.Provides {
			provides[provide] = p.ID
		}
	}

	consolidateDependencies(pkgs, provides)

	return pkgs, installedFiles, nil
}

func splitFileName(filename string) (name, ver, rel string, err error) {
	filename = strings.TrimSuffix(filename, ".rpm")

	archIndex := strings.LastIndex(filename, ".")
	if archIndex == -1 {
		return "", "", "", errors.New("bad name format")
	}

	relIndex := strings.LastIndex(filename[:archIndex], "-")
	if relIndex == -1 {
		return "", "", "", errors.New("bad name format")
	}
	rel = filename[relIndex+1 : archIndex]

	verIndex := strings.LastIndex(filename[:relIndex], "-")
	if verIndex == -1 {
		return "", "", "", errors.New("bad name format")
	}
	ver = filename[verIndex+1 : relIndex]

	name = filename[:verIndex]
	return name, ver, rel, nil
}

func writeToTempFile(tmpPath string, rc io.Reader) (string, error) {
	tmpDir, err := os.MkdirTemp(tmpPath, "rpm")
	if err != nil {
		return "", xerrors.Errorf("failed to create a temp dir: %w", err)
	}

	filePath := filepath.Join(tmpDir, "Packages")
	f, err := os.Create(filePath)
	if err != nil {
		return "", xerrors.Errorf("failed to create a package file: %w", err)
	}

	if _, err = io.Copy(f, rc); err != nil {
		return "", xerrors.Errorf("failed to copy a package file: %w", err)
	}

	if err = f.Close(); err != nil {
		return "", xerrors.Errorf("failed to close a temp file: %w", err)
	}

	return filePath, nil
}

func consolidateDependencies(pkgs []Package, provides map[string]string) {
	for i := range pkgs {
		pkgs[i].DependsOn = lo.FilterMap(pkgs[i].DependsOn, func(d string, _ int) (string, bool) {
			if pkgID, ok := provides[d]; ok && pkgs[i].ID != pkgID {
				return pkgID, true
			}
			return "", false
		})
		sort.Strings(pkgs[i].DependsOn)
		pkgs[i].DependsOn = slices.Compact(pkgs[i].DependsOn)

		if len(pkgs[i].DependsOn) == 0 {
			pkgs[i].DependsOn = nil
		}
	}
}
