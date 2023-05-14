package utils

import (
	"log"
	"regexp"
	"strings"
)

type OsRelease struct {
	Name      string
	Version   string
	Id        string
	VersionID string
	CPE       string
	HomeUrl   string
	Data      string
}

func ReleaseParse(data string) OsRelease {
	var release OsRelease
	release.Data = data
	splits := strings.Split(data, "\n")
	for _, s := range splits {
		parts := strings.Split(s, "=")
		switch parts[0] {
		case "NAME":
			release.Name = parts[1]
			continue
		case "VERSION":
			release.Version = parts[1]
			continue
		case "ID":
			release.Id = parts[1]
			continue
		case "VERSION_ID":
			release.VersionID = parts[1]
			continue
		case "CPE_NAME":
			release.CPE = parts[1]
			continue
		case "HOME_URL":
			release.HomeUrl = parts[1]
			continue
		}
	}
	return release
}

func (r *OsRelease) Check(exp *regexp.Regexp) bool {
	return exp.MatchString(r.Data)
}

func (r *OsRelease) Print() {
	log.Printf("\n--------------------------------------\n"+
		"%s\nname: %s\nversion: %s\nversion_id: %s\ncpe: %s\nhome url: %s\n--------------------------------------\n",
		BlueColor("OS RELEASE:"), BlueColor(r.Name), BlueColor(r.Version), BlueColor(r.VersionID), BlueColor(r.CPE), BlueColor(r.HomeUrl))
}

func (r *OsRelease) Validate() bool {
	return strings.Contains(strings.ToLower(r.Data), "alt")
}
