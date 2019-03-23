package models

import (
	"fmt"
	"strings"

	"github.com/willcj33/yaml-metadata-exercise/services"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

type ApplicationMetadata struct {
	Title       string        `validate:"required" yaml:"title"`
	Version     string        `validate:"required" yaml:"version"`
	Maintainers []*Maintainer `validate:"required,dive,required" yaml:"maintainers"`
	Company     string        `validate:"required" yaml:"company"`
	Website     string        `validate:"required,isURL" yaml:"website"`
	Source      string        `validate:"required,isURL" yaml:"source"`
	License     string        `validate:"required" yaml:"license"`
	Description string        `validate:"required" yaml:"description"`
}

func (am *ApplicationMetadata) GetID() string {
	s := fmt.Sprintf("%s::%s", strings.ToLower(strings.Replace(am.Title, " ", "_", -1)), am.Source)
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
