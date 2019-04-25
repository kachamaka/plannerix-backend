package lambdat

//This package is used to test if the lambda function works
//In live environment the time is very different

import (
	"encoding/json"
	"net"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djhworld/go-lambda-invoke/golambdainvoke"

	"github.com/aws/aws-lambda-go/events"
)

func ReadBody(body interface{}) (events.APIGatewayProxyRequest, error) {
	result, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyRequest{}, err
	}
	req.Body = string(result)
	return req, nil
}

func InvokeHandler(handler interface{}, request interface{}) (string, error) {
	os.Setenv("_LAMBDA_SERVER_PORT", "9000")
	go lambda.Start(handler)
	for {
		conn, _ := net.DialTimeout("tcp", net.JoinHostPort("", "9000"), time.Millisecond*1)
		if conn != nil {
			conn.Close()
			break
		}
	}
	input := golambdainvoke.Input{
		Port:    9000,
		Payload: request,
	}
	res, err := golambdainvoke.Run(input)
	if err != nil && err.Error() != string([]byte{60, 110, 105, 108, 62, 10}) {
		return "", err
	}
	return string(res), nil
}

var req = events.APIGatewayProxyRequest{
	Resource:                        "/{proxy+}",
	Path:                            "/Seattle",
	HTTPMethod:                      "POST",
	Headers:                         map[string]string{"day": "Friday"},
	MultiValueHeaders:               map[string][]string{},
	QueryStringParameters:           map[string]string{"time": "morning"},
	MultiValueQueryStringParameters: map[string][]string{},
	PathParameters:                  map[string]string{"proxy": "Seattle"},
	StageVariables:                  map[string]string{},
	RequestContext: events.APIGatewayProxyRequestContext{
		AccountID:    "123456789012",
		ResourceID:   "nl9h80",
		Stage:        "test-invoke-stage",
		RequestID:    "test-invoke-request",
		ResourcePath: "/{proxy+}",
		Authorizer:   map[string]interface{}{},
		HTTPMethod:   "POST",
		APIID:        "r275xc9bmd",
	},
	Body:            "",
	IsBase64Encoded: false,
}
