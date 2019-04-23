package apigw_test

import (
	"net/http"
	"testing"

	"github.com/Ryanair/goaws/lambda/apigw"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	calls []*apigw.Request
	*apigw.Response
	error
}

func (h *mockHandler) Handle(req *apigw.Request) (*apigw.Response, error) {
	h.calls = append(h.calls, req)
	return h.Response, h.error
}

func TestNewHandler_postRequest(t *testing.T) {
	// given
	request := &events.APIGatewayProxyRequest{
		Resource:   "/{proxy+}",
		HTTPMethod: "POST",
		Body:       "some random request body",
		Headers: map[string]string{
			"Accept": "json",
		},
	}
	mockHandler := &mockHandler{
		Response: apigw.NewResponse(apigw.StatusCreated, "API response message"),
		error:    nil,
	}
	wrappedHandler := apigw.WrapHandler(mockHandler, apigw.Body(), apigw.Headers())

	// when
	resp, err := wrappedHandler(request)

	// then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "API response message",
	}, resp)
	assert.Equal(t, &apigw.Request{
		Resource: "/{proxy+}",
		Method:   "POST",
		Headers: map[string]string{
			"Accept": "json",
		},
		Body: "some random request body",
	}, mockHandler.calls[0])
}

func TestNewHandler_getRequest(t *testing.T) {
	// given
	request := &events.APIGatewayProxyRequest{
		Resource:   "/{proxy+}",
		HTTPMethod: "GET",
		PathParameters: map[string]string{
			"userId": "randomUserId",
		},
		QueryStringParameters: map[string]string{
			"firstQueryParam": "firstQueryParamValue",
		},
		Headers: map[string]string{
			"Accept": "json",
		},
	}
	mockHandler := &mockHandler{
		Response: apigw.NewResponse(apigw.StatusOK, "API response message"),
		error:    nil,
	}
	wrappedHandler := apigw.WrapHandler(mockHandler, apigw.PathParams(), apigw.QueryParams(), apigw.Headers())

	// when
	resp, err := wrappedHandler(request)

	// then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "API response message",
	}, resp)
	assert.Equal(t, &apigw.Request{
		Resource: "/{proxy+}",
		Method:   "GET",
		PathParameters: map[string]string{
			"userId": "randomUserId",
		},
		QueryStringParameters: map[string]string{
			"firstQueryParam": "firstQueryParamValue",
		},
		Headers: map[string]string{
			"Accept": "json",
		},
	}, mockHandler.calls[0])
}

func TestNewHandler_getRequest_optionsNotSpecified(t *testing.T) {
	// given
	request := &events.APIGatewayProxyRequest{
		Resource:   "/{proxy+}",
		HTTPMethod: "GET",
		PathParameters: map[string]string{
			"userId": "randomUserId",
		},
		QueryStringParameters: map[string]string{
			"firstQueryParam": "firstQueryParamValue",
		},
		Headers: map[string]string{
			"Accept": "json",
		},
	}
	mockHandler := &mockHandler{
		Response: apigw.NewResponse(apigw.StatusOK, "API response message"),
		error:    nil,
	}
	wrappedHandler := apigw.WrapHandler(mockHandler)

	// when
	resp, err := wrappedHandler(request)

	// then
	assert.Nil(t, err)
	assert.Equal(t, events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "API response message",
	}, resp)
	assert.Equal(t, &apigw.Request{
		Resource: "/{proxy+}",
		Method:   "GET",
		PathParameters: map[string]string{
			"userId": "randomUserId",
		},
		QueryStringParameters: map[string]string{
			"firstQueryParam": "firstQueryParamValue",
		},
		Headers: map[string]string{
			"Accept": "json",
		},
	}, mockHandler.calls[0])
}
