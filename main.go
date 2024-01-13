package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/charmbracelet/huh"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const CTX_FILE_PATH = ".awsctx"

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	var files = config.DefaultSharedConfigFiles
	//var awsConfigData *ini.File
	var formOptions []huh.Option[string]
	for _, f := range files {
		d, err := ini.Load(f)
		if err != nil {
			fmt.Printf("Fail to read file: %v", err)
			os.Exit(1)
		}
		for _, s := range d.Sections() {
			if strings.HasPrefix(s.Name(), "profile") {
				formOptions = append(formOptions, huh.Option[string]{
					Key:   strings.Split(s.Name(), " ")[1],
					Value: strings.Split(s.Name(), " ")[1],
				})
			}
		}
		if len(formOptions) > 0 {
			break
		} else {
			log.Fatal("no profiles found")
		}
	}

	var profileToSet string
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&profileToSet).
				Options(
					formOptions...,
				).
				Title("Choose a profile"),
		),
	)

	err := f.Run()
	if err != nil {
		log.Fatal(err)
	}

	ctxFile := fmt.Sprintf("export AWS_PROFILE=%s", profileToSet)

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(homedir+"/"+CTX_FILE_PATH, []byte(ctxFile), 0644); err != nil {
		log.Fatal(err)
	}
}
