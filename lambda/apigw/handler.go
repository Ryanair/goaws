package apigw

import "github.com/aws/aws-lambda-go/events"

type handler interface {
	Handle(*Request) (*Response, error)
}

type EventConverter func(*events.APIGatewayProxyRequest) func(*Request)

type LambdaHandler func(*events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func WrapHandler(handler handler, eventConverters ...EventConverter) LambdaHandler {
	if len(eventConverters) == 0 {
		eventConverters = []EventConverter{Headers(), MultiValueHeaders(), QueryParams(), MultiValueQueryParams(), PathParams(), IsBase64Encoded(), Body()}
	}

	return func(event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var opts []func(*Request)
		for _, ec := range eventConverters {
			opts = append(opts, ec(event))
		}

		req := NewRequest(event.Resource, event.HTTPMethod, opts...)
		res, err := handler.Handle(req)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		return res.Convert(), nil
	}
}

func Headers() EventConverter {
	return func(event *events.APIGatewayProxyRequest) func(*Request) {
		return RequestHeaders(event.Headers)
	}
}

func MultiValueHeaders() EventConverter {
	return func(event *events.APIGatewayProxyRequest) func(*Request) {
		return RequestMultiValueHeaders(event.MultiValueHeaders)
	}
}

func QueryParams() EventConverter {
	return func(event *events.APIGatewayProxyRequest) func(*Request) {
		return RequestQueryParams(event.QueryStringParameters)
	}
}

func MultiValueQueryParams() EventConverter {
	return func(event *events.APIGatewayProxyRequest) func(*Request) {
		return RequestMultiValueQueryParams(event.MultiValueQueryStringParameters)
	}
}

func PathParams() EventConverter {
	return func(event *events.APIGatewayProxyRequest) func(*Request) {
		return RequestPathParams(event.PathParameters)
	}
}

func IsBase64Encoded() EventConverter {
	return func(event *events.APIGatewayProxyRequest) func(*Request) {
		return RequestIsBase64Encoded(event.IsBase64Encoded)
	}
}

func Body() EventConverter {
	return func(event *events.APIGatewayProxyRequest) func(*Request) {
		return RequestBody(event.Body)
	}
}
