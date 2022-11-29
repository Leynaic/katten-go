package models

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/fr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	frtrad "github.com/go-playground/validator/v10/translations/fr"
)

type ErrorResponse struct {
	FailedField string `json:"field"`
	Message     string `json:"message"`
}

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func InitErrors() {
	fr := fr.New()
	uni = ut.New(fr)
	trans = uni.GetFallback()
	validate = validator.New()
	frtrad.RegisterDefaultTranslations(validate, trans)
}

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(s)

	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&validator.InvalidValidationError{}) {
			validationErr := err.(*validator.InvalidValidationError)
			var element ErrorResponse
			element.FailedField = validationErr.Type.Name()
			element.Message = validationErr.Error()
			errors = append(errors, &element)
		} else {
			for key, err := range err.(validator.ValidationErrors).Translate(trans) {
				var element ErrorResponse
				fullField := strings.Split(key, ".")
				element.FailedField = fullField[len(fullField)-1]
				element.Message = err
				errors = append(errors, &element)
			}
		}
	}
	return errors
}
