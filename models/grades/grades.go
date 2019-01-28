package grades

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jinzhu/now"
	"gitlab.com/s-org-backend/models/database"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type SingleGrade struct {
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

func GetAllGrades(user_id string, conn *dynamodb.DynamoDB) ([]SingleGrade, error) {
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
		return []SingleGrade{}, err
	}
	grades := []SingleGrade{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &grades)
	if err != nil {
		log.Println("line 60")
		return []SingleGrade{}, nil
	}
	return grades, nil
}

func GetWeeklyGrades(user_id string, conn *dynamodb.DynamoDB) ([]SingleGrade, error) {
	grades, err := GetAllGrades(user_id, conn)
	if err != nil {
		log.Println("line 68")
		return []SingleGrade{}, nil
	}
	weeklyGrades := []SingleGrade{}
	mondayTime := now.BeginningOfWeek().UnixNano()
	for _, el := range grades {
		if el.Timestamp >= mondayTime {
			weeklyGrades = append(weeklyGrades, el)
		}
	}

	return weeklyGrades, nil
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
		log.Println("line 97", err)
		return err
	}
	return nil
}
