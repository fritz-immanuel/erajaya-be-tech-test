package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/notif"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/telegram"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/pkg/errors"
)

// FieldError represents error message for each field
//
//swagger:model
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents error message
//
//swagger:model
type ErrorResponse struct {
	Code    string        `json:"code"`
	Status  string        `json:"Status"`
	Message string        `json:"Message"`
	Fields  []*FieldError `json:"fields"`
	Data    *DataError    `json:"Data"`
}

type DataError struct {
	Message string `json:"Message"`
	Status  int    `json:"Status"`
}

// MakeFieldError create field error object
func MakeFieldError(field string, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

// Error writes error http response
func Error(c *gin.Context, data string, status int, err types.Error) {
	var errorCode string

	if status == 0 {
		status = http.StatusInternalServerError
	}

	switch status {
	case http.StatusUnauthorized:
		errorCode = "Unauthorized"
	case http.StatusNotFound:
		errorCode = "NotFound"
	case http.StatusBadRequest:
		errorCode = "BadRequest"
	case http.StatusUnprocessableEntity:
		errorCode = "ValidationError"
	case http.StatusInternalServerError:
		errorCode = "InternalServerError"
	case http.StatusNotImplemented:
		errorCode = "NotImplemented"
	}

	errorFields := []*FieldError{}

	switch err.Error.(type) {
	case validator.ValidationErrors:
		for _, err := range err.Error.(validator.ValidationErrors) {
			e := MakeFieldError(
				err.Field(),
				err.ActualTag())

			errorFields = append(errorFields, e)
		}

		data = "Unprocessable Entity"
		errorCode = "UnprocessableEntity"
		status = http.StatusUnprocessableEntity
	}

	c.JSON(status, ErrorResponse{
		Code:    errorCode,
		Status:  "Warning",
		Message: data,
		Fields:  errorFields,
		Data:    nil,
	})

	if err.Error != nil {
		log.Printf("INFO: %v\n", err.Error.Error())
		log.Printf("DETAIL [%s - %s]: %s\n", err.Path, err.Type, err.Message)

		if err.StatusCode >= 500 {
			msg := fmt.Sprintf(`URL: %s%%0A%%0AINFO: %v%%0A%%0ADETAIL [%s - %s]: %s`, strings.ReplaceAll(c.Request.URL.String(), "&", "%26"), err.Error.Error(), err.Path, err.Type, err.Message)

			if err := telegram.Send(msg); err != nil {
				fmt.Println("NOTIFY TO TELE ERROR: ", err)
			}
		}
	}
}

// AdvancedError writes error http response with params
func AdvancedError(w http.ResponseWriter, slackNotifier notif.Notifier, logNotifier notif.Notifier, status int, err *types.Error) {
	var errorCode string

	switch status {
	case http.StatusUnauthorized:
		errorCode = "Unauthorized"
	case http.StatusNotFound:
		errorCode = "NotFound"
	case http.StatusBadRequest:
		errorCode = "BadRequest"
	case http.StatusUnprocessableEntity:
		errorCode = "UnprocessableEntity"
	case http.StatusInternalServerError:
		errorCode = "InternalServerError"
	}

	errorFields := []*FieldError{}
	switch err.Error.(type) {
	case validator.ValidationErrors:
		for _, err := range err.Error.(validator.ValidationErrors) {
			errorFields = append(errorFields,
				MakeFieldError(err.Field(), err.ActualTag()))
		}
		errorCode = "BadRequest"
		status = http.StatusBadRequest
	}

	bytes := []byte(err.Params)
	params := map[string]interface{}{}
	json.Unmarshal(bytes, &params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusInternalServerError {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    errorCode,
			Message: "server error",
			Fields:  errorFields,
		})

		if logNotifier != nil {
			message := fmt.Sprintf("path    : %s\n\nerror   : %s\n\ntype    : %s\n\nmessage : %s\n\nparams  : %v\n\nstatus  : Internal Server Error\n", err.Path, err.Error, err.Type, err.Message, params)
			if err := logNotifier.Notify(fmt.Sprintf("```%s```", message)); err != nil {
				log.Println("Failed to notify log using slack: ", err)
			}
		}

		if err.IsIgnore == false {
			log.Printf("INFO: %v\n", err.Error.Error())
			log.Printf("DETAIL [%s - %s]: %s\nPARAMS: %v\n", err.Path, err.Type, err.Message, err.Params)
		}

		errMessage := fmt.Sprintf("%+v\n%s", err, string(debug.Stack()))
		if slackNotifier != nil {
			if err := slackNotifier.Notify(fmt.Sprintf("```%s```", errMessage)); err != nil {
				log.Println("Failed to notify using slack: ", err)
			}
		}
	} else {
		if len(errorFields) > 0 {
			json.NewEncoder(w).Encode(ErrorResponse{
				Code:    errorCode,
				Message: "validation error",
				Fields:  errorFields,
			})
			return
		}

		if logNotifier != nil {
			message := fmt.Sprintf("path    : %s\n\nerror   : %s\n\ntype    : %s\n\nmessage : %s\n\nparams  : %v\n", err.Path, err.Error, err.Type, err.Message, params)
			if err := logNotifier.Notify(fmt.Sprintf("```%s```", message)); err != nil {
				log.Println("Failed to notify log using slack: ", err)
			}
		}

		if err.IsIgnore == false {
			log.Printf("INFO: %v\n", err.Error.Error())
			log.Printf("DETAIL [%s - %s]: %s\nPARAMS: %v\n", err.Path, err.Type, err.Message, err.Params)
		}

		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    errorCode,
			Message: err.Error.Error(),
			Fields:  errorFields,
		})
	}
}

// ErrorRPC writes error rpc response
func ErrorRPC(n notif.Notifier, err types.Error) {
	if err.Error != nil {
		log.Printf("INFO: %v\n", err.Error.Error())
		log.Printf("DETAIL [%s - %s]: %s\n", err.Path, err.Type, err.Message)
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}

		var st errors.StackTrace
		if err, ok := err.Error.(stackTracer); ok {
			st = err.StackTrace()
			fmt.Printf("INFO: %+v\n", st[0])
		}

		if n != nil {
			errMessage := fmt.Sprintf("ERROR: %v\n", err)
			if len(st) > 0 {
				errMessage = fmt.Sprintf("\n\nStack Trace: %v\n", st[0])
			}

			if err := n.Notify(fmt.Sprintf("```%s```", errMessage)); err != nil {
				fmt.Println("NOTIFY TO SLACK ERROR: ", err)
			}
		}
	}
}
