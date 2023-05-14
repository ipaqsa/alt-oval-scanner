package scanner

import (
	"alt-oval-scanner/pkg/rpm"
	"alt-oval-scanner/pkg/utils"
	"archive/tar"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	"io"
	"log"
)

func (s *imageScanner) scan() ([]rpm.Package, utils.OsRelease, error) {
	if s.image == "" {
		return nil, utils.OsRelease{}, errors.New("image empty")
	}
	err := s.downloadImage()
	if err != nil {
		return nil, utils.OsRelease{}, err
	}
	err = s.layersPreprocess()
	if err != nil {
		return nil, utils.OsRelease{}, err
	}
	return s.packages, s.osRelease, err
}

func (s *imageScanner) downloadImage() error {
	log.Printf("download %s.....", s.image)
	ref, err := name.ParseReference(s.image)
	if err != nil {
		return err
	}
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}
	layers, err := img.Layers()
	if err != nil {
		return err
	}
	s.layers = layers
	return nil
}

func (s *imageScanner) layersPreprocess() error {
	var release = false
	var packages = false
	for _, l := range s.layers {
		if release && packages {
			return nil
		}
		err := s.walk(l, &release, &packages)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *imageScanner) walk(layer v1.Layer, release, packages *bool) error {
	layerReader, err := layer.Uncompressed()
	if err != nil {
		return err
	}
	tr := tar.NewReader(layerReader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		switch hdr.FileInfo().Name() {
		case "os-release":
			if *release {
				continue
			}
			*release = true
			if err = s.parseOSRelease(tr); err != nil {
				return err
			}
			s.osRelease.Print()
			if !s.osRelease.Validate() {
				log.Fatal("Use ALT Linux Image")
			}
			continue
		case "Packages":
			if *packages {
				continue
			}
			*packages = true
			p, i, err := rpm.ParsePkgInfo(s.tmpDir, tr)
			if err != nil {
				return err
			}
			log.Printf("find: %d packages", len(p))
			s.packages = p
			s.installedPackages = i
			continue
		}
	}
	return nil
}

func (s *imageScanner) parseOSRelease(reader *tar.Reader) error {
	file, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	s.osRelease = utils.ReleaseParse(string(file))
	return nil
}
