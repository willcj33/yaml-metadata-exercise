package models

import (
	"encoding/base64"
	"sort"
	"strconv"
	"strings"

	"github.com/willcj33/yaml-metadata-exercise/config"
	"github.com/willcj33/yaml-metadata-exercise/services"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

type ApplicationMetadata struct {
	Title       string        `validate:"required" yaml:"title" json:"title"`
	Version     string        `validate:"required" yaml:"version" json:"version"`
	Maintainers []*Maintainer `validate:"required,dive,required" yaml:"maintainers" json:"maintainers"`
	Company     string        `validate:"required" yaml:"company" json:"company"`
	Website     string        `validate:"required,isURL" yaml:"website" json:"website"`
	Source      string        `validate:"required,isURL" yaml:"source" json:"source"`
	License     string        `validate:"required" yaml:"license" json:"license"`
	Description string        `validate:"required" yaml:"description" json:"description"`
}

func (am *ApplicationMetadata) GetCamelField(key string) string {
	cleanStr := strings.Replace(strings.Replace(strings.Replace(key, "[", ".", 1), "]", "", 1), "applicationMetadata.", "", 1)
	spl := strings.Split(cleanStr, ".")
	switch spl[0] {
	case "title":
		return am.Title
	case "version":
		return am.Version
	case "maintainers":
		index, _ := strconv.ParseInt(spl[1], 10, 32)
		switch spl[2] {
		case "name":
			return am.Maintainers[index].Name
		case "email":
			return am.Maintainers[index].Email
		}
	case "company":
		return am.Company
	case "website":
		return am.Website
	case "source":
		return am.Source
	case "license":
		return am.License
	case "description":
		return am.Description
	}
	return ""
}

func (am *ApplicationMetadata) GetPascalField(key string) string {
	cleanStr := strings.Replace(strings.Replace(strings.Replace(key, "[", ".", 1), "]", "", 1), "ApplicationMetadata.", "", 1)
	spl := strings.Split(cleanStr, ".")
	switch spl[0] {
	case "Title":
		return am.Title
	case "Version":
		return am.Version
	case "Maintainers":
		index, _ := strconv.ParseInt(spl[1], 10, 32)
		switch spl[2] {
		case "Name":
			return am.Maintainers[index].Name
		case "Email":
			return am.Maintainers[index].Email
		}
	case "Company":
		return am.Company
	case "Website":
		return am.Website
	case "Source":
		return am.Source
	case "License":
		return am.License
	case "Description":
		return am.Description
	}
	return ""
}

func (am *ApplicationMetadata) GetID(cfg config.Config) string {
	if cfg.StorageMode == "single" {
		return "application_metadata"
	}
	identifierBuilder := []string{}
	for _, key := range cfg.IdentifierFields {
		if key == "maintainers.email" {
			for _, maintainer := range am.Maintainers {
				identifierBuilder = append(identifierBuilder, maintainer.Email)
			}
		} else if key == "maintainers.name" {
			for _, maintainer := range am.Maintainers {
				identifierBuilder = append(identifierBuilder, maintainer.Name)
			}
		} else {
			if am.GetCamelField(key) != "" {
				identifierBuilder = append(identifierBuilder, am.GetCamelField(key))
			}
		}
	}
	sort.Strings(identifierBuilder)
	s := base64.StdEncoding.EncodeToString([]byte(strings.Join(identifierBuilder, "::")))
	return s
}

func (am *ApplicationMetadata) FromYaml(content []byte) error {
	if err := yaml.Unmarshal(content, am); err != nil {
		return err
	}
	return nil
}

func (am *ApplicationMetadata) Validate() validator.ValidationErrorsTranslations {
	validationInstance := services.GetValidator()
	if err := validationInstance.Service.Struct(am); err != nil {
		return validationInstance.Translate(err)
	}
	return nil
}
