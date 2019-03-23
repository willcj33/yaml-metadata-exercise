package models

type Maintainer struct {
	Name  string `validate:"required" yaml:"name"`
	Email string `validate:"required,email" yaml:"email"`
}
