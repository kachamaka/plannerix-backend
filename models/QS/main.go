package qs

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

func NewResponse(status int, headers map[string]string, body interface{}) (Response, error) {
	j, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		return Response{}, err
	}
	if headers == nil {
		headers = map[string]string{
			"content-type": "application/json",
		}
	}
	return Response{
		Body:            string(j),
		Headers:         headers,
		IsBase64Encoded: false,
		StatusCode:      status,
	}, nil
}

func NewError(errMsg string, code int) (Response, error) {
	body := map[string]interface{}{
		"success": false,
		"errMsg":  errMsg,
		"code":    code,
	}
	j, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		log.Println("Error by encoding error message:", err)
		return Response{}, err
	}
	return Response{
		Body:            string(j),
		Headers:         map[string]string{"content-type": "application/json"},
		StatusCode:      200,
		IsBase64Encoded: false,
	}, nil
}

type Request events.APIGatewayProxyRequest

func (r Request) ReadTo(b interface{}) error {
	err := json.Unmarshal([]byte(r.Body), b)
	if err != nil {
		return err
	}
	return nil
}
