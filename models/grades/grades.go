package grades

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jinzhu/now"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Grade struct {
	Timestamp int64  `json:"timestamp"`
	Value     int    `json:"value"`
	Subject   string `json:"subject"`
}

func InputGrade(username string, timestamp int64, value int, subject string, conn *dynamodb.DynamoDB) error {
	j := fmt.Sprintf(`{"username": "%v","timestamp": %v, "value":%v, "subject":"%v"}`, username, timestamp, value, subject)
	body, err := database.MarshalJSONToDynamoMap(j)
	if err != nil {
		log.Println("line 28 error with marshal json to map")
		return errors.MarshalJsonToMapError
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("plannerix-grades"),
		Item:      body,
	}
	_, err = conn.PutItem(input)
	if err != nil {
		log.Println("line 37 error with put item")
		return errors.PutItemError
	}
	return nil
}

func GetAllGrades(username string, conn *dynamodb.DynamoDB) ([]Grade, error) {
	getItemInput := &dynamodb.QueryInput{
		TableName:              aws.String("plannerix-grades"),
		KeyConditionExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
	}
	output, err := conn.Query(getItemInput)
	if err != nil {
		log.Println("line 55 error with output")
		return nil, errors.OutputError
	}
	grades := []Grade{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &grades)
	if err != nil {
		log.Println("line 61 error with output")
		return nil, errors.UnmarshalListOfMapsError
	}
	return grades, nil
}

func GetWeeklyGrades(username string, conn *dynamodb.DynamoDB) ([]Grade, error) {
	// offSet, _ := time.ParseDuration("+02.00h")
	mondayTime := now.Monday().Unix()
	// log.Println("test", test)
	// mondayTime := now.BeginningOfWeek().UTC().AddDate(0, 0, -6).Unix()
	// log.Println(time.Unix(mondayTime, 0))
	// log.Println(mondayTime)

	filt := expression.Name("timestamp").GreaterThan(expression.Value(mondayTime)).
		And(expression.Name("username").Equal(expression.Value(username)))

	proj := expression.NamesList(expression.Name("timestamp"), expression.Name("value"), expression.Name("subject"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("line 78 couldn't build expression")
		return nil, errors.ExpressionBuilderError
	}

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("plannerix-grades"),
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

	weeklyGrades := []Grade{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &weeklyGrades)
	if err != nil {
		log.Println("line 99 error with unmarshal")
		return nil, errors.UnmarshalListOfMapsError
	}

	return weeklyGrades, nil
}

func GetYearGrades(username string, conn *dynamodb.DynamoDB) ([]Grade, error) {
	t := time.Now()
	var startDate int64
	var endDate int64
	if t.Month() >= 2 && t.Month() <= 8 {
		startDate = time.Date(t.Year(), 2, 1, 0, 0, 0, 0, time.UTC).Unix()
		endDate = time.Date(t.Year(), 9, 14, 23, 59, 59, 0, time.UTC).Unix()
	} else if t.Month() == 1 {
		startDate = time.Date(t.Year()-1, 9, 15, 0, 0, 0, 0, time.UTC).Unix()
		endDate = time.Date(t.Year(), 1, 31, 23, 59, 59, 0, time.UTC).Unix()
	} else {
		startDate = time.Date(t.Year(), 9, 15, 0, 0, 0, 0, time.UTC).Unix()
		endDate = time.Date(t.Year()+1, 1, 31, 23, 59, 59, 0, time.UTC).Unix()
	}

	filt := expression.Name("timestamp").Between(expression.Value(startDate), expression.Value(endDate)).
		And(expression.Name("username").Equal(expression.Value(username)))

	proj := expression.NamesList(expression.Name("timestamp"), expression.Name("value"), expression.Name("subject"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("line 128 couldn't build expression")
		return nil, errors.ExpressionBuilderError
	}

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("plannerix-grades"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	output, err := conn.Scan(getItemScanInput)
	if err != nil {
		log.Println("line 143 error with output")
		return nil, errors.OutputError
	}
	yearGrades := []Grade{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &yearGrades)
	if err != nil {
		log.Println("line 149 error with unmarshal")
		return nil, errors.UnmarshalListOfMapsError
	}
	return yearGrades, nil
}

func DeleteGrade(username string, timestamp int64, conn *dynamodb.DynamoDB) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("plannerix-grades"),
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
		log.Println("line 169 error with unmarshal")
		return errors.DeleteItemError
	}
	return nil
}

func AdaptTimestamp(timestamp int64) int64 {
	// offSet, _ := time.ParseDuration("+02.00h")
	// .Add(offSet)
	now := time.Now()
	gradeDate := time.Unix(timestamp, 0)
	// log.Println(now)
	// log.Println(gradeDate)
	if gradeDate.Year() < now.Year() {
		now = now.AddDate(gradeDate.Year()-now.Year(), 0, 0)
	}
	diff := gradeDate.YearDay() - now.YearDay()
	// log.Println(diff)
	newTimestamp := now.AddDate(0, 0, diff).Unix()
	return newTimestamp
}
