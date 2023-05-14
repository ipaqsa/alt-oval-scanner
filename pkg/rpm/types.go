package rpm

import (
	"github.com/opencontainers/go-digest"
	"io"
)

type RPMDB struct {
	io.Reader
}

type Package struct {
	ID         string   `json:",omitempty"`
	Name       string   `json:",omitempty"`
	Version    string   `json:",omitempty"`
	Release    string   `json:",omitempty"`
	Epoch      int      `json:",omitempty"`
	Arch       string   `json:",omitempty"`
	SrcName    string   `json:",omitempty"`
	SrcVersion string   `json:",omitempty"`
	SrcRelease string   `json:",omitempty"`
	SrcEpoch   int      `json:",omitempty"`
	Licenses   []string `json:",omitempty"`
	Maintainer string   `json:",omitempty"`

	BuildInfo *BuildInfo `json:",omitempty"` // only for Red Hat

	DependsOn []string `json:",omitempty"`

	FilePath string `json:",omitempty"`

	Digest digest.Digest `json:",omitempty"`

	Locations []Location `json:",omitempty"`
}

type BuildInfo struct {
	ContentSets []string `json:",omitempty"`
	Nvr         string   `json:",omitempty"`
	Arch        string   `json:",omitempty"`
}

type Location struct {
	StartLine int `json:",omitempty"`
	EndLine   int `json:",omitempty"`
}

type Version struct {
	epoch   int
	version string
	release string
}
