package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const VERSION = "1.0.0"

func main() {
	settings := GetSettings()

	if len(os.Args) != 2 {
		println("")
		println("ts-publishing-api-go version", VERSION)
		println("Please run as 'ts-publishing-api <path to product folder>'")
		println("See project README.md for more information.")
		println("")
		os.Exit(0)
	}

	directory := os.Args[1]
	productBundle := ReadInput(directory)
	var credentials Credentials

	if err := productBundle.Draft.createDraft(settings); err != nil {
		log.Fatal("Error creating draft: ", err)
	}

	for _, file := range productBundle.Files {
		err, fileId := credentials.Upload(directory, file.Name, settings)
		if err != nil {
			log.Fatal("Error uploading file: ", err)
		}
		file.FileId = fileId

		productBundle.Draft.addFile(file, settings)
	}

	for _, preview := range productBundle.Previews {
		if preview.Type == "thumbnail" {
			err, fileId := credentials.Upload(directory, preview.Name, settings)
			if err != nil {
				log.Fatal("Error uploading preview: ", err)
			}
			preview.FileId = fileId

			productBundle.Draft.addThumbnail(preview, settings)
		} else if preview.Type == "turntable" {
			files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", directory, preview.Name))
			if err != nil {
				log.Fatal("Error reading turntable directory: ", err)
			}

			for _, file := range files {
				if strings.HasPrefix(file.Name(), ".") {
					continue
				}
				err, fileId := credentials.Upload(directory, fmt.Sprintf("%s/%s", preview.Name, file.Name()), settings)
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

	err, productId := productBundle.Draft.publish(settings)
	if err != nil {
		log.Fatal("Error publishing product: ", err)
	}

	log.Printf("Successfully published product ID: %d", productId)
}
