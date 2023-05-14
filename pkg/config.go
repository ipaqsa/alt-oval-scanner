package pkg

import "strings"

type ConfigT struct {
	BranchesUrl  string
	BaseUrl      string
	PathToDB     string
	PathToTmpDir string
}

var Config = ConfigT{}

func (c *ConfigT) Valid() bool {
	if c.BaseUrl == "" || c.BranchesUrl == "" || c.PathToDB == "" || c.PathToTmpDir == "" {
		return false
	}
	splits := strings.Split(c.PathToDB, ".")
	if len(splits) < 2 {
		return false
	}
	if splits[len(splits)-1] != "db" {
		return false
	}
	return true
}
