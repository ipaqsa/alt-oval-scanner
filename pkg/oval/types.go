package oval

import "alt-oval-scanner/pkg/repository"

type Manager struct {
	branchesUrl string
	baseUrl     string
	repository  *repository.Repository
}

type Branches struct {
	Length   int      `json:"length"`
	Branches []string `json:"branches"`
}

type AdjTest struct {
	ID      string        `xml:"id,attr" json:",omitempty"`
	Version string        `xml:"version,attr" json:",omitempty"`
	Check   string        `xml:"check,attr" json:",omitempty"`
	Comment string        `xml:"comment,attr" json:",omitempty"`
	Object  RpmInfoObject `xml:"object" json:",omitempty"`
	State   RpmInfoState  `xml:"state" json:",omitempty"`
}

type AdjTestBase struct {
	ID      string                  `xml:"id,attr" json:",omitempty"`
	Version string                  `xml:"version,attr" json:",omitempty"`
	Check   string                  `xml:"check,attr" json:",omitempty"`
	Comment string                  `xml:"comment,attr" json:",omitempty"`
	Object  TextFileContent54Object `xml:"object" json:",omitempty"`
	State   TextFileContent54State  `xml:"state" json:",omitempty"`
}

type Vuln struct {
	Title            string
	Version          string
	InstalledVersion string
	References       []Reference
	CVEs             []CVE
	Severity         string
}
