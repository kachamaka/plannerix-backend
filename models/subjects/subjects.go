package subjects

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
)

type Period struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Subject   string `json:"subject"`
}

type ScheduleDay struct {
	Periods []Period `json:"periods"`
}

type ScheduleData struct {
	Username string   `json:"username"`
	Day      int      `json:"day"`
	Periods  []Period `json:"periods"`
}

type SubjectData struct {
	Username string   `json:"username"`
	Subjects []string `json:"subjects"`
}

func NewSchedule(un string, schedule []ScheduleDay, conn *dynamodb.DynamoDB) error {
	for i := 0; i < len(schedule); i++ {
		periodDay := ScheduleData{
			Day:      i + 1,
			Username: un,
			Periods:  schedule[i].Periods,
		}
		periodDayMap, err := dynamodbattribute.MarshalMap(periodDay)
		if err != nil {
			return errors.MarshalMapError
		}
		input := &dynamodb.PutItemInput{
			TableName: aws.String("s-org-schedules"),
			Item:      periodDayMap,
		}
		_, err = conn.PutItem(input)
		if err != nil {
			log.Println(err)
			return errors.PutItemError
		}
	}
	return nil
}

func NewSubjects(un string, subjects []string, conn *dynamodb.DynamoDB) error {
	userSubjects := SubjectData{
		Username: un,
		Subjects: subjects,
	}
	userSubjectsMap, err := dynamodbattribute.MarshalMap(userSubjects)
	if err != nil {
		return errors.MarshalMapError
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("s-org-subjects"),
		Item:      userSubjectsMap,
	}
	_, err = conn.PutItem(input)
	if err != nil {
		return errors.PutItemError
	}
	return nil
}

func GetSubjects(username string, conn *dynamodb.DynamoDB) ([]string, error) {
	getItemInput := &dynamodb.QueryInput{
		TableName:              aws.String("s-org-subjects"),
		KeyConditionExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
	}
	output, err := conn.Query(getItemInput)
	if err != nil {
		log.Println("line 84", err)
		return nil, errors.OutputError
	}
	subjectsData := []SubjectData{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &subjectsData)
	if err != nil {
		log.Println("line 90")
		return nil, errors.UnmarshalListOfMapsError
	}
	return subjectsData[0].Subjects, nil
}

func GetSchedule(username string, conn *dynamodb.DynamoDB) ([]ScheduleData, error) {
	getItemInput := &dynamodb.QueryInput{
		TableName:              aws.String("s-org-schedules"),
		KeyConditionExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
	}
	output, err := conn.Query(getItemInput)
	if err != nil {
		log.Println("line 82", err)
		return nil, errors.OutputError
	}
	schedule := []ScheduleData{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &schedule)
	if err != nil {
		log.Println("line 88")
		return nil, errors.UnmarshalListOfMapsError
	}
	return schedule, nil
}

func UpdateSchedule(username string, schedule []ScheduleData, conn *dynamodb.DynamoDB) error {
	for i := 0; i < len(schedule); i++ {
		periodsMarshal, err := dynamodbattribute.MarshalList(schedule[i].Periods)
		if err != nil {
			log.Println("line 124 err with marshal list", err)
			return errors.MarshalListError
		}
		updateItemInput := &dynamodb.UpdateItemInput{
			TableName: aws.String("s-org-schedules"),
			Key: map[string]*dynamodb.AttributeValue{
				"username": {
					S: aws.String(username),
				},
				"day": {
					N: aws.String(strconv.Itoa(i + 1)),
				},
			},
			UpdateExpression: aws.String("set periods = :periods"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":periods": {
					L: periodsMarshal,
				},
			},
			ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
		}

		_, err = conn.UpdateItem(updateItemInput)
		if err != nil {
			log.Println("line 148 err with update item", err)
			return errors.UpdateItemError
		}
	}
	return nil
}

func GetNextPeriod(username string, conn *dynamodb.DynamoDB) (Period, error) {
	// .Add(time.Hour * time.Duration(24))
	location, _ := time.LoadLocation("Europe/Sofia")
	t := time.Now().In(location)
	// log.Println(t)
	nowHour := t.Hour()
	nowMinutes := t.Minute()
	nowDay := int(t.Weekday())
	if nowDay == 0 || nowDay == 6 {
		return Period{}, nil
	}
	// log.Println(nowHour, nowMinutes, nowDay)
	//^^^^ correct

	//down - test
	// nowHour := 7
	// nowMinutes := 30
	// nowDay := 4
	//end test

	filt := expression.Name("day").Equal(expression.Value(nowDay)).
		And(expression.Name("username").Equal(expression.Value(username)))
	proj := expression.NamesList(expression.Name("periods"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		return Period{}, errors.ExpressionBuilderError
	}

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("s-org-schedules"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	output, err := conn.Scan(getItemScanInput)
	if err != nil {
		log.Println("line 152", err)
		return Period{}, errors.OutputError
	}

	todaySchedule := []ScheduleData{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &todaySchedule)
	if err != nil {
		log.Println("line 159", err)
		return Period{}, errors.UnmarshalListOfMapsError
	}
	todayPeriods := todaySchedule[0].Periods
	// log.Println(todayPeriods, "periods")
	nextPeriod := Period{}
	for i, period := range todayPeriods {
		periodTime := strings.Split(period.StartTime, ":")
		periodStartHour, _ := strconv.Atoi(periodTime[0])
		periodStartMinute, _ := strconv.Atoi(periodTime[1])
		log.Println(periodStartHour, periodStartMinute)
		if nowHour == periodStartHour {
			if nowMinutes < periodStartMinute {
				nextPeriod = period
				break
			}
		} else if nowHour+1 == periodStartHour {
			nextPeriod = period
			break
		} else if nowHour < periodStartHour {
			if i == 0 {
				nextPeriod = period
				break
			}
		}
	}

	return nextPeriod, nil

}
