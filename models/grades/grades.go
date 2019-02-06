package grades

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jinzhu/now"
	"gitlab.com/s-org-backend/models/database"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Grade struct {
	Timestamp int64  `json:"timestamp"`
	Value     int    `json:"value"`
	Subject   string `json:"subject"`
}

func InputGrade(user_id string, timestamp int64, value int, subject string, conn *dynamodb.DynamoDB) error {
	j := fmt.Sprintf(`{"user_id": "%v","timestamp": %v, "value":%v, "subject":"%v"}`, user_id, timestamp, value, subject)
	body, err := database.MarshalJSONToDynamoMap(j)
	if err != nil {
		log.Println("line 26")
		return err
	}
	input := &dynamodb.PutItemInput{
		TableName:           aws.String("s-org-grades"),
		ConditionExpression: aws.String("attribute_not_exists(user_id)"),
		Item:                body,
	}
	_, err = conn.PutItem(input)
	if err != nil {
		log.Println("line 36")
		return err
	}
	return nil
}

func GetAllGrades(user_id string, conn *dynamodb.DynamoDB) ([]Grade, error) {
	getItemInput := &dynamodb.QueryInput{
		TableName:              aws.String("s-org-grades"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user_id": {
				S: aws.String(user_id),
			},
		},
	}
	output, err := conn.Query(getItemInput)
	if err != nil {
		log.Println("line 54", err)
		return []Grade{}, err
	}
	grades := []Grade{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &grades)
	if err != nil {
		log.Println("line 60")
		return []Grade{}, nil
	}
	return grades, nil
}

func GetWeeklyGrades(user_id string, conn *dynamodb.DynamoDB) ([]Grade, error) {
	mondayTime := now.BeginningOfWeek().UnixNano()

	filt := expression.Name("timestamp").GreaterThan(expression.Value(mondayTime)).
		And(expression.Name("user_id").Equal(expression.Value(user_id)))

	proj := expression.NamesList(expression.Name("timestamp"), expression.Name("value"), expression.Name("subject"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("s-org-grades"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	output, err := conn.Scan(getItemScanInput)
	if err != nil {
		log.Println("line 68")
		return []Grade{}, nil
	}

	weeklyGrades := []Grade{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &weeklyGrades)
	if err != nil {
		log.Println("line 60")
		return []Grade{}, nil
	}

	return weeklyGrades, nil
}

func GetYearGrades(user_id string, conn *dynamodb.DynamoDB) ([]Grade, error) {
	t := time.Now()
	var startDate int64
	var endDate int64
	if t.Month() >= 2 && t.Month() <= 8 {
		startDate = time.Date(t.Year(), 2, 1, 0, 0, 0, 0, time.UTC).UnixNano()
		endDate = time.Date(t.Year(), 8, 31, 23, 59, 59, 0, time.UTC).UnixNano()
	} else if t.Month() == 1 {
		startDate = time.Date(t.Year()-1, 9, 1, 0, 0, 0, 0, time.UTC).UnixNano()
		endDate = time.Date(t.Year(), 1, 31, 23, 59, 59, 0, time.UTC).UnixNano()
	} else {
		startDate = time.Date(t.Year(), 9, 1, 0, 0, 0, 0, time.UTC).UnixNano()
		endDate = time.Date(t.Year()+1, 1, 31, 23, 59, 59, 0, time.UTC).UnixNano()
	}

	filt := expression.Name("timestamp").Between(expression.Value(startDate), expression.Value(endDate)).
		And(expression.Name("user_id").Equal(expression.Value(user_id)))

	proj := expression.NamesList(expression.Name("timestamp"), expression.Name("value"), expression.Name("subject"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("s-org-grades"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	output, err := conn.Scan(getItemScanInput)
	if err != nil {
		log.Println("line 117", err)
		return []Grade{}, err
	}
	yearGrades := []Grade{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &yearGrades)
	if err != nil {
		log.Println("line 123")
		return []Grade{}, nil
	}
	return yearGrades, nil
}

func DeleteGrade(user_id string, timestamp int64, conn *dynamodb.DynamoDB) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("s-org-grades"),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(user_id),
			},
			"timestamp": {
				N: aws.String(strconv.FormatInt(timestamp, 10)),
			},
		},
	}
	_, err := conn.DeleteItem(deleteItemInput)
	if err != nil {
		log.Println("line 143", err)
		return err
	}
	return nil
}
