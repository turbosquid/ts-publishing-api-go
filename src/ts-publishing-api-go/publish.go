package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/google/jsonapi"
)

type Thumbnail struct {
	Id     int    `jsonapi:"primary,thumbnail"`
	FileId int    `jsonapi:"attr,file_id"`
	Type   string `jsonapi:"attr,thumbnail_type,omitempty"`
}
type Turntable struct {
	Id      int    `jsonapi:"primary,turntable"`
	FileIds []int  `jsonapi:"attr,file_ids"`
	Type    string `jsonapi:"attr,thumbnail_type,omitempty"`
}
type ProductFile struct {
	Id              int    `jsonapi:"primary,product_file"`
	FileId          int    `jsonapi:"attr,file_id"`
	Format          string `jsonapi:"attr,file_format"`
	FormatVersion   string `jsonapi:"attr,format_version,omitempty"`
	Renderer        string `jsonapi:"attr,renderer,omitempty"`
	RendererVersion string `jsonapi:"attr,renderer_version,omitempty"`
	Native          bool   `jsonapi:"attr,is_native,omitempty"`
}
type CustomerFile struct {
	Id          int    `jsonapi:"primary,customer_file"`
	FileId      int    `jsonapi:"attr,file_id"`
	Description string `jsonapi:"attr,file_format"`
}
type PromotionalFile struct {
	Id          int    `jsonapi:"primary,promotional_file"`
	FileId      int    `jsonapi:"attr,file_id"`
	Description string `jsonapi:"attr,file_format"`
}
type TextureFile struct {
	Id          int    `jsonapi:"primary,texture_file"`
	FileId      int    `jsonapi:"attr,file_id"`
	Description string `jsonapi:"attr,file_format"`
}
type ViewerFile struct {
	Id          int    `jsonapi:"primary,viewer_file"`
	FileId      int    `jsonapi:"attr,file_id"`
	Description string `jsonapi:"attr,file_format"`
}
type Certification struct {
	Id   string `jsonapi:"primary,certification"`
	Type string `jsonapi:"attr,certification_id"`
}
type Product struct {
	Id    int    `jsonapi:"primary,product"`
	Draft *Draft `jsonapi:"relation,draft"`
}

func (draft *Draft) createDraft(settings Settings) error {
	if settings.Debug {
		log.Printf("Create Draft")
	}
	var message bytes.Buffer
	if err := jsonapi.MarshalPayload(&message, draft); err != nil {
		log.Fatal("Error building Create Draft message: ", err)
		return err
	}

	url := fmt.Sprintf("%s/api/drafts", settings.Server)
	req, err := http.NewRequest("POST", url, bytes.NewReader(message.Bytes()))
	if err != nil {
		log.Fatal("Error building request for Create Draft: ", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.Token))
	req.Header.Add("Accept", "application/vnd.api+json; com.turbosquid.api.version=1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing request for Create Draft: ", err)
		return err
	}

	defer resp.Body.Close()

	if settings.Debug {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		bodyStr := buf.String()

		log.Printf(resp.Status)
		log.Printf(bodyStr)
	}

	err = jsonapi.UnmarshalPayload(resp.Body, draft)
	if draft.Id > 0 {
		log.Printf("Draft ID %d", draft.Id)
	}
	return err
}

func (draft *Draft) addFile(file File, settings Settings) error {
	if settings.Debug {
		log.Printf("Adding file: %d", file.FileId)
	}
	var message bytes.Buffer
	if file.Type == "product_file" {
		draftFile := &ProductFile{
			FileId:          file.FileId,
			Format:          file.Format,
			FormatVersion:   file.FormatVersion,
			Renderer:        file.Renderer,
			RendererVersion: file.RendererVersion,
			Native:          file.Native,
		}
		if err := jsonapi.MarshalPayload(&message, draftFile); err != nil {
			log.Fatal("Error building product_file message: ", err)
			return err
		}
	} else if file.Type == "customer_file" {
		draftFile := &CustomerFile{
			FileId:      file.FileId,
			Description: file.Description,
		}
		if err := jsonapi.MarshalPayload(&message, draftFile); err != nil {
			log.Fatal("Error building customer_file message: ", err)
			return err
		}
	} else if file.Type == "promotional_file" {
		draftFile := &PromotionalFile{
			FileId:      file.FileId,
			Description: file.Description,
		}
		if err := jsonapi.MarshalPayload(&message, draftFile); err != nil {
			log.Fatal("Error building promotional_file message: ", err)
			return err
		}
	} else if file.Type == "texture_file" {
		draftFile := &TextureFile{
			FileId:      file.FileId,
			Description: file.Description,
		}
		if err := jsonapi.MarshalPayload(&message, draftFile); err != nil {
			log.Fatal("Error building texture_file message: ", err)
			return err
		}
	} else if file.Type == "viewer_file" {
		draftFile := &ViewerFile{
			FileId:      file.FileId,
			Description: file.Description,
		}
		if err := jsonapi.MarshalPayload(&message, draftFile); err != nil {
			log.Fatal("Error building viewer_file message: ", err)
			return err
		}
	}

	url := fmt.Sprintf("%s/api/drafts/%d/%ss", settings.Server, draft.Id, file.Type)
	req, err := http.NewRequest("POST", url, bytes.NewReader(message.Bytes()))
	if err != nil {
		log.Fatal("Error building request for Add File: ", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.Token))
	req.Header.Add("Accept", "application/vnd.api+json; com.turbosquid.api.version=1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing request for Add File: ", err)
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 {
		log.Fatal("Failed to add file")
	}

	return nil
}

