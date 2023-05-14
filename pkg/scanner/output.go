package scanner

import (
	"alt-oval-scanner/pkg/oval"
	"alt-oval-scanner/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

func (s *Scanner) Output(vulns map[string][]oval.Vuln) {
	if s.output == PrintFormat {
		s.print(vulns)
	}
	if s.output == JsonFormat {
		s.json(vulns)
	}
}

func (s *Scanner) print(vulns map[string][]oval.Vuln) {
	count := 0
	for _, vul := range vulns {
		count += len(vul)
	}
	log.Printf("Vulns: %s", utils.RedColor(strconv.Itoa(count)))
	for v, k := range vulns {
		fmt.Printf("Package: %s (v. %s):", utils.BlueColor(v), utils.RedColor(k[0].InstalledVersion))
		fmt.Printf("%s", utils.BlueColor("\n [--------------------------\n"))
		for _, v := range k {
			for _, c := range v.CVEs {
				fmt.Printf("  %s - %s - %s - %s\n", utils.BlueColor(c.CveID), impactColor(c.Impact), utils.GreenColor(c.Cvss3), c.Href)
			}
			for _, r := range v.References {
				if r.Source == "CVE" {
					continue
				}
				fmt.Printf("  %s - %s - %s\n", utils.BlueColor(r.RefID), impactColor(v.Severity), r.RefURL)
			}
		}
		fmt.Printf("%s", utils.BlueColor(" --------------------------]\n"))
	}
}
func impactColor(val string) string {
	switch val {
	case "High":
		return utils.RedColor(val)
	case "Low":
		return utils.GreenColor(val)
	default:
		return utils.BlueColor(val)
	}
}

func (s *Scanner) json(vulns map[string][]oval.Vuln) {
	bytes, err := json.Marshal(vulns)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(s.outputFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
