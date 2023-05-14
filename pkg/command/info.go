package command

import (
	"fmt"
	"os"
	"runtime"
)

var Info InfoT

func initInfo() {
	Info.OS = runtime.GOOS
	Info.Arch = runtime.GOARCH
	if _, err := os.Lstat("/.dockerenv"); err != nil && os.IsNotExist(err) {
		Info.Container = "outside"
	} else {
		Info.Container = "inside"
	}
	Info.Hostname, _ = os.Hostname()
	printInfo()
}

func printInfo() {
	fmt.Printf("v%s\nHostname > %s\nDocker container > %s\nOS > %s/%s\nPath to config > %s\n",
		getVersion(), Info.Hostname,
		Info.Container, Info.OS, Info.Arch, *configPath)
}
