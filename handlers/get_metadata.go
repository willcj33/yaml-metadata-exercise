package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/willcj33/yaml-metadata-exercise/config"
	"github.com/willcj33/yaml-metadata-exercise/db"
	"gopkg.in/yaml.v2"
)

//GetMetadata queries for application metadata
func GetMetadata(w http.ResponseWriter, r *http.Request, store *db.MetadataStore, config config.Config) {
	urlQuery := r.URL.Query()
	q := urlQuery.Get("query")
	field := urlQuery.Get("field")
	fill := urlQuery.Get("fill")
	finalQ := ""
	if field != "" {
		finalQ = fmt.Sprintf("%s:\"%s\"", field, q)
	} else {
		finalQ = fmt.Sprintf("\"%s\"", q)
	}
	if res, err := store.Query(finalQ); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		ids := make([]string, len(res.Hits))
		for i, hit := range res.Hits {
			ids[i] = hit.ID
		}
		if fill != "" {
			var results = make([]interface{}, len(res.Hits))
			for i, hit := range res.Hits {
				results[i] = store.GetDocument(hit.ID)
			}
			var b []byte
			switch urlQuery.Get("format") {
			case "json":
				b, _ = json.Marshal(results)
			default:
				b, _ = yaml.Marshal(results)
			}
			w.Write(b)
			return
		}
		var results = make([]map[string]interface{}, len(res.Hits))
		for i, hit := range res.Hits {
			results[i] = make(map[string]interface{}, len(hit.Locations))
			for key, loc := range hit.Locations {
				var idx *uint64
				for _, lValue := range loc {
					for _, l2Value := range lValue {
						if l2Value.ArrayPositions != nil {
							for _, p := range l2Value.ArrayPositions {
								idx = &p
							}
						}
					}
				}
				if idx != nil && reflect.TypeOf(hit.Fields[key]).String() != "string" {
					results[i][key] = hit.Fields[key].([]interface{})[int(*idx)]
				} else {
					results[i][key] = hit.Fields[key]
				}
			}
		}
		var b []byte
		w.WriteHeader(http.StatusOK)
		switch urlQuery.Get("format") {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			b, _ = json.Marshal(results)
		default:
			w.Header().Set("Content-Type", "application/x-yaml")
			b, _ = yaml.Marshal(results)
		}
		w.Write(b)
	}
}
