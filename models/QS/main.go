package qs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/aws/aws-lambda-go/events"
)

var Unsupported = errors.New("Unsupported body structure")

func GetBody(req interface{}, data interface{}) error {
	r := Request{}
	err := mapstructure.Decode(req, &r)
	if err != nil {
		fmt.Println(err)
	}
	if !reflect.DeepEqual(r, Request{}) {
		err := r.ReadTo(data)
		return err
	}
	err = mapstructure.Decode(req, data)
	if err != nil {
		return err
	}
	return nil
}

func equal(v interface{}, p interface{}) bool {
	vt := reflect.TypeOf(v)
	pt := reflect.TypeOf(p)
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	if reflect.ValueOf(p).Kind() == reflect.Ptr {
		pt = pt.Elem()
	}
	nv := reflect.New(vt)
	np := reflect.New(pt)
	fmt.Println(nv, np)
	return true
}

type Response events.APIGatewayProxyResponse

func NewResponse(status int, body interface{}) (Response, error) {
	j, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		return Response{}, err
	}
	headers := map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "http://localhost:4200",
		"Access-Control-Allow-Credentials": "true",
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
		Body: string(j),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "http://localhost:4200",
			"Access-Control-Allow-Credentials": "true",
		},
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
