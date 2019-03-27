package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/willcj33/yaml-metadata-exercise/config"
	"github.com/willcj33/yaml-metadata-exercise/db"
	"github.com/willcj33/yaml-metadata-exercise/models"
)

//DeleteMetadata updates application metadata
func DeleteMetadata(w http.ResponseWriter, r *http.Request, store *db.MetadataStore, config config.Config) {
	var requestBody []byte
	applicationMetadata := &models.ApplicationMetadata{}

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	} else {
		requestBody = body
	}

	if err := applicationMetadata.FromYaml(requestBody); err != nil {
		http.Error(w, "can't parse yaml", http.StatusBadRequest)
		return
	}
	store.Delete(applicationMetadata.GetID(config))
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
