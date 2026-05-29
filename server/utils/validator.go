package utils

import "github.com/go-playground/validator/v10"

var GlobalValidator = validator.New(validator.WithRequiredStructEnabled())
