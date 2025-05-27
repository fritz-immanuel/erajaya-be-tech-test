package library

import (
	"regexp"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// isEmailValid checks if the email provided passes the required structure and length.
func IsEmailValid(e string) (string, *types.Error) {
	if len(e) < 3 && len(e) > 254 {
		return "", &types.Error{
			Path:       ".IsEmailValid()",
			Message:    "Email is not valid",
			Error:      nil,
			StatusCode: 0,
			Type:       "validation-error",
		}
	}

	if emailRegex.MatchString(e) {
		return e, nil
	}

	return "", &types.Error{
		Path:       ".IsEmailValid()",
		Message:    "Email is not valid",
		Error:      nil,
		StatusCode: 0,
		Type:       "validation-error",
	}
}
