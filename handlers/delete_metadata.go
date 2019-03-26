package handlers

import (
	"net/http"

	"github.com/willcj33/yaml-metadata-exercise/config"
	"github.com/willcj33/yaml-metadata-exercise/db"
)

//PutMetadata updates application metadata
func DeleteMetadata(w http.ResponseWriter, r *http.Request, store *db.MetadataStore, config config.Config) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Data"))
}
