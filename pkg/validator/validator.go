package validator

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error      bool          `json:"error,omitempty"`
	FieldError []*FieldError `json:"field_error,omitempty"`
}

type FieldError struct {
	FailedField string `json:"failed_field,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Value       string `json:"value,omitempty"`
}

var validate = validator.New()

func ValidateStruct(s interface{}) *ErrorResponse {
	var errors []*FieldError
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element FieldError
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return &ErrorResponse{
		Error:      true,
		FieldError: errors,
	}
}
