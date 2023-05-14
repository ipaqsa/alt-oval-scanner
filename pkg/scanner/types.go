package scanner

import (
	"alt-oval-scanner/pkg/oval"
	"alt-oval-scanner/pkg/rpm"
	"alt-oval-scanner/pkg/utils"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

type OutputFormat string

type Scanner struct {
	manager    *oval.Manager
	tpmDir     string
	output     OutputFormat
	outputFile string
}

type imageScanner struct {
	tmpDir            string
	image             string
	layers            []v1.Layer
	osRelease         utils.OsRelease
	installedPackages []string
	packages          []rpm.Package
}

type hostScanner struct {
	tmpDir string
}
