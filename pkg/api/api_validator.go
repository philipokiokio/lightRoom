package api

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

// InitializeValidator initializes the validator instance.
func InitializeValidator() {
	validate = validator.New()
}
