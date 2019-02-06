package subjects

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Period struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Subject   string `json:"subject"`
}

type ScheduleDay struct {
	Periods []Period `json:"periods"`
}
type Schedule struct {
	Username string   `json:"username"`
	Day      int      `json:"day"`
	Periods  []Period `json:"periods"`
}

func NewSchedule(un string, schedule []ScheduleDay, conn *dynamodb.DynamoDB) error {
	for i := 0; i < len(schedule); i++ {
		periodDay := Schedule{
			Day:      i + 1,
			Username: un,
			Periods:  schedule[i].Periods,
		}
		periodDayMap, err := dynamodbattribute.MarshalMap(periodDay)
		input := &dynamodb.PutItemInput{
			TableName: aws.String("s-org-schedules"),
			Item:      periodDayMap,
		}
		_, err = conn.PutItem(input)
		if err != nil {
			return err
		}
	}
	return nil
}
