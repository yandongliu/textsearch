package indexer

import (
	"github.com/yandongliu/textsearch/entities"
	"github.com/yandongliu/textsearch/serializer"
	"github.com/yandongliu/textsearch/tokenizer"
)

type Indexer struct {
	Index *entities.Index
}

func NewEmptyIndexer() *Indexer {
	idx := make(entities.Index)
	return &Indexer{Index: &idx}
}

func IndexerFromFile(path string) (*Indexer, error) {
	idx := new(entities.Index)
	err := serializer.ReadGob(path, idx)
	if err != nil {
		return nil, err
	}
	return &Indexer{Index: idx}, err
}

func (indexer *Indexer) WriteIndexToFile(path string) error {
	err := serializer.WriteGob(path, indexer.Index)
	return err
}

// add adds documents to the index.
func (indexer *Indexer) AddDocs(docs []entities.Document) {
	for _, doc := range docs {
		indexer.AddDoc(doc)
	}
}

func (indexer *Indexer) AddDoc(doc entities.Document) {
	for _, token := range tokenizer.Analyze(doc.Text) {
		ids := (*indexer.Index)[token]
		if ids != nil && ids[len(ids)-1] == doc.ID {
			continue
		}
		(*indexer.Index)[token] = append(ids, doc.ID)
	}
}

// intersection returns the set intersection between a and b.
// a and b have to be sorted in ascending order and contain no duplicates.
func intersection(a []int, b []int) []int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	r := make([]int, 0, maxLen)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			r = append(r, a[i])
			i++
			j++
		}
	}
	return r
}

// search queries the index for the given text.
func (index *Indexer) Search(text string) []int {
	var r []int
	for _, token := range tokenizer.Analyze(text) {
		if ids, ok := (*index.Index)[token]; ok {
			if r == nil {
				r = ids
			} else {
				r = intersection(r, ids)
			}
		} else {
			return nil
		}
	}
	return r
}
