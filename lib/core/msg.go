package core

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var NameOptional = []validation.Rule{
	validation.Length(1, 64),
	validation.Match(regexp.MustCompile("^[0-9A-Za-z_.-]*$")),
}

var NameRequired = append(NameOptional, validation.Required)
