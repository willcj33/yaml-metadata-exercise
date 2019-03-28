package db

import (
	"encoding/json"
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
	var search *bleve.SearchRequest
	search = bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
	search.IncludeLocations = true
	search.Fields = []string{"*"}
	return store.index.Search(search)
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
	titleMapping.Analyzer = "en"
	metadataMapping.AddFieldMappingsAt("title", titleMapping)

	versionMapping := bleve.NewTextFieldMapping()
	versionMapping.Analyzer = "en"
	metadataMapping.AddFieldMappingsAt("version", versionMapping)

	maintainerMapping := bleve.NewDocumentMapping()
	maintainerNameMapping := bleve.NewTextFieldMapping()
	maintainerNameMapping.Analyzer = "en"
	maintainerMapping.AddFieldMappingsAt("name", maintainerNameMapping)
	maintainerEmailMapping := bleve.NewTextFieldMapping()
	maintainerEmailMapping.Analyzer = "en"
	maintainerMapping.AddFieldMappingsAt("email", maintainerEmailMapping)
	metadataMapping.AddSubDocumentMapping("maintainer", maintainerMapping)

	companyMapping := bleve.NewTextFieldMapping()
	companyMapping.Analyzer = "en"
	metadataMapping.AddFieldMappingsAt("company", companyMapping)

	websiteMapping := bleve.NewTextFieldMapping()
	websiteMapping.Analyzer = "en"
	metadataMapping.AddFieldMappingsAt("website", websiteMapping)

	sourceMapping := bleve.NewTextFieldMapping()
	sourceMapping.Analyzer = "en"
	metadataMapping.AddFieldMappingsAt("source", sourceMapping)

	licenseMapping := bleve.NewTextFieldMapping()
	licenseMapping.Analyzer = "en"
	metadataMapping.AddFieldMappingsAt("license", licenseMapping)

	descriptionMapping := bleve.NewTextFieldMapping()
	descriptionMapping.Analyzer = "en"
	metadataMapping.AddFieldMappingsAt("description", descriptionMapping)

	indexMapping.AddDocumentMapping("metadata", metadataMapping)
	return indexMapping
}
