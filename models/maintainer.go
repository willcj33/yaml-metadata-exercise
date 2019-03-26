package models

type Maintainer struct {
	Name  string `validate:"required" yaml:"name" json:"name"`
	Email string `validate:"required,email" yaml:"email" json:"email"`
}
