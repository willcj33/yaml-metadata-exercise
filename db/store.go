package db

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/willcj33/yaml-metadata-exercise/models"

	"github.com/blevesearch/bleve"

	"github.com/blevesearch/bleve/mapping"
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
		if index, err := bleve.Open(config.IndexName); err != nil {
			mapping := createMappings()
			newIndex, _ := bleve.New(config.IndexName, mapping)
			storeInstance = &MetadataStore{
				index: newIndex,
			}
		} else {
			storeInstance = &MetadataStore{
				index: index,
			}
		}
	})
	return storeInstance
}

//Query searches the metadata
func (store *MetadataStore) Query(q string) (*bleve.SearchResult, error) {
	resultMap := sync.Map{}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		var search *bleve.SearchRequest
		spl := strings.Split(q, ":")
		finalQ := spl[0]
		if len(spl) > 1 {
			finalQ = spl[1]
		}
		search = bleve.NewSearchRequest(bleve.NewWildcardQuery(fmt.Sprintf("*%s*", finalQ)))
		search.IncludeLocations = true
		search.Fields = []string{"*"}
		res, _ := store.index.Search(search)
		resultMap.Store("wildcard", res)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		var search *bleve.SearchRequest
		search = bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
		search.IncludeLocations = true
		search.Fields = []string{"*"}
		res, _ := store.index.Search(search)
		resultMap.Store("queryString", res)
		wg.Done()
	}()

	wg.Wait()

	var results *bleve.SearchResult
	resultMap.Range(func(k interface{}, v interface{}) bool {
		value := v.(*bleve.SearchResult)
		if results == nil || results.Total == 0 || (value.Total < results.Total && value.Total != 0) {
			results = value
		}
		return true
	})
	return results, nil
}

//Write searches the metadata
func (store *MetadataStore) Write(id string, data interface{}) error {
	b, _ := json.Marshal(data)
	store.index.SetInternal([]byte(id), b)
	return store.index.Index(id, data)
}

//GetDocument gets the original document stored
func (store *MetadataStore) GetDocument(id string) *models.ApplicationMetadata {
	b, _ := store.index.GetInternal([]byte(id))
	returnObj := &models.ApplicationMetadata{}
	json.Unmarshal(b, returnObj)
	return returnObj
}

//Delete deletes the original document stored and the index
func (store *MetadataStore) Delete(id string) {
	store.index.Delete(id)
	store.index.DeleteInternal([]byte(id))
}

func createMappings() *mapping.IndexMappingImpl {
	indexMapping := bleve.NewIndexMapping()

	metadataMapping := bleve.NewDocumentMapping()

	titleMapping := bleve.NewTextFieldMapping()
	metadataMapping.AddFieldMappingsAt("title", titleMapping)

	versionMapping := bleve.NewTextFieldMapping()
	metadataMapping.AddFieldMappingsAt("version", versionMapping)

	maintainerMapping := bleve.NewDocumentMapping()
	maintainerNameMapping := bleve.NewTextFieldMapping()
	maintainerMapping.AddFieldMappingsAt("name", maintainerNameMapping)
	maintainerEmailMapping := bleve.NewTextFieldMapping()
	maintainerMapping.AddFieldMappingsAt("email", maintainerEmailMapping)
	metadataMapping.AddSubDocumentMapping("maintainer", maintainerMapping)

	companyMapping := bleve.NewTextFieldMapping()
	metadataMapping.AddFieldMappingsAt("company", companyMapping)

	websiteMapping := bleve.NewTextFieldMapping()
	metadataMapping.AddFieldMappingsAt("website", websiteMapping)

	sourceMapping := bleve.NewTextFieldMapping()
	metadataMapping.AddFieldMappingsAt("source", sourceMapping)

	licenseMapping := bleve.NewTextFieldMapping()
	metadataMapping.AddFieldMappingsAt("license", licenseMapping)

	descriptionMapping := bleve.NewTextFieldMapping()
	metadataMapping.AddFieldMappingsAt("description", descriptionMapping)

	indexMapping.AddDocumentMapping("metadata", metadataMapping)
	return indexMapping
}
