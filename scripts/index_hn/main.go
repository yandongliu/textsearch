package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/yandongliu/textsearch/common"
	"github.com/yandongliu/textsearch/entities"
	"github.com/yandongliu/textsearch/indexer"
	"github.com/yandongliu/textsearch/serializer"
)

var logger = common.GetDevLogger()

func main() {
	var path, query string
	flag.StringVar(&path, "p", "/home/yandong/work/whatsup-go/data/hackernews/", "hackernews files path")
	flag.StringVar(&query, "q", "ChatGPT", "search query")
	flag.Parse()

	// Load HN docs from disk
	start := time.Now()
	docs, err := loadHNJsonDocs(path)
	if err != nil {
		log.Fatal(err)
	}
	err = serializer.WriteGob("./docs.gob", docs)
	if err != nil {
		log.Fatal(err)
	}
	logger.Infof("Loaded %d documents in %v", len(docs), time.Since(start))

	start = time.Now()
	docsIndex := indexer.NewEmptyIndexer()
	docsIndex.AddDocs(docs)
	logger.Infof("Indexed %d documents in %v", len(docs), time.Since(start))

	start = time.Now()
	err = docsIndex.WriteIndexToFile("./index.gob")
	if err != nil {
		logger.Error("indexing error: ", err)
	}
	logger.Infof("Wrote index to disk in %v", time.Since(start))

	start = time.Now()
	indexer2, err := indexer.IndexerFromFile("./index.gob")
	if err != nil {
		logger.Error("loading index error: ", err)
	} else {
		logger.Infof("Loaded index from disk in %v", time.Since(start))
	}

	start = time.Now()
	matchedIDs := indexer2.Search(query)
	logger.Infof("Found %d documents in %v", len(matchedIDs), time.Since(start))

	for _, id := range matchedIDs[:3] {
		doc := docs[id]
		logger.Infof("%d\t%v\n", id, doc)
	}
}

// Load Hacker News json files
func loadHNJsonDocs(path string) ([]entities.Document, error) {
	logger.Debugln("path: ", path)
	hndocs := make([]entities.HNDocument, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			_hndocs, err := readJsonFile(path + f.Name())
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			hndocs = append(hndocs, _hndocs...)
		}
	}
	logger.Info("hndocs size: ", len(hndocs))
	docs := make([]entities.Document, len(hndocs))
	for i, hndoc := range hndocs {
		docs[i].ID = i
		docs[i].Type = entities.DocTypeHN
		docs[i].Title = hndoc.Title
		docs[i].URL = hndoc.Url
		docs[i].Text = hndoc.Title
	}

	return docs, nil
}

func readJsonFile(path string) ([]entities.HNDocument, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal("Error when opening file: ", err)
	}

	var docs []entities.HNDocument
	err = json.Unmarshal(content, &docs)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return docs, err
}
