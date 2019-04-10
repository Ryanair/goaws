package apigw

import (
	"github.com/aws/aws-lambda-go/events"
)

const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusAccepted            = 202
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusInternalServerError = 500
)

type Response struct {
	StatusCode        int
	Headers           map[string]string
	MultiValueHeaders map[string][]string
	Body              string
	IsBase64Encoded   bool
}

func NewResponse(statusCode int, body string, options ...func(*Response)) *Response {
	response := Response{
		StatusCode: statusCode,
		Body:       body,
	}
	for _, option := range options {
		option(&response)
	}
	return &response
}

func ResponseHeaders(headers map[string]string) func(*Response) {
	return func(response *Response) {
		response.Headers = headers
	}
}

func ResponseMultiValueHeaders(multiValueHeaders map[string][]string) func(*Response) {
	return func(response *Response) {
		response.MultiValueHeaders = multiValueHeaders
	}
}

func ResponseIsBase64Encoded(base64Encoded bool) func(*Response) {
	return func(response *Response) {
		response.IsBase64Encoded = base64Encoded
	}
}

func (api *Response) Convert() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        api.StatusCode,
		Headers:           api.Headers,
		MultiValueHeaders: api.MultiValueHeaders,
		Body:              api.Body,
		IsBase64Encoded:   api.IsBase64Encoded,
	}
}
