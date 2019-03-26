package models

import (
	"encoding/base64"
	"sort"
	"strings"

	"github.com/fatih/structs"
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

func (am *ApplicationMetadata) GetID(cfg config.Config) string {
	if cfg.StorageMode == "single" {
		return "application_metadata"
	}
	objectMap := structs.Map(am)
	identifierBuilder := []string{}
	for _, key := range cfg.IdentifierFields {
		if key == "Maintainers.Email" {
			for _, maintainer := range am.Maintainers {
				identifierBuilder = append(identifierBuilder, maintainer.Email)
			}
		} else if key == "Maintainers.Name" {
			for _, maintainer := range am.Maintainers {
				identifierBuilder = append(identifierBuilder, maintainer.Name)
			}
		} else {
			if objectMap[key] != nil {
				identifierBuilder = append(identifierBuilder, objectMap[key].(string))
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
