package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/jsonapi"
)

type Credentials struct {
	Id           string     `jsonapi:"primary,upload_credential"`
	KeyPrefix    string     `jsonapi:"attr,key_prefix"`
	Bucket       string     `jsonapi:"attr,bucket"`
	AccessKey    string     `jsonapi:"attr,access_key"`
	SecretKey    string     `jsonapi:"attr,secret_key"`
	SessionToken string     `jsonapi:"attr,session_token"`
	Expiration   *time.Time `jsonapi:"attr,expiration,iso8601"`
	Region       string     `jsonapi:"attr,region"`
	Session      *awssession.Session
}

type Upload struct {
	Id        string `jsonapi:"primary,upload"`
	UploadKey string `jsonapi:"attr,upload_key"`
	Status    string `jsonapi:"attr,status,omitempty"`
	Message   string `jsonapi:"attr,message,omitempty"`
	FileId    int    `jsonapi:"attr,file_id,omitempty"`
}

func (credentials *Credentials) Upload(directory string, filePath string, settings Settings) (error, int) {
	if err := credentials.checkExpired(settings); err != nil {
		log.Fatalf("Failure getting credentials: %s", err)
	}
	log.Printf("Uploading file %s", filePath)

	err, upload := credentials.UploadFile(fmt.Sprintf("%s/%s", directory, filePath))
	if err != nil {
		log.Fatalf("Failure uploading file: %s", err)
	}

	if settings.Debug {
		log.Printf("Processing file %s", filePath)
	}
	if err = processUpload(settings, &upload); err != nil {
		log.Fatalf("Failure processing upload: %s", err)
	}

	if settings.Debug {
		log.Printf("Polling process file %s: %s", filePath, upload.Id)
	}
	start := time.Now()
	sleep := 1
	for {
		upload.Poll(settings)
		t := time.Now()
		if (upload.Status != "queued" && upload.Status != "processing") || int(t.Sub(start).Seconds()) > settings.UploadTimeout {
			break
		} else {
			time.Sleep(time.Duration(sleep))
			sleep = sleep + 1
		}
	}

	if upload.Status != "success" {
		log.Fatalf("Upload process failed for %s: %s", filePath, upload.Status)
	}

	return err, upload.FileId
}

func (credentials *Credentials) UploadFile(source string) (error, Upload) {
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(credentials.Session)
	var upload Upload

	f, err := os.Open(source)
	if err != nil {
		return err, upload
	}
	defer f.Close()

	_, filename := filepath.Split(source)
	upload.UploadKey = fmt.Sprintf("%s%s", credentials.KeyPrefix, filename)

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(credentials.Bucket),
		Key:    aws.String(upload.UploadKey),
		Body:   f,
	})
	return err, upload
}

func (credentials *Credentials) checkExpired(settings Settings) error {
	var err error
	if credentials.Expiration == nil || withinSeconds(credentials.Expiration, 15) {
		err = credentials.updateCredentials(settings)
	}
	return err
}

func (credentials *Credentials) updateCredentials(settings Settings) error {
	if settings.Debug {
		log.Printf("Update credentials")
	}
	url := fmt.Sprintf("%s/api/uploads/credentials", settings.Server)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal("Error building request for Upload Credentials", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.Token))
	req.Header.Add("Accept", "application/vnd.api+json; com.turbosquid.api.version=1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing request for Upload Credentials", err)
		return err
	}

	defer resp.Body.Close()

	if err = jsonapi.UnmarshalPayload(resp.Body, credentials); err != nil {
		return err
	}

	credentials.Session, err = awssession.NewSession(&aws.Config{
		Region:      aws.String(credentials.Region),
		Credentials: awscreds.NewStaticCredentials(credentials.AccessKey, credentials.SecretKey, credentials.SessionToken),
	})
	return err
}

func processUpload(settings Settings, upload *Upload) error {
	url := fmt.Sprintf("%s/api/uploads", settings.Server)
	var message bytes.Buffer
	if err := jsonapi.MarshalPayload(&message, upload); err != nil {
		log.Fatal("Error building upload process message", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(message.Bytes()))
	if err != nil {
		log.Fatal("Error building request for Upload Process", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.Token))
	req.Header.Add("Accept", "application/vnd.api+json; com.turbosquid.api.version=1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing request for Upload Process", err)
		return err
	}

	defer resp.Body.Close()

	if err = jsonapi.UnmarshalPayload(resp.Body, upload); err != nil {
		return err
	}
	if settings.Debug {
		log.Printf("Upload: %s", upload.Id)
	}

	return err
}

func (upload *Upload) Poll(settings Settings) error {
	url := fmt.Sprintf("%s/api/uploads/%s", settings.Server, upload.Id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error building request for Upload Poll", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.Token))
	req.Header.Add("Accept", "application/vnd.api+json; com.turbosquid.api.version=1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing request for Upload Poll", err)
		return err
	}
	defer resp.Body.Close()

	err = jsonapi.UnmarshalPayload(resp.Body, upload)

	return err
}

func withinSeconds(expiration *time.Time, seconds int) bool {
	t := time.Now()
	diff := expiration.Sub(t)
	return int(diff.Seconds()) < seconds
}
