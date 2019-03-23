package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/willcj33/yaml-metadata-exercise/db"
	"github.com/willcj33/yaml-metadata-exercise/models"
)

//PostMetadata creates application metadata
func PostMetadata(w http.ResponseWriter, r *http.Request, store *db.MetadataStore) {
	var requestBody []byte
	applicationMetadata := &models.ApplicationMetadata{}

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	} else {
		requestBody = body
	}

	if err := applicationMetadata.FromYaml(requestBody); err != nil {
		log.Printf("Error parsing yaml: %v", err)
		http.Error(w, "can't parse yaml", http.StatusBadRequest)
		return
	}

	if validationErrors := applicationMetadata.Validate(); validationErrors != nil && len(validationErrors) > 0 {
		log.Printf("Error parsing yaml: %v", validationErrors)
		http.Error(w, fmt.Sprintf("%v", validationErrors), http.StatusBadRequest)
		return
	}
	id := applicationMetadata.GetID()
	store.Write(id, applicationMetadata)
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))

}