func (draft *Draft) addThumbnail(preview Preview, settings Settings) error {
	if settings.Debug {
		log.Printf("Adding preview: %s", preview.Name)
	}
	thumbnail := &Thumbnail{
		FileId: preview.FileId,
		Type:   preview.ThumbnailType,
	}

	var message bytes.Buffer
	if err := jsonapi.MarshalPayload(&message, thumbnail); err != nil {
		log.Fatal("Error building thumbnail message: ", err)
		return err
	}

	url := fmt.Sprintf("%s/api/drafts/%d/%ss", settings.Server, draft.Id, preview.Type)
	req, err := http.NewRequest("POST", url, bytes.NewReader(message.Bytes()))
	if err != nil {
		log.Fatal("Error building request for thumbnail: ", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.Token))
	req.Header.Add("Accept", "application/vnd.api+json; com.turbosquid.api.version=1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing request for Add Thumbnail: ", err)
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 {
		log.Fatal("Failed to add preview")
	}

	return nil
}

func (draft *Draft) addTurntable(preview Preview, settings Settings) error {
	if settings.Debug {
		log.Printf("Adding turntable: %s", preview.Name)
	}
	turntable := &Turntable{
		FileIds: preview.FileIds,
		Type:    preview.ThumbnailType,
	}

	var message bytes.Buffer
	if err := jsonapi.MarshalPayload(&message, turntable); err != nil {
		log.Fatal("Error building turntable message: ", err)
		return err
	}

	url := fmt.Sprintf("%s/api/drafts/%d/%ss", settings.Server, draft.Id, preview.Type)
	req, err := http.NewRequest("POST", url, bytes.NewReader(message.Bytes()))
	if err != nil {
		log.Fatal("Error building request for turntable: ", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.Token))
	req.Header.Add("Accept", "application/vnd.api+json; com.turbosquid.api.version=1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing request for Add Turntable: ", err)
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 {
		log.Fatal("Failed to add turntable")
	}

	return nil
}

func (draft *Draft) certifications(settings Settings, certifications []string) error {

	for _, certificationType := range certifications {
		if settings.Debug {
			log.Printf("Add certification: %s", certificationType)
		}
		url := fmt.Sprintf("%s/api/drafts/%d/certifications", settings.Server, draft.Id)

		certification := &Certification{
			Type: certificationType,
		}
		var message bytes.Buffer
		if err := jsonapi.MarshalPayload(&message, certification); err != nil {
			log.Fatal("Error building certification message", err)
			return err
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(message.Bytes()))
		if err != nil {
			log.Fatal("Error building request for certification", err)
			return err
		}
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.Token))
		req.Header.Add("Accept", "application/vnd.api+json; com.turbosquid.api.version=1")
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error performing request for Certification", err)
			return err
		}

		if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 {
			log.Fatal("Failed to set certification")
		}
	}

	return nil
}

func (draft *Draft) publish(settings Settings) (error, int) {
	log.Printf("Publish draft")
	url := fmt.Sprintf("%s/api/products", settings.Server)

	var product Product
	product.Draft = draft

	var message bytes.Buffer
	if err := jsonapi.MarshalPayloadWithoutIncluded(&message, &product); err != nil {
		log.Fatal("Error building product message: ", err)
		return err, 0
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(message.Bytes()))
	if err != nil {
		log.Fatal("Error building request for publish", err)
		return err, 0
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", settings.Token))
	req.Header.Add("Accept", "application/vnd.api+json; com.turbosquid.api.version=1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error performing request for Publish", err)
		return err, 0
	}

	defer resp.Body.Close()

	if err = jsonapi.UnmarshalPayload(resp.Body, product); err != nil {
		return err, 0
	}
	return nil, product.Id
}
