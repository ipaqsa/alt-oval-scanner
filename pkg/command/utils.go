package command

import (
	"alt-oval-scanner/pkg"
	"alt-oval-scanner/pkg/scanner"
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"strings"
)

func pathToConfIsValid(path string) error {
	if path == "" {
		return errors.New("path to config is empty")
	}
	splits := strings.Split(path, ".")
	if len(splits) < 2 {
		return errors.New("path to config is invalid")
	}
	format := splits[len(splits)-1]
	if format != "yml" {
		return errors.New("config format is not yml")
	}
	return nil
}

func splitPath(cpath string) (string, string, string) {
	ext := path.Ext(cpath)
	splits := strings.Split(cpath, "/")
	filename := splits[len(splits)-1]
	cpath = strings.TrimSuffix(cpath, filename)
	filename = strings.TrimSuffix(filename, ext)
	return cpath, filename, ext
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	log.Println(err.Error())
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func initConfiguration() {
	err := pathToConfIsValid(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	initInfo()
	path, filename, ext := splitPath(*configPath)
	if !exists(*configPath) {
		log.Fatal("path to config is invalid")
	}
	if path != "" {
		viper.AddConfigPath(path)
	} else {
		viper.AddConfigPath(".")
	}
	viper.SetConfigName(filename)
	viper.SetConfigType(ext[1:])
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&pkg.Config)
	if err != nil {
		log.Fatal(err)
	}
	if !pkg.Config.Valid() {
		log.Fatal("config is invalid")
	}
}

func toOutputFormat(data string) (string, scanner.OutputFormat, error) {
	if strings.Contains(data, "json") {
		path, _, ext := splitPath(data)
		if ext != ".json" {
			return "", "", errors.New("format is not json")
		}
		if path != "" && !exists(path) {
			return "", "", errors.New("path to output file is invalid")
		}
		return data, scanner.JsonFormat, nil
	}
	return "", scanner.PrintFormat, nil
}
