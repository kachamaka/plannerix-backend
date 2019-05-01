package schedule

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

var (
	errIDNotFound = errors.New("This id does not exist")
	errNoSuchDay  = errors.New("Such day doesn't exist")
)

var DAYS = map[string]int{
	"monday":    1,
	"tuesday":   2,
	"wednesday": 3,
	"thursday":  4,
	"friday":    5,
	"saturday":  6,
	"sunday":    7,
}

type Schedule struct { // ?? should it exist
	Schedule   map[string]DailySchedule
	subjectIDs map[string]bool // this is all ids in this schedule
	Conn       *dynamodb.DynamoDB
	UserID     string
}

func (s *Schedule) LoadSubjectIDs() {
	fmt.Println(s.Schedule)
	s.subjectIDs = map[string]bool{}
	for _, dailyS := range s.Schedule {
		for _, lesson := range dailyS.Lessons {
			fmt.Println(s.subjectIDs[lesson.SubjectID])
			if !s.subjectIDs[lesson.SubjectID] {
				s.subjectIDs[lesson.SubjectID] = true
			}
		}
	}
	fmt.Println(s.subjectIDs)
}

func (s *Schedule) CheckIfIDsExist() error {
	for k := range s.subjectIDs {
		in := &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"user_id": {
					S: aws.String(s.UserID),
				},
				"id": {
					S: aws.String(k),
				},
			},
			TableName:            aws.String("plannerix-subjects"),
			ProjectionExpression: aws.String("#name"),
			ExpressionAttributeNames: map[string]*string{
				"#name": aws.String("name"),
			},
		}
		out, err := s.Conn.GetItem(in)
		if err != nil {
			return err
		}
		if len(out.Item) == 0 {
			return errors.Wrap(errIDNotFound, k)
		}
	}
	return nil
}

func (s Schedule) InsertItem() error {
	for k, v := range s.Schedule {
		if _, ok := DAYS[k]; !ok {
			return errNoSuchDay
		}
		marshaled, err := dynamodbattribute.MarshalMap(map[string]interface{}{
			"user_id": s.UserID,
			"day":     DAYS[k],
			"lessons": v.Lessons,
		})
		if err != nil {
			return err
		}
		input := &dynamodb.PutItemInput{
			Item:      marshaled,
			TableName: aws.String("plannerix-schedule"),
		}
		_, err = s.Conn.PutItem(input)
		if err != nil {
			return err
		}
	}
	return nil
}

type DailySchedule struct {
	Lessons []Lesson `json:"lessons"`
}

type Lesson struct {
	Start     int    `json:"start"`
	Duration  int    `json:"duration"`
	SubjectID string `json:"id"`
}

// to save
// Unmarshal -> delete previous schedule -> set new schedule

//to get
// get from db -> parse in Schedule -> get Subjects -> join on subject id -> return to user
// Note: maybe use context package for performence boost
