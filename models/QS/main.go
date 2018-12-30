package qs

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/aws/aws-lambda-go/events"
)

var Unsupported = errors.New("Unsupported body structure")

func GetBody(req interface{}, data interface{}) error {
	switch v := req.(type) {
	case Request:
		err := v.ReadTo(data)
		if err != nil {
			return err
		}
	case map[string]interface{}:
		err := mapstructure.Decode(v, data)
		if err != nil {
			return err
		}
	default:
		return Unsupported
	}
	return nil
}

type Response events.APIGatewayProxyResponse

func NewResponse(status int, body interface{}) (Response, error) {
	j, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		return Response{}, err
	}
	headers := map[string]string{"Content-Type": "application/json"}

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
