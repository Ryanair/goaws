package apigw

const (
	GetMethod    = "GET"
	PostMethod   = "POST"
	PutMethod    = "PUT"
	DeleteMethod = "DELETE"
)

type Request struct {
	Resource                        string
	Path                            string
	Method                          string
	Headers                         map[string]string
	MultiValueHeaders               map[string][]string
	QueryStringParameters           map[string]string
	MultiValueQueryStringParameters map[string][]string
	PathParameters                  map[string]string
	StageVariables                  map[string]string
	Body                            string
	IsBase64Encoded                 bool
}

func NewRequest(resource, method string, options ...func(*Request)) *Request {
	req := Request{
		Resource: resource,
		Method:   method,
	}
	for _, option := range options {
		option(&req)
	}
	return &req
}

func RequestHeaders(headers map[string]string) func(*Request) {
	return func(request *Request) {
		request.Headers = headers
	}
}

func RequestMultiValueHeaders(multiValueHeaders map[string][]string) func(*Request) {
	return func(request *Request) {
		request.MultiValueHeaders = multiValueHeaders
	}
}

func RequestQueryParams(params map[string]string) func(*Request) {
	return func(request *Request) {
		request.QueryStringParameters = params
	}
}

func RequestMultiValueQueryParams(params map[string][]string) func(*Request) {
	return func(request *Request) {
		request.MultiValueQueryStringParameters = params
	}
}

func RequestPathParams(params map[string]string) func(*Request) {
	return func(request *Request) {
		request.PathParameters = params
	}
}

func RequestBody(body string) func(*Request) {
	return func(request *Request) {
		request.Body = body
	}
}

func RequestIsBase64Encoded(base64Encoded bool) func(*Request) {
	return func(request *Request) {
		request.IsBase64Encoded = base64Encoded
	}
}
