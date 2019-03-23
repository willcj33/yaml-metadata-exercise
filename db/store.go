package db

import (
	"sync"

	"github.com/blevesearch/bleve"

	"github.com/willcj33/yaml-metadata-exercise/config"
)

//MetadataStore contains storage utilities for the metadata
//TODO: potentially move this to a real DB in the future
type MetadataStore struct {
	index bleve.Index
}

var storeInstance *MetadataStore
var storeOnce sync.Once

//GetStore gets the singleton store instance
func GetStore(config config.Config) *MetadataStore {
	storeOnce.Do(func() {
		mapping := bleve.NewIndexMapping()
		i, _ := bleve.New("applicationMetadata.bleve", mapping)
		storeInstance = &MetadataStore{
			index: i,
		}
	})
	return storeInstance
}

//Query searches the metadata
func (store *MetadataStore) Query(q string) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(q)
	search := bleve.NewSearchRequest(query)
	return store.index.Search(search)
}

//Write searches the metadata
func (store *MetadataStore) Write(id string, data interface{}) error {
	return store.index.Index(id, data)
}
