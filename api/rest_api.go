package api

// RestAPI  - Base API struct.
type RestAPI struct {
	method         string
	endpoint       string
	requestObject  interface{}
	responseObject interface{}
	errorObject    interface{}
	statusCode     int
	rawResponse    []byte
	err            error
}

// NewRestAPI - Returns a new object of the RestAPI.
func NewRestAPI(
	method string,
	endpoint string,
	requestObject interface{},
	responseObject interface{},
	errorObject interface{},
) *RestAPI {
	return &RestAPI{method, endpoint, requestObject, responseObject, errorObject, 0, nil, nil}
}

// RequestObject - Returns the request object of the RestAPI
func (b *RestAPI) RequestObject() interface{} {
	return b.requestObject
}

// ResponseObject - Returns the ResponseObject interface.
func (b *RestAPI) ResponseObject() interface{} {
	return b.responseObject
}

// ErrorObject - Returns the ErrorObject interface.
func (b *RestAPI) ErrorObject() interface{} {
	return b.errorObject
}

// Method - Returns the Method string, i.e. GET, PUT, POST.
func (b *RestAPI) Method() string {
	return b.method
}

// Endpoint - Returns the Endpoint url string.
func (b *RestAPI) Endpoint() string {
	return b.endpoint
}

// StatusCode - Returns the status code of the api.
func (b *RestAPI) StatusCode() int {
	return b.statusCode
}

// RawResponse - Returns the rawResponse object as byte type.
func (b *RestAPI) RawResponse() []byte {
	return b.rawResponse
}

// Error - Returns the err the api.
func (b *RestAPI) Error() error {
	return b.err
}

// SetStatusCode - Sets the statusCode from api object.
func (b *RestAPI) SetStatusCode(statusCode int) {
	b.statusCode = statusCode
}

// SetRawResponse - Sets the rawResponse on api object.
func (b *RestAPI) SetRawResponse(rawResponse []byte) {
	b.rawResponse = rawResponse
}

// SetError - Sets the err on api object.
func (b *RestAPI) SetError(err error) {
	b.err = err
}

// SetResponseObject - Sets the responseObject
func (b *RestAPI) SetResponseObject(res interface{}) {
	b.responseObject = res
}

// SetErrorObject - Sets the errorObject
func (b *RestAPI) SetErrorObject(res interface{}) {
	b.errorObject = res
}
