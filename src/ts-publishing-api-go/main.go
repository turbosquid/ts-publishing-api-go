package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const VERSION = "1.2.0"

type Params struct {
	Path    string
	Publish bool
}

func ParseParams() Params {
	var params Params
	binName := filepath.Base(os.Args[0])
	flag.StringVar(&params.Path, "path", "", "Path to product folder")
	flag.BoolVar(&params.Publish, "publish", false, "Publish draft after creation.")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\n%s %s:\n", binName, VERSION)
		fmt.Fprintf(flag.CommandLine.Output(), "See project README.md for more information.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nArguments:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExample:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\t%s -path <project path>\n", binName)
	}
	flag.Parse()

	if params.Path == "" {
		flag.Usage()
		os.Exit(0)
	}

	return params
}

func main() {
	settings := GetSettings()

	params := ParseParams()

	productBundle := ReadInput(params.Path)
	var credentials Credentials

	if err := productBundle.Draft.createDraft(settings); err != nil {
		log.Fatal("Error creating draft: ", err)
	}

	for _, file := range productBundle.Files {
		log.Printf("Uploading file: %s", file.Name)
		err, fileId := credentials.Upload(productBundle.Directory, file.Name, settings)
		if err != nil {
			log.Fatal("Error uploading file: ", err)
		}
		file.FileId = fileId

		productBundle.Draft.addFile(file, settings)
	}

	for _, preview := range productBundle.Previews {
		if preview.Type == "thumbnail" {
			err, fileId := credentials.Upload(productBundle.Directory, preview.Name, settings)
			if err != nil {
				log.Fatal("Error uploading preview: ", err)
			}
			preview.FileId = fileId

			productBundle.Draft.addThumbnail(preview, settings)
		} else if preview.Type == "turntable" {
			files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", productBundle.Directory, preview.Name))
			if err != nil {
				log.Fatal("Error reading turntable directory: ", err)
			}

			for _, file := range files {
				if strings.HasPrefix(file.Name(), ".") {
					continue
				}
				err, fileId := credentials.Upload(productBundle.Directory, fmt.Sprintf("%s/%s", preview.Name, file.Name()), settings)
				if err != nil {
					log.Fatal("Error uploading turntable file: ", err)
				}
				preview.FileIds = append(preview.FileIds, fileId)
			}

			productBundle.Draft.addTurntable(preview, settings)
		}
	}

	if err := productBundle.Draft.certifications(settings, productBundle.Certifications); err != nil {
		log.Fatal("Error setting certifications: ", err)
	}

	if params.Publish {
		err, productId := productBundle.Draft.publish(settings)
		if err != nil {
			log.Fatal("Error publishing product: ", err)
		}
		log.Printf("Successfully published product ID: %d", productId)
	}
}
