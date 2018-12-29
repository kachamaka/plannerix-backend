package database

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var conn *dynamodb.DynamoDB

//MarshalJSONToDynamoMap reads any json string (`{"hello":"world"}`) and parses it to a vaild `map[string]*dynamodb.AttributeValue` in hopes to make it simpler
func MarshalJSONToDynamoMap(j string) (map[string]*dynamodb.AttributeValue, error) {
	var js map[string]interface{}
	if json.Unmarshal([]byte(j), &js) != nil {
		return nil, errors.New("invalid json string body")
	}
	v, err := dynamodbattribute.MarshalMap(js)
	if err != nil {
		return nil, err
	}
	return v, nil
}

//SetConn accept a pointer to a pointer of a dynamodb.DynamoDB (**dynamodb.DynamoDB) checks if there is a connection. If there is none creates a new connection.
//Created for AWS Lambda to use a connection.
func SetConn(conn **dynamodb.DynamoDB) {
	//Credit Martin
	config := &aws.Config{
		Region: aws.String("eu-central-1"),
	}
	user := os.Getenv("USER")
	if user == "trayan" || user == "Lion" || user == "kachamaka-PC" {
		config.Endpoint = aws.String("http://localhost:8000")
	}
	if *conn == nil {
		s := session.Must(session.NewSession(config))
		svc := dynamodb.New(s)
		*conn = svc
	}
}
