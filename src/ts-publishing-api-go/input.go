package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type ProductBundle struct {
	Draft          Draft     `json:"product"`
	Files          []File    `json:"files"`
	Previews       []Preview `json:"previews"`
	Certifications []string  `json:"certifications"`
}

type Draft struct {
	Id           int      `jsonapi:"primary,draft,omitempty"`
	Name         string   `json:"name" jsonapi:"attr,name"`
	Type         string   `json:"product_type" jsonapi:"attr,product_type"`
	PriceUsd     float32  `json:"price_usd"`
	Price        Price    `jsonapi:"attr,price"`
	Description  string   `json:"description" jsonapi:"attr,description"`
	Status       string   `json:"status" jsonapi:"attr,status"`
	License      string   `json:"license" jsonapi:"attr,license"`
	Tags         []string `json:"tags" jsonapi:"attr,tags"`
	Animated     bool     `json:"animated" jsonapi:"attr,animated"`
	Geometry     string   `json:"geometry" jsonapi:"attr,geometry"`
	Materials    bool     `json:"materials" jsonapi:"attr,materials"`
	Polygons     int      `json:"polygons" jsonapi:"attr,polygons"`
	Rigged       bool     `json:"rigged" jsonapi:"attr,rigged"`
	Textures     bool     `json:"textures" jsonapi:"attr,textures"`
	UnwrappedUVs string   `json:"unwrapped_u_vs" jsonapi:"attr,unwrapped_u_vs"`
	UVMapped     bool     `json:"uv_mapped" jsonapi:"attr,uv_mapped"`
	Vertices     int      `json:"vertices" jsonapi:"attr,vertices"`
}

type Price struct {
	Value       int    `json:"value"`
	Currency    string `json:"currency"`
	Denominator int    `json:"demonminator"`
}

type File struct {
	FileId          int
	Name            string `json:"file_name"`
	Type            string `json:"type"`
	Format          string `json:"file_format"`
	FormatVersion   string `json:"format_version"`
	Renderer        string `json:"renderer"`
	RendererVersion string `json:"renderer_version"`
	Native          bool   `json:"is_native"`
	Description     string `json:"description"`
}

type Preview struct {
	FileId        int
	FileIds       []int
	Name          string `json:"file_name"`
	Type          string `json:"type"`
	ThumbnailType string `json:"thumbnail_type"`
}

func ReadInput(directory string) ProductBundle {
	productPath := fmt.Sprintf("%s/product.json", directory)
	jsonFile, err := ioutil.ReadFile(productPath)
	if err != nil {
		log.Fatalf(fmt.Sprintf("unable to read %s", productPath), err)
	}

	var pb ProductBundle
	if err = json.Unmarshal([]byte(jsonFile), &pb); err != nil {
		log.Fatalf("Unable to parse json file: ", err)
	}

	pb.Draft.Price.Currency = "USD"
	pb.Draft.Price.Denominator = 100
	pb.Draft.Price.Value = int(pb.Draft.PriceUsd * 100.0)

	return pb
}
