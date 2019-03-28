package models

import (
	"encoding/base64"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/go-playground/validator.v9"

	"github.com/willcj33/yaml-metadata-exercise/config"
)

func TestGetCamelField(t *testing.T) {
	m := &ApplicationMetadata{
		Title: "some_title",
		Maintainers: []*Maintainer{
			&Maintainer{
				Email: "test@test.com",
			},
		},
	}

	f1 := m.GetCamelField("title")
	f2 := m.GetCamelField("maintainers[0].email")

	if f1 != "some_title" {
		t.Errorf("Name failed: expected %s, go %s", "some_title", f1)
	}

	if f2 != "test@test.com" {
		t.Errorf("Maintainers Email failed: expected %s, go %s", "test@test.com", f2)
	}
}

func TestGetPascalField(t *testing.T) {
	m := &ApplicationMetadata{
		Title: "some_title",
		Maintainers: []*Maintainer{
			&Maintainer{
				Email: "test@test.com",
			},
		},
	}

	f1 := m.GetPascalField("Title")
	f2 := m.GetPascalField("Maintainers[0].Email")

	if f1 != "some_title" {
		t.Errorf("Name failed: expected %s, go %s", "some_title", f1)
	}

	if f2 != "test@test.com" {
		t.Errorf("Maintainers Email failed: expected %s, go %s", "test@test.com", f2)
	}
}

func TestGetID(t *testing.T) {
	expected := base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{"some_title", "test@test.com"}, "::")))

	m := &ApplicationMetadata{
		Title: "some_title",
		Maintainers: []*Maintainer{
			&Maintainer{
				Email: "test@test.com",
			},
		},
	}

	id := m.GetID(config.Config{
		IdentifierFields: []string{"title", "maintainers.email"},
	})

	if id != expected {
		t.Errorf("ID failed: expected %s, got %s", expected, id)
	}
}

func TestRequiredValidation(t *testing.T) {
	vMap := validator.ValidationErrorsTranslations{
		"ApplicationMetadata.Version":     "Version is a required field",
		"ApplicationMetadata.Maintainers": "Maintainers is a required field",
		"ApplicationMetadata.Company":     "Company is a required field",
		"ApplicationMetadata.Website":     "Website is a required field",
		"ApplicationMetadata.Source":      "Source is a required field",
		"ApplicationMetadata.License":     "License is a required field",
		"ApplicationMetadata.Description": "Description is a required field",
		"ApplicationMetadata.Title":       "Title is a required field",
	}

	m := &ApplicationMetadata{}

	errs := m.Validate()

	if !reflect.DeepEqual(errs, vMap) {
		t.Errorf("ID failed: expected %v, got %v", vMap, errs)
	}
}

func TestInvalidEmail(t *testing.T) {
	vMap := validator.ValidationErrorsTranslations{
		"ApplicationMetadata.Title":                "Title is a required field",
		"ApplicationMetadata.Version":              "Version is a required field",
		"ApplicationMetadata.Maintainers[0].Name":  "Name is a required field",
		"ApplicationMetadata.Maintainers[0].Email": "Email must be a valid email address",
		"ApplicationMetadata.Company":              "Company is a required field",
		"ApplicationMetadata.Website":              "Website is a required field",
		"ApplicationMetadata.Source":               "Source is a required field",
		"ApplicationMetadata.License":              "License is a required field",
		"ApplicationMetadata.Description":          "Description is a required field",
	}

	m := &ApplicationMetadata{
		Maintainers: []*Maintainer{
			&Maintainer{
				Email: "bad",
			},
		},
	}

	errs := m.Validate()

	if !reflect.DeepEqual(errs, vMap) {
		t.Errorf("ID failed: expected %v, got %v", vMap, errs)
	}
}
