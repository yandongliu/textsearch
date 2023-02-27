package main

import (
	"flag"
	"log"
	"time"

	"github.com/yandongliu/textsearch/common"
	"github.com/yandongliu/textsearch/entities"
	"github.com/yandongliu/textsearch/indexer"
	"github.com/yandongliu/textsearch/serializer"
)

var logger = common.GetDevLogger()

func main() {
	var query string
	flag.StringVar(&query, "q", "openai", "search query")
	flag.Parse()

	// Load docs from disk
	start := time.Now()
	var docs = new([]entities.Document)
	err := serializer.ReadGob("./docs.gob", docs)
	if err != nil {
		logger.Error("loading docs error: ", err)
	} else {
		logger.Infof("Loaded %d docs in %v", len(*docs), time.Since(start))
	}

	// Load index from disk
	start = time.Now()
	docsIndex, err := indexer.IndexerFromFile("./index.gob")
	if err != nil {
		logger.Error("loading index error: ", err)
	} else {
		logger.Infof("loaded index in %v ", time.Since(start))
	}

	start = time.Now()
	matchedIDs := docsIndex.Search(query)
	logger.Infof("Search found %d documents in %v", len(matchedIDs), time.Since(start))

	logger.Infof("Top 5 matches for %s:", query)
	for _, id := range matchedIDs[:5] {
		doc := (*docs)[id]
		log.Printf("%d\t%v\n", id, doc)
	}
}
