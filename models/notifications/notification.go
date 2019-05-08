package notifications

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//RANGE is also the rate of the cheduler function
const RANGE = 30

func UpdateSubscriptionOfUser(p profile.Payload, subscription string, conn *dynamodb.DynamoDB) error {
	updateInput := dynamodb.UpdateItemInput{
		TableName: aws.String("s-org-users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: &p.Username,
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sub": {
				S: &subscription,
			},
		},
		UpdateExpression: aws.String("set subscription = :sub"),
	}
	_, err := conn.UpdateItem(&updateInput)
	return err
}

func GetUsersInRange(minutes int, conn *dynamodb.DynamoDB) ([]FirstLessonNotificationItem, error) {
	scanInput := dynamodb.ScanInput{
		TableName: aws.String("plannerix-first-lesson-notifications"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":from": &dynamodb.AttributeValue{
				N: aws.String(strconv.Itoa(minutes - RANGE)),
			},
			":to": &dynamodb.AttributeValue{
				N: aws.String(strconv.Itoa(minutes + RANGE)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#time": aws.String("time"),
		},
		FilterExpression: aws.String("#time >= :from AND #time <= :to"),
	}
	out, err := conn.Scan(&scanInput)
	if err != nil {
		return nil, err
	}

	result := []FirstLessonNotificationItem{}
	dynamodbattribute.UnmarshalListOfMaps(out.Items, &result)
	return result, nil
}

func AddUserNotificationItem(body FirstLessonNotificationItem, conn *dynamodb.DynamoDB) error {
	marshled, err := dynamodbattribute.MarshalMap(body)
	if err != nil {
		return err
	}
	putInput := dynamodb.PutItemInput{
		TableName: aws.String("plannerix-first-lesson-notifications"),
		Item:      marshled,
	}
	_, err = conn.PutItem(&putInput)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOldItemForDay(day time.Weekday, userID string, conn *dynamodb.DynamoDB) error {
	items, err := getItemsForUserInDay(day, userID, conn)
	if err != nil {
		return err
	}
	for _, v := range items {
		err := deleteItem(v, conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func getItemsForUserInDay(day time.Weekday, userID string, conn *dynamodb.DynamoDB) ([]FirstLessonNotificationItem, error) {
	queryInput := dynamodb.QueryInput{
		TableName: aws.String("plannerix-first-lesson-notifications"),
		IndexName: aws.String("userIDIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"time": &dynamodb.Condition{
				ComparisonOperator: aws.String("BETWEEN"),
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{
						N: aws.String(strconv.Itoa((int(day) - 1) * minutesInADay)),
					},
					&dynamodb.AttributeValue{
						N: aws.String(strconv.Itoa(int(day) * minutesInADay)),
					},
				},
			},
			"user_id": &dynamodb.Condition{
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{
						S: &userID,
					},
				},
			},
		},
	}
	out, err := conn.Query(&queryInput)
	if err != nil {
		return nil, err
	}
	result := []FirstLessonNotificationItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &result)
	return result, err
}

func deleteItem(item FirstLessonNotificationItem, conn *dynamodb.DynamoDB) error {
	deleteInput := dynamodb.DeleteItemInput{
		TableName: aws.String("plannerix-first-lesson-notifications"),
		Key: map[string]*dynamodb.AttributeValue{
			"time": &dynamodb.AttributeValue{
				N: aws.String(strconv.Itoa(item.Minutes)),
			},
			"user_id": &dynamodb.AttributeValue{
				S: &item.UserID,
			},
		},
	}
	_, err := conn.DeleteItem(&deleteInput)
	return err
}

type FirstLessonNotificationItem struct {
	Minutes int    `json:"time"`
	UserID  string `json:"user_id"`
}
