package events

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Event struct {
	EventID     string `json:"event_id"`
	ID          string `json:"id"`
	SubjectID   string `json:"subject_id"`
	Type        int    `json:"subjectType"`
	Description string `json:"description"`
	Timestamp   int64  `json:"eventTime"`
}

func CreateEvent(id string, subject string, subjectType int, description string, timestamp int64, conn *dynamodb.DynamoDB) error {
	eventID := generateEventID()
	j := fmt.Sprintf(`{"event_id": "%v", "id": "%v","eventTime": %v, "subject_id":"%v", "subjectType":%v, "description":"%v"}`, eventID, id, timestamp, subject, subjectType, description)
	body, err := database.MarshalJSONToDynamoMap(j)
	if err != nil {
		log.Println("line 27 error with marshal json to map")
		return errors.MarshalJsonToMapError
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("plannerix-events"),
		Item:      body,
	}
	_, err = conn.PutItem(input)
	if err != nil {
		log.Println(err)
		log.Println("line 36 error with put item")
		return errors.PutItemError
	}
	return nil
}

func GetAllEvents(id string, conn *dynamodb.DynamoDB) ([]Event, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("plannerix-events"),
		IndexName: aws.String("queryIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(id)}},
			},
		},
	}

	res, err := conn.Query(queryInput)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	// log.Println(res, "test")

	events := []Event{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &events)

	if err != nil {
		log.Println("line 99 error with unmarshal")
		return nil, errors.UnmarshalListOfMapsError
	}

	return events, nil
}

func GetWeeklyEvents(id string, conn *dynamodb.DynamoDB) ([]Event, error) {

	location, _ := time.LoadLocation("Europe/Sofia")
	now := time.Now().In(location)
	nowTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	// nowAfterTwoWeeks := time.Now().AddDate(0, 0, 7).Unix()
	nowAfterTwoWeeksTime := time.Date(now.Year(), now.Month(), now.Day()+7, 23, 59, 59, 0, location)

	log.Println(id)
	log.Println(nowTime.Unix())
	log.Println(nowAfterTwoWeeksTime.Unix())
	filt := expression.Name("eventTime").GreaterThanEqual(expression.Value(nowTime.Unix())).
		And(expression.Name("eventTime").LessThan(expression.Value(nowAfterTwoWeeksTime.Unix()))).
		And(expression.Name("id").Equal(expression.Value(id)))

	proj := expression.NamesList(expression.Name("eventTime"), expression.Name("subject"), expression.Name("subjectType"), expression.Name("description"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("line 78 couldn't build expression")
		return nil, errors.ExpressionBuilderError
	}

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("plannerix-events"),
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
	log.Println(output)

	weeklyEvents := []Event{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &weeklyEvents)
	if err != nil {
		log.Println("line 99 error with unmarshal")
		return nil, errors.UnmarshalListOfMapsError
	}

	return weeklyEvents, nil
}

func EditEvent(event_id string, id string, subject string, subjectType int, description string, timestamp int64, conn *dynamodb.DynamoDB) error {
	if description == "" {
		description = " "
	}
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("plannerix-events"),
		Key: map[string]*dynamodb.AttributeValue{
			"event_id": {
				S: aws.String(event_id),
			},
			"id": {
				S: aws.String(id),
			},
		},
		UpdateExpression: aws.String("set subject_id = :subject_id, subjectType = :subjectType, description = :description, eventTime = :eventTime"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":subject_id": {
				S: aws.String(subject),
			},
			":subjectType": {
				N: aws.String(strconv.Itoa(subjectType)),
			},
			":description": {
				S: aws.String(description),
			},
			":eventTime": {
				N: aws.String(strconv.FormatInt(timestamp, 10)),
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

func DeleteEvent(event_id string, id string, conn *dynamodb.DynamoDB) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("plannerix-events"),
		Key: map[string]*dynamodb.AttributeValue{
			"event_id": {
				S: aws.String(event_id),
			},
			"id": {
				S: aws.String(id),
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
	location, _ := time.LoadLocation("Europe/Sofia")
	now := time.Now().In(location)
	testDate := time.Unix(timestamp, 0).In(location)
	// log.Println(now)
	// log.Println(testDate)
	diff := testDate.YearDay() - now.YearDay()
	// log.Println(diff)
	newTimestamp := now.AddDate(0, 0, diff).Unix()
	return newTimestamp
}

func generateEventID() string {
	h := fnv.New64a()
	t := time.Now().String()
	h.Write([]byte(t))
	// h.Write([]byte(username))
	h.Write([]byte("kowalski event anal"))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

// func AdaptTimestamp(timestamp int64) int64 {
// 	randInts := random.RangeInt(0, 86399000000000, 1)
// 	randInt := randInts[0]
// 	newTimestamp := int64(randInt) + timestamp
// 	return newTimestamp
// }
