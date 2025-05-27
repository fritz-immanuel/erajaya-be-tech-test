package appcontext

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	// KeyURLPath represents the url path key in http server context
	KeyURLPath contextKey = "URLPath"

	// KeyHTTPMethodName represents the method name key in http server context
	KeyHTTPMethodName contextKey = "HTTPMethodName"

	// KeySessionID represents the current logged-in SessionID
	KeySessionID contextKey = "SessionID"

	// KeyUserID represents the current logged-in User ID
	KeyUserID contextKey = "UserID"

	// KeyUserName represents the current logged-in User Name
	KeyUserName contextKey = "UserName"

	// KeyEmployeeID represents the current logged-in Employee ID
	KeyEmployeeID contextKey = "EmployeeID"

	// KeyEmployeeName represents the current logged-in Employee Name
	KeyEmployeeName contextKey = "EmployeeName"

	// KeyCompanyID represents the current logged-in CompanyID
	KeyCompanyID contextKey = "CompanyID"

	// KeyBusinessID represents the current logged-in BusinessID
	KeyBusinessID contextKey = "BusinessID"

	// KeyIsSupervisor represents the current logged-in Supervisor Status
	KeyIsSupervisor contextKey = "IsSupervisor"

	KeyTokenExpiryTime contextKey = "Exp"

	// KeyLoginToken represents the current logged-in token
	KeyLoginToken contextKey = "LoginToken"

	// Type represents the current logged-in UserID's CustomerID from customer-payfazz
	KeyType contextKey = "Type"

	// KeyVersionCode represents the current version code of request
	KeyVersionCode contextKey = "VersionCode"

	// KeyCurrentXAccessToken represents the current access token of request
	KeyCurrentXAccessToken contextKey = "CurrentAccessToken"

	KeyRequestStatus contextKey = "KeyRequestStatus"

	// KeyRequestHeader represents the header of the request
	KeyRequestHeader contextKey = "KeyRequestHeader"

	// KeyRequestBody represents the body of the request
	KeyRequestBody contextKey = "KeyRequestBody"
)

// RequestStatus gets request status from context
func RequestStatus(ctx *gin.Context) *string {
	requestStatus := (*ctx).Value(KeyRequestStatus)
	if requestStatus != nil {
		v := requestStatus.(string)
		return &v
	}
	return nil
}

// RequestHeader gets client request header
func RequestHeader(ctx *gin.Context) string {
	requestHeader := (*ctx).Value(KeyRequestHeader)
	if requestHeader != nil {
		v := requestHeader.(string)
		return v
	}
	return ""
}

// RequestBody gets client request body
func RequestBody(ctx *gin.Context) interface{} {
	requestBody := (*ctx).Value(KeyRequestBody)
	if requestBody != nil {
		v := requestBody.(interface{})
		return v
	}
	return nil
}

// URLPath gets the data url path from the context
func URLPath(ctx *gin.Context) *string {
	urlPath := ctx.Value(fmt.Sprintf("%s", KeyURLPath))
	if urlPath != nil {
		v := urlPath.(string)
		return &v
	}
	return nil
}

// HTTPMethodName gets the data http method from the context
func HTTPMethodName(ctx *gin.Context) *string {
	httpMethodName := ctx.Value(fmt.Sprintf("%s", KeyHTTPMethodName))
	if httpMethodName != nil {
		v := httpMethodName.(string)
		return &v
	}
	return nil
}

// SessionID gets the data session id from the context
func SessionID(ctx *gin.Context) *string {
	sessionID := ctx.Value(fmt.Sprintf("%s", KeySessionID))
	if sessionID != nil {
		v := sessionID.(string)
		return &v
	}
	return nil
}

// UserID gets current userId logged in from the context
func UserID(ctx *gin.Context) *string {
	userID := ctx.Value(fmt.Sprintf("%v", KeyUserID))
	if userID != nil {
		if reflect.ValueOf(userID).Kind().String() == "string" {
			v := userID.(string)
			return &v
		} else {
			v := fmt.Sprintf("%v", userID)
			return &v
		}
	}
	return nil
}

