package errorutility

import (
	"encoding/json"
	"fmt"

	"strings"

	"github.com/go-playground/validator/v10"
)

func ParseError(errors ...error) map[string][]string {
	var errorList map[string][]string = map[string][]string{}

	for _, err := range errors {
		switch typedError := any(err).(type) {
		case validator.ValidationErrors:
			for _, e := range typedError {
				parseFieldError(e, errorList)
			}
		case *json.UnmarshalTypeError:
			parseUnmarshallingError(typedError, errorList)
		default:
			fmt.Printf("Error type %s\n", typedError)
			_, ok := errorList["others"]
			if ok {
				errorList["others"] = append(errorList["others"], err.Error())
			} else {
				errorList["others"] = []string{err.Error()}
			}
		}

	}

	return errorList
}

func parseFieldError(err validator.FieldError, errorData map[string][]string) {
	var errorMessage string
	switch err.Tag() {
	case "required":
		errorMessage = fmt.Sprintf("%s field is required.", err.Field())
	case "lt":
		errorMessage = fmt.Sprintf("%s field must be lower than predetermined value.", err.Field())
	case "lte":
		errorMessage = fmt.Sprintf("%s field must be lower than or equal to predetermined value.", err.Field())
	case "gt":
		errorMessage = fmt.Sprintf("%s field must be greater than predetermined value.", err.Field())
	case "gte":
		errorMessage = fmt.Sprintf("%s field must be greater than or equal to predetermined value.", err.Field())
	case "oneof":
		errorMessage = fmt.Sprintf("Value for field %s does not match the allowed values.", err.Field())
	case "email":
		errorMessage = fmt.Sprintf("Value for field %s must be a valid email.", err.Field())
	case "max":
		errorMessage = fmt.Sprint("Maximum length exceeded.")
	default:
		errorMessage = fmt.Errorf("%v", err).Error()
	}

	_, ok := errorData[err.Field()]
	if ok {
		errorData[err.Field()] = append(errorData[err.Field()], errorMessage)
	} else {
		errorData[err.Field()] = []string{errorMessage}
	}
}

func parseUnmarshallingError(err *json.UnmarshalTypeError, errorData map[string][]string) {
	_, ok := errorData[err.Field]
	if ok {
		errorData[err.Field] = append(errorData[err.Field], fmt.Sprintf("The field %s must be a %s", err.Field, err.Type.String()))
	} else {
		errorData[err.Field] = []string{fmt.Sprintf("The field %s must be a %s", err.Field, err.Type.String())}
	}
}

func DetectOtherError(err error) error {
	if strings.Contains(err.Error(), "mismatch between encode and decode") {
		return &HashidError{}
	}
	return err
}
