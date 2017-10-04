package rest

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var user = "nsxUser"
var password = "nsxPass"
var ignoreSSL = true
var debug = true
var client Client

var server *httptest.Server

const (
	unauthorizedStatusCode = http.StatusForbidden
	unauthorizedResponse   = "unauthorized"
)

func hasHeader(req *http.Request, name string, value string) bool {
	return req.Header.Get(name) == value
}

func setup(statusCode int, responseBody string) {
	basicAuthHeaderValue := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))
	server = httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !hasHeader(r, "Authorization", basicAuthHeaderValue) {
				w.WriteHeader(unauthorizedStatusCode)
				fmt.Fprint(w, unauthorizedResponse)
				return
			}
			w.WriteHeader(statusCode)
			fmt.Fprintln(w, responseBody)
		}))
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	client = Client{
		URL:       server.URL,
		User:      user,
		Password:  password,
		IgnoreSSL: ignoreSSL,
		Debug:     debug,
		Headers:   headers,
	}
}

func setupWrongHeader(statusCode int, responseBody string) {
	basicAuthHeaderValue := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))
	server = httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !hasHeader(r, "Authorization", basicAuthHeaderValue) {
				w.WriteHeader(unauthorizedStatusCode)
				fmt.Fprint(w, unauthorizedResponse)
				return
			}
			w.WriteHeader(statusCode)
			fmt.Fprintln(w, responseBody)
		}))
	headers := make(map[string]string)
	headers["Content-Type"] = "foo/bar"
	client = Client{
		URL:       server.URL,
		User:      user,
		Password:  password,
		IgnoreSSL: ignoreSSL,
		Debug:     debug,
		Headers:   headers,
	}
}

func TestHappyCase(t *testing.T) {
	setup(200, "pong")
	client = Client{
		URL:       server.URL,
		User:      user,
		Password:  password,
		IgnoreSSL: ignoreSSL,
		Debug:     debug,
	}
	apiRequest := NewBaseAPI(http.MethodGet, "/", nil, nil, nil)

	err := client.Do(apiRequest)

	assert.Nil(t, err)
}

// TODO: add TestFailWhenNotValidSSLCerts(t *testing.T)

func TestBasicAuthFailure(t *testing.T) {
	setup(0, "")
	client = Client{
		URL:       server.URL,
		User:      "invalidUser",
		Password:  "invalidPass",
		IgnoreSSL: ignoreSSL,
		Debug:     debug,
	}

	apiRequest := NewBaseAPI(http.MethodGet, "/", nil, nil, nil)
	err := client.Do(apiRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	assert.Equal(t, unauthorizedStatusCode, apiRequest.StatusCode())
	assert.Equal(t, unauthorizedResponse, string(apiRequest.RawResponse()))
}

func TestHttpReq(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	client := Client{
		URL:   ts.URL,
		Debug: true,
	}
	api := NewBaseAPI(http.MethodGet, "/", nil, new(string), nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, api.StatusCode())
	resp := *api.ResponseObject().(*string)
	assert.Equal(t, "Hello, client\n", resp)
	ts.Close()
}

type JSONFoo struct {
	Fields map[string]string `json:"fields"`
}

type XMLFoo struct {
	XMLName xml.Name `xml:"fields"`
	Foo     string   `xml:"foo"`
}

func TestHttpJSONResp(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"fields":{"foo":"bar"}}`))
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	client := Client{URL: ts.URL, Debug: true}

	api := NewBaseAPI(http.MethodGet, "/", nil, new(JSONFoo), nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, api.StatusCode())
	resp := *api.ResponseObject().(*JSONFoo)
	assert.Equal(t, "bar", resp.Fields["foo"])
}

func TestHttpXMLResp(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><fields><foo>bar</foo></fields>`))
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/xml"
	client := Client{URL: ts.URL, Debug: true, Headers: headers}

	out := new(XMLFoo)
	api := NewBaseAPI(http.MethodGet, "/", nil, out, nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, api.StatusCode())
	assert.Equal(t, "bar", out.Foo)
}

func TestHttpOctetStreamResp(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte(`ouitybw50ybvqy9yt8b6983p8v93`))
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/octet-stream"
	client := Client{URL: ts.URL, Debug: true, Headers: headers}

	out := new([]byte)
	api := NewBaseAPI(http.MethodGet, "/", nil, out, nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, api.StatusCode())
	assert.Equal(t, []byte("ouitybw50ybvqy9yt8b6983p8v93"), *out)
}

func TestHttpPlainTextResp(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/plain")
		w.Write([]byte("plain text"))
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/octet-stream"
	client := Client{URL: ts.URL, Debug: true, Headers: headers}

	out := new(string)
	api := NewBaseAPI(http.MethodGet, "/", nil, out, nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, api.StatusCode())
	assert.Equal(t, "plain text", *out)
}

func TestHttpNoBodyResp(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()

	client := Client{URL: ts.URL, Debug: true}

	api := NewBaseAPI(http.MethodGet, "/", nil, new([]byte), nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, api.StatusCode())
	buffer := *api.ResponseObject().(*[]byte)
	assert.Equal(t, []uint8([]byte(nil)), buffer)
}

type ErrStruct struct {
	ErrID   string `json:"error_id"`
	ErrCode string `json:"error_code,omitempty"`
	ErrText string `json:"error_text,omitempty"`
}

func TestHttpErrorObjectBack(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error_id":"resource.not_found","error_text":"Resource '/api/tm/3.8/config/active/monitor' does not exist"}`))
	}))
	defer ts.Close()

	client := Client{URL: ts.URL, Debug: true}

	api := NewBaseAPI(http.MethodGet, "/", nil, new([]byte), new(ErrStruct))
	err := client.Do(api)
	if err != nil {
		assert.Equal(t, http.StatusInternalServerError, api.StatusCode())
		errStruct := api.ErrorObject().(*ErrStruct)
		assert.Equal(t, "resource.not_found", errStruct.ErrID)
	}
}

type ReqBody struct {
	FieldOne string `json:"field_1"`
	FieldTwo string `json:"field_2"`
}

func TestHttpRequestPayload(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "PUT", r.Method)
		r.ParseForm()

		var reqPayload ReqBody
		decBody := json.NewDecoder(r.Body)
		decBody.Decode(&reqPayload)
		//assert.Equal(t, nil, err)
		assert.Equal(t, "Foo", reqPayload.FieldOne)
		assert.Equal(t, "Bar", reqPayload.FieldTwo)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"field_1":"Fooo", "field_2":"Baar"}`))
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	client := Client{
		URL:     ts.URL,
		Headers: headers,
		Timeout: 10,
	}

	var reqPayload ReqBody
	reqPayload.FieldOne = "Foo"
	reqPayload.FieldTwo = "Bar"
	assert.Equal(t, "Foo", reqPayload.FieldOne)
	assert.Equal(t, "Bar", reqPayload.FieldTwo)

	api := NewBaseAPI(http.MethodPut, "/", reqPayload, new(ReqBody), new(ErrStruct))
	err := client.Do(api)
	if err != nil {
		assert.Equal(t, http.StatusInternalServerError, api.StatusCode())
		errStruct := api.ErrorObject().(*ErrStruct)
		assert.Equal(t, "001", errStruct.ErrID)
		assert.Equal(t, "12345", errStruct.ErrCode)
	}

	respBody := *api.ResponseObject().(*ReqBody)
	assert.Equal(t, "Fooo", respBody.FieldOne)
	assert.Equal(t, "Baar", respBody.FieldTwo)
}
