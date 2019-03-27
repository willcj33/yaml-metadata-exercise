package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/willcj33/yaml-metadata-exercise/config"
	"github.com/willcj33/yaml-metadata-exercise/db"
	"github.com/willcj33/yaml-metadata-exercise/models"
	yaml "gopkg.in/yaml.v2"
)

//PostMetadata creates application metadata
func PostMetadata(w http.ResponseWriter, r *http.Request, store *db.MetadataStore, config config.Config) {
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

	if validationErrors := applicationMetadata.Validate(); validationErrors != nil && len(validationErrors) > 0 {
		var returnArray = make([]string, len(validationErrors))
		i := 0
		for k, v := range validationErrors {
			if applicationMetadata.GetPascalField(k) == "" {
				returnArray[i] = fmt.Sprintf("%s", v)
			} else {
				returnArray[i] = fmt.Sprintf("%s -- %s", v, applicationMetadata.GetPascalField(k))
			}
			i++
		}
		var errorObject = map[string][]string{"Request contains errors": returnArray}
		var b []byte
		switch r.URL.Query().Get("format") {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			b, _ = json.Marshal(errorObject)
		default:
			w.Header().Set("Content-Type", "application/x-yaml")
			b, _ = yaml.Marshal(errorObject)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(b))
		return
	}

	id := applicationMetadata.GetID(config)
	store.Write(id, applicationMetadata)
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}
