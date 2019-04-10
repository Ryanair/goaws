package apigw_test

import (
	"fmt"
	"testing"

	"github.com/ryanair/go-aws/lambda/apigw"
	"github.com/stretchr/testify/assert"
)

func TestRequests(t *testing.T) {
	var (
		get  = "get"
		post = "post"
		put  = "put"
		del  = "del"
	)
	hdrs := func(key, value string) func(*apigw.Request) {
		headers := map[string]string{key: value}
		return apigw.RequestHeaders(headers)
	}

	multiHdrs := func(key, value string) func(*apigw.Request) {
		headers := map[string][]string{key: {value, value}}
		return apigw.RequestMultiValueHeaders(headers)

	}

	testdata := []struct {
		method string
		kv     string
		body   string
		base64 bool
	}{
		{method: apigw.GetMethod, kv: get, body: "", base64: true},
		{method: apigw.PostMethod, kv: post, body: post, base64: false},
		{method: apigw.PutMethod, kv: put, body: put, base64: false},
		{method: apigw.DeleteMethod, kv: del, body: "", base64: false},
	}

	for _, data := range testdata {
		resource := fmt.Sprintf("/%v", data.kv)
		base64 := apigw.RequestIsBase64Encoded(data.base64)
		hdrs := hdrs(data.kv, data.kv)
		multiHdrs := multiHdrs(data.kv, data.kv)
		body := apigw.RequestBody(data.body)
		pathParams := apigw.RequestPathParams(map[string]string{data.kv: data.kv})
		multiPathParams := apigw.RequestMultiValueQueryParams(map[string][]string{data.kv: {data.kv, data.kv}})
		queryParams := apigw.RequestQueryParams(map[string]string{data.kv: data.kv})

		request := apigw.NewRequest(resource, data.method, body, hdrs, multiHdrs, base64, pathParams, multiPathParams, queryParams)

		assert.Equal(t, resource, request.Resource)
		assert.Equal(t, data.method, request.Method)
		assert.Equal(t, data.body, request.Body)
		assert.Equal(t, data.base64, request.IsBase64Encoded)
		assert.Equal(t, data.kv, request.PathParameters[data.kv])
		assert.Equal(t, data.kv, request.QueryStringParameters[data.kv])
		assert.Equal(t, []string{data.kv, data.kv}, request.MultiValueQueryStringParameters[data.kv])
	}
}
