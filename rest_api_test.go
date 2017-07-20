package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ReqError - Base Error structure
type ReqError struct {
	Error ReqErrorItem `json:"error"`
}

// ReqErrorItem - Error item
type ReqErrorItem struct {
	ErrorID   string                 `json:"error_id"`
	ErrorText string                 `json:"error_text"`
	ErrorInfo map[string]interface{} `json:"error_info"`
}

func TestBaseApi(t *testing.T) {
	method := "myMethod"
	endpoint := "/myEndpoint"
	requestObject := string("")
	responseObject := string("")
	errorObject := new(ReqError)
	errorObject.Error.ErrorID = "resource.validation_failed"
	errorObject.Error.ErrorText = "Some of the properties in the resource failed validation."
	errorObject.Error.ErrorInfo = make(map[string]interface{})

	errorMap := make(map[string]interface{})
	errorItem := make(map[string]string)
	errorItem["error_id"] = "num.range"
	errorItem["error_text"] = "Value must be in range 1000 - 2000."
	errorMap["key1"] = errorItem
	errorObject.Error.ErrorInfo["basic"] = errorMap

	statusCode := 200
	rawResponse := []byte("some server response in []byte")
	err := errors.New("an error")

	api := NewBaseAPI(method, endpoint, requestObject, responseObject, errorObject)

	api.SetStatusCode(statusCode)
	api.SetRawResponse(rawResponse)
	api.SetResponseObject(responseObject)
	api.SetErrorObject(errorObject)
	api.SetError(err)

	assert.Equal(t, method, api.Method())
	assert.Equal(t, endpoint, api.Endpoint())
	assert.Equal(t, requestObject, api.RequestObject())
	assert.Equal(t, responseObject, api.ResponseObject())

	assert.Equal(t, statusCode, api.StatusCode())
	assert.Equal(t, rawResponse, api.RawResponse())
	assert.Equal(t, err, api.Error())
}

func TestReqError(t *testing.T) {
	var errObj ReqError

	jsonErr := []byte(`
{
        "error": {
                "error_id": "resource.validation_failed",
                "error_text": "Some of the properties in the resource failed validation.",
                "error_info": {
                        "basic": {
                                "key1": {
                                        "error_id": "num.range",
                                        "error_text": "Value must be in range 1000 - 2000."
                                }
                        }
                }
        }
}`)

	fmt.Println("Error structure as JSON:\n", string(jsonErr))

	err := json.Unmarshal(jsonErr, &errObj)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("Error struct:\n%+v", errObj)
	assert.Equal(t, "resource.validation_failed", errObj.Error.ErrorID)
}
