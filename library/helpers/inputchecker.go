package helpers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func MultiValueFilterCheck(input string) (string, *types.Error) {
	result := ""

	if input != "" {
		sanitizeSpace := strings.ReplaceAll(input, " ", "")
		splitComa := strings.Split(sanitizeSpace, ",")

		for _, s := range splitComa {
			num, err := strconv.Atoi(s)
			if err != nil {
				return "", &types.Error{
					Message:    "Unknown input data type",
					Error:      fmt.Errorf("unknown input data type"),
					StatusCode: http.StatusBadRequest,
				}
			}

			if reflect.TypeOf(num).String() == "int" {
				if result == "" {
					result = fmt.Sprintf("%d", num)
				} else {
					result = fmt.Sprintf("%s,%d", result, num)
				}
			}
		}
	}

	return result, nil
}

func ValidateStruct(input interface{}) *types.Error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	errValidation := validate.Struct(input)
	if errValidation != nil {
		missingTag := strings.Split(strings.Split(errValidation.Error(), ".")[1], "'")[0]
		msg := fmt.Sprintf(`'%s' is required`, missingTag)

		return &types.Error{
			Path:       ".Helpers->ValidateStruct()",
			Message:    msg,
			Error:      errValidation,
			StatusCode: http.StatusBadRequest,
			Type:       "validation-error",
		}
	}

	return nil
}

func ValidateUUID(input string) (string, *types.Error) {
	if _, err := uuid.Parse(input); err != nil {
		return "", &types.Error{
			Path:       ".Helpers->ValidateUUID()",
			Message:    "Unknown input data type",
			Error:      fmt.Errorf("unknown input data type"),
			StatusCode: http.StatusBadRequest,
			Type:       "validation-error",
		}
	}

	return input, nil
}

func MultiValueUUIDCheck(input string) (string, *types.Error) {
	result := ""

	if input != "" {
		explodeUUID := strings.Split(input, ",")
		for _, v := range explodeUUID {
			if _, err := uuid.Parse(v); err != nil {
				return "", &types.Error{
					Path:       ".Helpers->MultiValueUUIDCheck()",
					Message:    "Unknown input data type",
					Error:      fmt.Errorf("unknown input data type"),
					StatusCode: http.StatusBadRequest,
					Type:       "validation-error",
				}
			}

			result = result + v + ","
		}
		return strings.TrimSuffix(result, ","), nil
	}

	return result, nil
}
