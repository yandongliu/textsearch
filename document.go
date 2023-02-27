package main

import (
	"compress/gzip"
	"encoding/xml"
	"os"
	"strconv"

	"github.com/yandongliu/textsearch/entities"
)

// loadDocuments loads a Wikipedia abstract dump and returns a slice of documents.
// Dump example: https://dumps.wikimedia.org/enwiki/latest/enwiki-latest-abstract1.xml.gz
func loadXMLDocuments(path string) ([]entities.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	dec := xml.NewDecoder(gz)
	dump := struct {
		Documents []entities.Document `xml:"doc"`
	}{}
	if err := dec.Decode(&dump); err != nil {
		return nil, err
	}
	docs := dump.Documents
	for i := range docs {
		docs[i].ID = i
	}
	return docs, nil
}

func loadJsonDocuments() ([]entities.Document, error) {
	n := 5
	docs := make([]entities.Document, n)
	for i := 0; i < n; i++ {
		docs[i].ID = i
		docs[i].Title = "Title" + strconv.Itoa(i)
		docs[i].URL = "URL" + strconv.Itoa(i)
		docs[i].Text = "Text" + strconv.Itoa(i)
	}

	return docs, nil
}
