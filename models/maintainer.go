package models

//Maintainer is the holder struct for the maintainer information on the application metadata
type Maintainer struct {
	Name  string `validate:"required" yaml:"name" json:"name"`
	Email string `validate:"required,email" yaml:"email" json:"email"`
}
