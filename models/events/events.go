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
	Type        int    `json:"subjectType"`
	Description string `json:"description"`
	Timestamp   int64  `json:"timestamp"`
}

func CreateEvent(username string, subject string, subjectType int, description string, timestamp int64, conn *dynamodb.DynamoDB) error {
	j := fmt.Sprintf(`{"username": "%v","timestamp": %v, "subject":"%v", "subjectType":%v, "description":"%v"}`, username, timestamp, subject, subjectType, description)
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

	proj := expression.NamesList(expression.Name("timestamp"), expression.Name("subject"), expression.Name("subjectType"), expression.Name("description"))

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
		log.Println(err)
		return nil, errors.UnmarshalListOfMapsError
	}

	return events, nil
}

func GetWeeklyEvents(username string, conn *dynamodb.DynamoDB) ([]Event, error) {
	location, _ := time.LoadLocation("Europe/Sofia")
	now := time.Now().In(location)
	nowTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	// nowAfterTwoWeeks := time.Now().AddDate(0, 0, 7).Unix()
	nowAfterTwoWeeksTime := time.Date(now.Year(), now.Month(), now.Day()+7, 23, 59, 59, 0, location)

	filt := expression.Name("timestamp").GreaterThanEqual(expression.Value(nowTime.Unix())).
		And(expression.Name("timestamp").LessThan(expression.Value(nowAfterTwoWeeksTime.Unix()))).
		And(expression.Name("username").Equal(expression.Value(username)))

	proj := expression.NamesList(expression.Name("timestamp"), expression.Name("subject"), expression.Name("subjectType"), expression.Name("description"))

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

func EditEvent(username string, subject string, subjectType int, description string, timestamp int64, conn *dynamodb.DynamoDB) error {
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
		UpdateExpression: aws.String("set subject = :subject, subjectType = :subjectType, description = :description"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":subject": {
				S: aws.String(subject),
			},
			":subjectType": {
				N: aws.String(strconv.Itoa(subjectType)),
			},
			":description": {
				S: aws.String(description),
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
	}

	output, err := conn.UpdateItem(updateItemInput)
	log.Println(output, "output")
	if err != nil {
		log.Println(err)
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
			"timestamp": {
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

func AdaptTimestamp(timestamp int64) int64 {
	offSet, _ := time.ParseDuration("+02.00h")
	now := time.Now().Add(offSet)
	testDate := time.Unix(timestamp, 0).UTC().Add(offSet)
	// log.Println(now)
	// log.Println(testDate)
	diff := testDate.YearDay() - now.YearDay()
	// log.Println(diff)
	newTimestamp := now.AddDate(0, 0, diff).Unix()
	return newTimestamp
}

// func AdaptTimestamp(timestamp int64) int64 {
// 	randInts := random.RangeInt(0, 86399000000000, 1)
// 	randInt := randInts[0]
// 	newTimestamp := int64(randInt) + timestamp
// 	return newTimestamp
// }
