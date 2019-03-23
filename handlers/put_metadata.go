package handlers

import (
	"net/http"

	"github.com/willcj33/yaml-metadata-exercise/db"
)

//PutMetadata updates application metadata
func PutMetadata(w http.ResponseWriter, r *http.Request, store *db.MetadataStore) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Data"))
}
