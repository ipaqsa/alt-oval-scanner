package command

var version string

func getVersion() string {
	return version
}

func setVersion(v string) {
	version = v
}
