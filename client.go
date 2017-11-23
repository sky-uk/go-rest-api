package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/sky-uk/go-rest-api/contenttype"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Client struct.
type Client struct {
	URL        string
	User       string
	Password   string
	IgnoreSSL  bool
	Debug      bool
	Headers    map[string]string
	Timeout    time.Duration // in seconds
	StatusCode int
}

func (restClient *Client) formatRequestPayload(api *BaseAPI) (io.Reader, error) {

	var requestPayload io.Reader

	var reqBytes []byte
	if api.RequestObject() != nil {
		var err error
		contentType := contenttype.GetType(restClient.Headers["Content-Type"])

		switch contentType {

		case "json":
			reqBytes, err = json.Marshal(api.RequestObject())
			if err != nil {
				log.Fatal("[ERROR] ", err)
				return nil, err
			}

		case "xml":
			reqBytes, err = xml.Marshal(api.RequestObject())
			if err != nil {
				log.Fatal("[ERROR] ", err)
				return nil, err
			}

		case "octet-stream", "plain", "html":
			reqBytes = api.RequestObject().([]byte)

		}
		requestPayload = bytes.NewReader(reqBytes)
	}

	if restClient.Debug {
		log.Println("[TRACE] --------------------------------------------------------------")
		log.Println("[TRACE] Request payload:")
		log.Println("[TRACE] ", string(reqBytes))
		log.Println("[TRACE] --------------------------------------------------------------")
	}

	return requestPayload, nil
}

// Do - makes the API call.
func (restClient *Client) Do(api *BaseAPI) error {

	requestURL := fmt.Sprintf("%s%s", restClient.URL, api.Endpoint())
	if restClient.Debug {
		log.Printf("[TRACE] Going to perform request:[%s] %s\n", api.Method(), requestURL)
	}

	if restClient.Headers == nil {
		restClient.Headers = make(map[string]string)
	}

	_, ok := restClient.Headers["Content-Type"]
	if !ok {
		restClient.Headers["Content-Type"] = "text/plain"
	}

	requestPayload, err := restClient.formatRequestPayload(api)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(api.Method(), requestURL, requestPayload)
	if err != nil {
		log.Println("[ERROR] Error building the request: ", err)
		return err
	}

	if restClient.User != "" {
		req.SetBasicAuth(restClient.User, restClient.Password)
	}

	for headerKey, headerValue := range restClient.Headers {
		req.Header.Set(headerKey, headerValue)
	}

	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: restClient.IgnoreSSL},
		MaxIdleConns:      10,
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: true,
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   restClient.Timeout * time.Second,
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println("[ERROR] Error executing request: ", err)
		return err
	}
	defer res.Body.Close()
	restClient.StatusCode = res.StatusCode
	return restClient.handleResponse(api, res)
}

func (restClient *Client) handleResponse(apiObj *BaseAPI, res *http.Response) error {

	apiObj.SetStatusCode(res.StatusCode)
	bodyText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERROR] Error reading response: ", err)
		return err
	}

	if len(bodyText) > 0 {
		contentType := contenttype.GetType(res.Header.Get("Content-Type"))
		if restClient.Debug {
		}

		if restClient.Debug {
			log.Println("[TRACE] --------------------------------------------------------------")
			log.Println("[TRACE] Response content type: ", contentType)
			log.Println("[TRACE] Response payload:")
			log.Println("[TRACE] ", string(bodyText))
			log.Println("[TRACE] --------------------------------------------------------------")
		}
		apiObj.SetRawResponse(bodyText)

		switch contentType {
		case "json":
			if apiObj.StatusCode() >= http.StatusOK && apiObj.StatusCode() < http.StatusBadRequest {
				err := json.Unmarshal(bodyText, apiObj.ResponseObject())
				if err != nil {
					log.Println("[ERROR] Error unmarshalling response: ", err)
					return err
				}
			} else {
				if apiObj.ErrorObject() != nil {
					err := json.Unmarshal(bodyText, apiObj.ErrorObject())
					if err != nil {
						log.Printf("[ERROR] Error unmarshalling error response:\n%v", err)
						return err
					}
				}
				errMsg := fmt.Sprintf("Response status code: %d", apiObj.StatusCode())
				return errors.New(errMsg)
			}

		case "xml":
			if apiObj.StatusCode() >= http.StatusOK && apiObj.StatusCode() < http.StatusBadRequest {
				err := xml.Unmarshal(bodyText, apiObj.ResponseObject())
				if err != nil {
					log.Println("[ERROR] Error unmarshalling response: ", err)
					return err
				}
			} else {
				if apiObj.ErrorObject() != nil {
					err := xml.Unmarshal(bodyText, apiObj.ErrorObject())
					if err != nil {
						log.Printf("[ERROR] Error unmarshalling error response:\n%v", err)
					}
				}
				errMsg := fmt.Sprintf("Response status code: %d", apiObj.StatusCode())
				return errors.New(errMsg)
			}

		case "octet-stream":
			if apiObj.ResponseObject() != nil {
				if pstream, is := apiObj.ResponseObject().(*[]byte); is {
					*pstream = bodyText
				} else {
					log.Println("[WARN] Response object expected to be *[]byte")
				}
			}

		case "plain", "html":
			if apiObj.ResponseObject() != nil {
				if pstream, is := apiObj.ResponseObject().(*string); is {
					*pstream = string(bodyText)
				} else {
					log.Println("[WARN] Response object expected to be *string")
				}
			}

		default:
			log.Printf("[WARN] Content type %s not supported yet", contentType)
		}
	} else {
	}

	return nil
}
