package handler

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/opchaves/gin-web-app/app/model"
)

// parseError takes an error or multiple errors and attempts to determine the best path to convert them into
// human readable strings
func parseError(errs ...error) []model.FieldError {
	var out []model.FieldError
	for _, err := range errs {
		switch typedError := any(err).(type) {
		case validator.ValidationErrors:
			// if the type is validator.ValidationErrors then it's actually an array of validator.FieldError so we'll
			// loop through each of those and convert them one by one
			for _, e := range typedError {
				fieldErr := model.FieldError{
					Field:   e.Field(),
					Message: parseFieldError(e),
				}
				out = append(out, fieldErr)
			}
		case *json.UnmarshalTypeError:
			// similarly, if the error is an unmarshalling error we'll parse it into another, more readable string format
			out = append(out, parseMarshallingError(*typedError))
		default:
			out = append(out, model.FieldError{
				Field:   "Error",
				Message: err.Error(),
			})
		}
	}
	return out
}

func parseFieldError(e validator.FieldError) string {
	// workaround to the fact that the `gt|gtfield=Start` gets passed as an entire tag for some reason
	// https://github.com/go-playground/validator/issues/926
	fieldPrefix := fmt.Sprintf("The field %s", e.Field())
	tag := strings.Split(e.Tag(), "|")[0]
	switch tag {
	case "required_without":
		return fmt.Sprintf("%s is required if %s is not supplied", fieldPrefix, e.Param())
	case "email":
		return fmt.Sprintf("email is not valid")
	case "min":
		if e.Type().Name() == "string" {
			return fmt.Sprintf("cannot have less than %s characters", e.Param())
		}
		return fmt.Sprintf("min value is %s, %s", e.Param(), e.Type().String())
	case "max":
		if e.Type().Name() == "string" {
			return fmt.Sprintf("cannot be longer than %s characters", e.Param())
		}
		return fmt.Sprintf("max value is %s", e.Param())
	case "lt", "ltfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be less than %s", fieldPrefix, param)
	case "gt", "gtfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be greater than %s", fieldPrefix, param)
	default:
		// if it's a tag for which we don't have a good format string yet we'll try using the default english translator
		english := en.New()
		translator := ut.New(english, english)
		if translatorInstance, found := translator.GetTranslator("en"); found {
			return e.Translate(translatorInstance)
		} else {
			return fmt.Errorf("%v", e).Error()
		}
	}
}

func parseMarshallingError(e json.UnmarshalTypeError) model.FieldError {
	return model.FieldError{
		Field:   e.Field,
		Message: fmt.Sprintf("must be a %s", e.Type.String()),
	}
}
