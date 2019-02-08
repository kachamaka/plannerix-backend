package events

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/s-org-backend/models/database"
)

type Event struct {
	Subject     string `json:"subject"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Timestamp   int64  `json:"time"`
}

func CreateEvent(username string, subject string, title string, description string, timestamp int64, conn *dynamodb.DynamoDB) error {
	j := fmt.Sprintf(`{"username": "%v","timestamp": %v, "subject":"%v", "title":"%v", "description":"%v"}`, username, timestamp, subject, title, description)
	body, err := database.MarshalJSONToDynamoMap(j)
	if err != nil {
		log.Println("line 27 error with marshal json to map")
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("s-org-events"),
		Item:      body,
	}
	_, err = conn.PutItem(input)
	if err != nil {
		log.Println("line 36 error with put item")
		return err
	}
	return nil
}

func GetAllEvents(username string, conn *dynamodb.DynamoDB) ([]Event, error) {

	filt := expression.Name("username").Equal(expression.Value(username))

	proj := expression.NamesList(expression.Name("timestamp"), expression.Name("subject"), expression.Name("title"), expression.Name("description"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("line 51 error with building expression")
		return nil, err
	}

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("s-org-events"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	output, err := conn.Scan(getItemScanInput)
	if err != nil {
		log.Println("line 65 error with output")
		return []Event{}, nil
	}

	events := []Event{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &events)
	if err != nil {
		log.Println("line 72 error with unmarshal")
		return []Event{}, nil
	}

	return events, nil
}

func EditEvent(username string, subject string, title string, description string, timestamp int64, conn *dynamodb.DynamoDB) error {
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("s-org-events"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
			"timestamp": {
				N: aws.String(strconv.FormatInt(timestamp, 10)),
			},
		},
		UpdateExpression: aws.String("set subject = :subject, title = :title, description = :description"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":subject": {
				S: aws.String(subject),
			},
			":title": {
				S: aws.String(title),
			},
			":description": {
				S: aws.String(description),
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
	}

	_, err := conn.UpdateItem(updateItemInput)
	if err != nil {
		log.Println("line 107 error with update item")
		return err
	}
	return nil
}

func DeleteEvent(username string, timestamp int64, conn *dynamodb.DynamoDB) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("s-org-events"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
			"eventTime": {
				N: aws.String(strconv.FormatInt(timestamp, 10)),
			},
		},
	}
	_, err := conn.DeleteItem(deleteItemInput)
	if err != nil {
		log.Println("line 127 error with delete item", err)
		return err
	}
	return nil
}
