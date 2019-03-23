package handlers

import (
	"log"
	"net/http"

	"github.com/willcj33/yaml-metadata-exercise/db"
	"gopkg.in/yaml.v2"
)

//GetMetadata queries for application metadata
func GetMetadata(w http.ResponseWriter, r *http.Request, store *db.MetadataStore) {
	q := r.URL.Query().Get("query")
	if res, err := store.Query(q); err != nil {
		log.Printf("Error searching for metadata: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		ids := make([]string, len(res.Hits))
		for i, hit := range res.Hits {
			ids[i] = hit.ID
		}
		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		output, _ := yaml.Marshal(ids)
		w.Write(output)
	}
}
