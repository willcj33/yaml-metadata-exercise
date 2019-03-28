package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/willcj33/yaml-metadata-exercise/config"
	"github.com/willcj33/yaml-metadata-exercise/models"
)

var store *MetadataStore

func setup() {
	store = GetStore(config.Config{
		IndexName: "test.bleve",
	})
}

func teardown() {
	fmt.Println("done")
	dir, _ := os.Getwd()
	os.RemoveAll(fmt.Sprintf("%s/test.bleve", dir))
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestWrite(t *testing.T) {
	err := store.Write("test", &models.ApplicationMetadata{Title: "Will Jaynes"})
	if err != nil {
		t.Errorf("Failed: expected to not err, got: %s", err.Error())
	}
}

func TestQuery(t *testing.T) {
	res, err := store.Query("will")
	if err != nil {
		t.Errorf("Failed: expected to not err, got: %s", err.Error())
	}
	if len(res.Hits) != 1 {
		t.Errorf("Failed: expected 1 hit, got: %d", len(res.Hits))
	}
}

func TestQueryField(t *testing.T) {
	res, err := store.Query("title:will")
	if err != nil {
		t.Errorf("Failed: expected to not err, got: %s", err.Error())
	}
	if len(res.Hits) != 1 {
		t.Errorf("Failed: expected 1 hit, got: %d", len(res.Hits))
	}
}

func TestDelete(t *testing.T) {
	store.Delete("test")
	res, err := store.Query("will")
	if err != nil {
		t.Errorf("Failed: expected to not err, got: %s", err.Error())
	}
	if len(res.Hits) != 0 {
		t.Errorf("Failed: expected 0 hits, got: %d", len(res.Hits))
	}
}
