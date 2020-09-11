package main

import (
	"bufio"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Settings struct {
	Token         string `yaml:"token"`
	Server        string `yaml:"server,omitempty"`
	Debug         bool   `yaml:"debug,omitempty"`
	UploadTimeout int    `yaml:"upload_timeout,omitempty"`
}

func GetSettings() Settings {
	var s Settings
	yamlFile, err := ioutil.ReadFile("settings.yml")

	if err != nil {
		reader := bufio.NewReader(os.Stdin)
		println("Settings file does not exist")
		println("Please enter your TurboSquid API Key:")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) > 0 {
			s.Token = text
			yamlFile, err := yaml.Marshal(s)
			if err != nil {
				log.Fatal("Error creating settings yml text", err)
			}
			if err = ioutil.WriteFile("settings.yml", yamlFile, 0644); err != nil {
				log.Fatal("Error writing setting.yml file", err)
			}
			println("API Key saved to settings.yml")
		} else {
			log.Fatal("Invalid api key entered")
		}
	}

	err = yaml.Unmarshal(yamlFile, &s)
	if err != nil {
		log.Fatalf("settings.yml is not properly formatted", err)
	}
	if s.Server == "" {
		s.Server = "https://api.turbosquid.com"
	}
	if s.UploadTimeout == 0 {
		s.UploadTimeout = 90
	}

	if s.Token == "" {
		log.Fatalf("settings.yml must contain a valid API Token")
	}

	return s
}
