// +build local ci

package apigw_test

import (
	"testing"

	"github.com/Ryanair/goaws/lambda/apigw"

	"github.com/stretchr/testify/assert"
)

func TestResponses(t *testing.T) {
	var (
		ok        = "ok"
		created   = "created"
		accepted  = "accepted"
		badReq    = "bad_request"
		notFound  = "not_found"
		conflict  = "conflict"
		intSrvErr = "internal_server_error"
	)
	hdrs := func(key, value string) func(*apigw.Response) {
		headers := map[string]string{key: value}
		return apigw.ResponseHeaders(headers)
	}

	multiHdrs := func(key, value string) func(*apigw.Response) {
		headers := map[string][]string{key: {value, value}}
		return apigw.ResponseMultiValueHeaders(headers)

	}
	testdata := []struct {
		status     int
		body       string
		kv         string
		headers    func(*apigw.Response)
		multiValue func(*apigw.Response)
	}{
		{status: apigw.StatusOK, body: ok, kv: ok, headers: hdrs(ok, ok), multiValue: multiHdrs(ok, ok)},
		{status: apigw.StatusCreated, body: created, kv: created, headers: hdrs(created, created), multiValue: multiHdrs(created, created)},
		{status: apigw.StatusAccepted, body: accepted, kv: accepted, headers: hdrs(accepted, accepted), multiValue: multiHdrs(accepted, accepted)},
		{status: apigw.StatusBadRequest, body: badReq, kv: badReq, headers: hdrs(badReq, badReq), multiValue: multiHdrs(badReq, badReq)},
		{status: apigw.StatusNotFound, body: notFound, kv: notFound, headers: hdrs(notFound, notFound), multiValue: multiHdrs(notFound, notFound)},
		{status: apigw.StatusConflict, body: conflict, kv: conflict, headers: hdrs(conflict, conflict), multiValue: multiHdrs(conflict, conflict)},
		{status: apigw.StatusInternalServerError, body: intSrvErr, kv: intSrvErr, headers: hdrs(intSrvErr, intSrvErr), multiValue: multiHdrs(intSrvErr, intSrvErr)},
	}

	for _, data := range testdata {

		response := apigw.NewResponse(data.status, data.body, data.headers, apigw.ResponseIsBase64Encoded(true), data.multiValue)

		assert.Equal(t, data.status, response.StatusCode)
		assert.Equal(t, data.kv, response.Headers[data.kv])
		assert.True(t, response.IsBase64Encoded)
		multi := map[string][]string{data.kv: {data.kv, data.kv}}
		assert.Equal(t, multi, response.MultiValueHeaders)

		converted := response.Convert()

		assert.Equal(t, response.StatusCode, converted.StatusCode)
		assert.Equal(t, response.Body, converted.Body)
		assert.Equal(t, response.Headers, converted.Headers)
		assert.Equal(t, response.IsBase64Encoded, converted.IsBase64Encoded)
	}
}
