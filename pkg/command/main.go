package command

import (
	"alt-oval-scanner/pkg"
	"alt-oval-scanner/pkg/scanner"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	setVersion("0.0.1")

	RootCmd.AddCommand(imageCmd)
	RootCmd.AddCommand(fsCmd)

	configPath = RootCmd.PersistentFlags().StringP("config", "c", "", "path to config")
	outputFormat = RootCmd.PersistentFlags().StringP("output", "o", "print", "path to output json file")

	cobra.OnInitialize(initConfiguration)
}

var RootCmd = &cobra.Command{
	Use: "alt-oval-scanner",
}

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Scan image",
	Run:   image,
}

var fsCmd = &cobra.Command{
	Use:   "host",
	Short: "Scan host",
	Run:   host,
}

func image(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("bad arguments, use: image registry.altlinux.org/alt/alt")
	}
	file, format, err := toOutputFormat(*outputFormat)
	if err != nil {
		log.Fatal(err)
	}
	sc, err := scanner.NewScanner(pkg.Config.BranchesUrl, pkg.Config.BaseUrl, pkg.Config.PathToDB,
		pkg.Config.PathToTmpDir, file, format)
	if err != nil {
		log.Fatal(err)
	}
	err = sc.ScanImage(args[0])
	if err != nil {
		log.Fatal(err)
	}
}

func host(cmd *cobra.Command, args []string) {
	file, format, err := toOutputFormat(*outputFormat)
	if err != nil {
		log.Fatal(err)
	}
	sc, err := scanner.NewScanner(pkg.Config.BranchesUrl, pkg.Config.BaseUrl, pkg.Config.PathToDB,
		pkg.Config.PathToTmpDir, file, format)
	if err != nil {
		log.Fatal(err)
	}
	err = sc.ScanHost()
	if err != nil {
		log.Fatal(err)
	}
}
