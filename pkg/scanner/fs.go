package scanner

import (
	"alt-oval-scanner/pkg/rpm"
	"alt-oval-scanner/pkg/utils"
	"io"
	"log"
	"os"
)

func (h *hostScanner) scan() ([]rpm.Package, utils.OsRelease, error) {
	release, err := parseRelease()
	if err != nil {
		return nil, utils.OsRelease{}, nil
	}
	release.Print()
	if !release.Validate() {
		log.Fatal("scanner support only ALT Linux")
	}
	packages, err := parsePackages(h.tmpDir)
	if err != nil {
		return nil, utils.OsRelease{}, nil
	}
	log.Printf("find: %d packages", len(packages))
	return packages, release, err
}

func parseRelease() (utils.OsRelease, error) {
	f, err := os.OpenFile("/etc/os-release", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return utils.OsRelease{}, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return utils.OsRelease{}, err
	}
	r := utils.ReleaseParse(string(data))
	return r, nil
}

func parsePackages(tmpDir string) ([]rpm.Package, error) {
	f, err := os.OpenFile("/var/lib/rpm/Packages", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	p, _, err := rpm.ParsePkgInfo(tmpDir, f)
	if err != nil {
		return nil, err
	}
	return p, nil
}
