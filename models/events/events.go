package events

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
		return errors.MarshalJsonToMapError
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("s-org-events"),
		Item:      body,
	}
	_, err = conn.PutItem(input)
	if err != nil {
		log.Println("line 36 error with put item")
		return errors.PutItemError
	}
	return nil
}

func GetAllEvents(username string, conn *dynamodb.DynamoDB) ([]Event, error) {

	filt := expression.Name("username").Equal(expression.Value(username))

	proj := expression.NamesList(expression.Name("timestamp"), expression.Name("subject"), expression.Name("title"), expression.Name("description"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("line 51 error with building expression")
		return nil, errors.ExpressionBuilderError
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
		return []Event{}, errors.OutputError
	}

	events := []Event{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &events)
	if err != nil {
		log.Println("line 72 error with unmarshal")
		return nil, errors.UnmarshalListOfMapsError
	}

	return events, nil
}

func GetWeeklyEvents(username string, conn *dynamodb.DynamoDB) ([]Event, error) {
	now := time.Now().UnixNano()
	nowAfterTwoWeeks := time.Now().AddDate(0, 0, 14).UnixNano()

	filt := expression.Name("timestamp").GreaterThan(expression.Value(now)).
		And(expression.Name("timestamp").LessThan(expression.Value(nowAfterTwoWeeks))).
		And(expression.Name("username").Equal(expression.Value(username)))

	proj := expression.NamesList(expression.Name("timestamp"), expression.Name("subject"), expression.Name("title"), expression.Name("description"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("line 78 couldn't build expression")
		return nil, errors.ExpressionBuilderError
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
		log.Println("line 92 error with output")
		return nil, errors.OutputError
	}

	weeklyEvents := []Event{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &weeklyEvents)
	if err != nil {
		log.Println("line 99 error with unmarshal")
		return nil, errors.UnmarshalListOfMapsError
	}

	return weeklyEvents, nil
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
		return errors.UpdateItemError
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
		return errors.DeleteItemError
	}
	return nil
}
