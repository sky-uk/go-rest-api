package rest

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/sky-uk/go-rest-api/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var user = "nsxUser"
var password = "nsxPass"
var ignoreSSL = true
var debug = true
var client *Client

var server *httptest.Server

const (
	unauthorizedStatusCode = http.StatusForbidden
	unauthorizedResponse   = "anauthorized"
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
	client = NewClient(server.URL, user, password, ignoreSSL, debug, headers)
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
	client = NewClient(server.URL, user, password, ignoreSSL, debug, headers)
}

func TestHappyCase(t *testing.T) {
	setup(200, "pong")
	client = NewClient(server.URL, user, password, ignoreSSL, debug, nil)
	apiRequest := api.NewRestAPI(http.MethodGet, "/", nil, nil, nil)

	err := client.Do(apiRequest)

	assert.Nil(t, err)
}

// TODO: add TestFailWhenNotValidSSLCerts(t *testing.T)

func TestBasicAuthFailure(t *testing.T) {
	setup(0, "")
	client = NewClient(server.URL, "invalidUser", "invalidPass", ignoreSSL, debug, nil)

	apiRequest := api.NewRestAPI(http.MethodGet, "/", nil, nil, nil)
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

	client := NewClient(ts.URL, "", "", false, true, nil)
	api := api.NewRestAPI(http.MethodGet, "/", nil, new(string), nil)
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

func TestHttpJSONReq(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"fields":{"foo":"bar"}}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "", "", false, true, nil)

	api := api.NewRestAPI(http.MethodGet, "/", nil, new(JSONFoo), nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, api.StatusCode())
	resp := *api.ResponseObject().(*JSONFoo)
	assert.Equal(t, "bar", resp.Fields["foo"])
}

func TestHttpXMLReq(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><fields><foo>bar</foo></fields>`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "", "", false, true, nil)

	api := api.NewRestAPI(http.MethodGet, "/", nil, new(XMLFoo), nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, api.StatusCode())
	resp := *api.ResponseObject().(*XMLFoo)
	assert.Equal(t, "bar", resp.Foo)
}

func TestHttpOctetStreamReq(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte(`ouitybw50ybvqy9yt8b6983p8v93`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "", "", false, true, nil)

	api := api.NewRestAPI(http.MethodGet, "/", nil, new([]byte), nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, api.StatusCode())
	buffer := *api.ResponseObject().(*[]byte)
	assert.Equal(t, []byte("ouitybw50ybvqy9yt8b6983p8v93"), buffer)
}
