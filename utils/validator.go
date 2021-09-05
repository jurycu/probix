package utils

import (
	"github.com/go-playground/validator/v10"
)

func Validator(structInput interface{}) error {
	validate := validator.New()
	err := validate.Struct(structInput)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return err
		}
	}
	return nil
}
