package scanner

import (
	"alt-oval-scanner/pkg/oval"
	"alt-oval-scanner/pkg/rpm"
	"alt-oval-scanner/pkg/utils"
	"log"
)

func NewScanner(branchesUrl, baseUrl, pathToDB, tmpDir, outputFile string, format OutputFormat) (*Scanner, error) {
	var scanner = Scanner{}
	var err error
	scanner.manager, err = oval.NewOvalManager(branchesUrl, baseUrl, pathToDB)
	if err != nil {
		return nil, err
	}
	scanner.tpmDir = tmpDir
	scanner.output = format
	scanner.outputFile = outputFile
	err = scanner.manager.Download()
	if err != nil {
		return nil, err
	}
	return &scanner, nil
}

func (s *Scanner) ScanImage(image string) error {
	is := &imageScanner{}
	is.image = image
	is.tmpDir = s.tpmDir
	packages, release, err := is.scan()
	if err != nil {
		return err
	}
	vulns, err := s.check(packages, release)
	if err != nil {
		return err
	}
	s.Output(vulns)
	return err
}
func (s *Scanner) ScanHost() error {
	hs := &hostScanner{}
	hs.tmpDir = s.tpmDir
	packages, release, err := hs.scan()
	if err != nil {
		return err
	}
	vulns, err := s.check(packages, release)
	if err != nil {
		return err
	}
	s.Output(vulns)
	return err
}

func (s *Scanner) check(packages []rpm.Package, release utils.OsRelease) (map[string][]oval.Vuln, error) {
	ovals, err := s.manager.OVALs()
	if err != nil {
		return nil, err
	}
	log.Printf("start check....")
	var vulns = map[string][]oval.Vuln{}
	for _, oval := range ovals {
		err = s.manager.OvalCheck(vulns, oval, packages, release)
		if err != nil {
			return nil, err
		}
	}
	return vulns, nil
}