// UserName gets current userName logged in from the context
func UserName(ctx *gin.Context) *string {
	userID := ctx.Value(fmt.Sprintf("%v", KeyUserName))
	if userID != nil {
		if reflect.ValueOf(userID).Kind().String() == "string" {
			v := userID.(string)
			return &v
		} else {
			v := fmt.Sprintf("%v", userID)
			return &v
		}
	}
	return nil
}

// EmployeeID gets current userId logged in from the context
func EmployeeID(ctx *gin.Context) *string {
	userID := ctx.Value(fmt.Sprintf("%v", KeyEmployeeID))
	if userID != nil {
		if reflect.ValueOf(userID).Kind().String() == "string" {
			v := userID.(string)
			return &v
		} else {
			v := fmt.Sprintf("%v", userID)
			return &v
		}
	}
	return nil
}

// EmployeeName gets current userName logged in from the context
func EmployeeName(ctx *gin.Context) *string {
	employeeName := ctx.Value(fmt.Sprintf("%v", KeyEmployeeName))
	if employeeName != nil {
		if reflect.ValueOf(employeeName).Kind().String() == "string" {
			v := employeeName.(string)
			return &v
		} else {
			v := fmt.Sprintf("%v", employeeName)
			return &v
		}
	}
	return nil
}

// TypeID gets current TypeID logged in from the context
func Type(ctx *gin.Context) *string {
	typeData := ctx.Value(fmt.Sprintf("%s", KeyType))

	if typeData != nil {
		v := typeData.(string)
		return &v
	}
	return nil
}

// CompanyID gets current prefered CompanyID of Employee or User ID
func CompanyID(ctx *gin.Context) *string {
	companyID := ctx.Value(fmt.Sprintf("%v", KeyCompanyID))
	if companyID != nil {
		if reflect.ValueOf(companyID).Kind().String() == "string" {
			v := companyID.(string)
			return &v
		} else {
			v := fmt.Sprintf("%v", companyID)
			return &v
		}
	}
	return nil
}

// BusinessID gets current prefered BusinessID of Employee or User ID
func BusinessID(ctx *gin.Context) *string {
	businessID := ctx.Value(fmt.Sprintf("%v", KeyBusinessID))
	if businessID != nil {
		if reflect.ValueOf(businessID).Kind().String() == "string" {
			v := businessID.(string)
			return &v
		} else {
			v := fmt.Sprintf("%v", businessID)
			return &v
		}
	}
	return nil
}

func IsSupervisor(ctx *gin.Context) int {
	isSupervisor := ctx.Value(fmt.Sprintf("%s", KeyIsSupervisor))
	if isSupervisor != nil {
		v := int(isSupervisor.(float64))
		return v
	}

	return 0
}

func TokenExpiryTime(ctx *gin.Context) *time.Time {
	tokenExpiryTime := ctx.Value(fmt.Sprintf("%v", KeyTokenExpiryTime))
	if tokenExpiryTime != nil {
		if reflect.ValueOf(tokenExpiryTime).Kind().String() == "string" {
			v := tokenExpiryTime.(string)
			parsed, _ := time.Parse("2006-01-02T15:04:05.999999-07:00", v)
			return &parsed
		} else {
			v := fmt.Sprintf("%v", tokenExpiryTime)
			parsed, _ := time.Parse("2006-01-02T15:04:05.999999-07:00", v)
			return &parsed
		}
	}

	return nil
}

// VersionCode gets current version code of request
func VersionCode(ctx *gin.Context) int {
	versionCode := ctx.Value(fmt.Sprintf("%s", KeyVersionCode))
	if versionCode != nil {
		v := int(versionCode.(float64))
		return v
	}
	return 0
}

// CurrentXAccessToken gets current x access token code of request
func CurrentXAccessToken(ctx *gin.Context) string {
	currentAccessToken := ctx.Value(fmt.Sprintf("%s", KeyCurrentXAccessToken))
	if currentAccessToken != nil {
		v := currentAccessToken.(string)
		return v
	}
	return ""
}
