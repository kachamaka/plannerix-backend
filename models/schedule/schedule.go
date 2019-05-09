package schedule

import (
	"log"
	"strconv"
	"time"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/notifications"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

var (
	errIDNotFound     = errors.New("This id does not exist")
	errNoSuchDay      = errors.New("Such day doesn't exist")
	errNoLessonsToday = errors.New("There are no lessons today")
)

const timeBeforeFirstLesson = 60

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
	s.subjectIDs = map[string]bool{}
	for _, dailyS := range s.Schedule {
		for _, lesson := range dailyS.Lessons {
			// fmt.Println(s.subjectIDs[lesson.SubjectID])
			if !s.subjectIDs[lesson.SubjectID] {
				s.subjectIDs[lesson.SubjectID] = true
			}
		}
	}
}

func (s *Schedule) CheckIfIDsExist() error {
	for k := range s.subjectIDs {
		log.Println(k, s.UserID)
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
		log.Println(k, in)
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
		log.Println(v.Lessons, len(v.Lessons) == 0)
		if len(v.Lessons) == 0 {
			continue
		}
		for i := range v.Lessons {
			v.Lessons[i].Subject = nil
		}
		marshaled, err := dynamodbattribute.MarshalMap(map[string]interface{}{
			"user_id":    s.UserID,
			"day":        DAYS[k],
			"allLessons": v.Lessons,
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

type QueryOut struct {
	Day int `json:"day"`
	DailySchedule
	UserID string `json:"user_id"`
}

func (s *Schedule) GetSchedule() error {
	input := &dynamodb.QueryInput{
		TableName: aws.String("plannerix-schedule"),
		KeyConditions: map[string]*dynamodb.Condition{
			"user_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(s.UserID),
					},
				},
			},
		},
	}
	out, err := s.Conn.Query(input)
	if err != nil {
		return err
	}
	log.Println(out, "out")
	sch := []QueryOut{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &sch)
	if err != nil {
		return err
	}
	log.Println(sch)
	schWtihStrings := map[string]DailySchedule{}
	for _, v := range sch {
		day := getDayFromInt(v.Day)
		schWtihStrings[day] = v.DailySchedule
	}
	log.Println(schWtihStrings)
	s.Schedule = schWtihStrings
	return nil
}

func getDayFromInt(i int) string {
	for k, v := range DAYS {
		if v == i {
			return k
		}
	}
	return ""
}

func (s *Schedule) MergeSubjects(subjects []Subject) {
	for day := range s.Schedule {
		for i := range s.Schedule[day].Lessons {
			index := searchSubjectsById(s.Schedule[day].Lessons[i].SubjectID, subjects)
			s.Schedule[day].Lessons[i].Subject = &Subject{
				Name: subjects[index].Name,
				ID:   subjects[index].ID,
			}
		}
	}
}

type DailySchedule struct {
	Lessons []Lesson `json:"allLessons"`
}

type Lesson struct {
	Start     int      `json:"start"`
	Duration  int      `json:"duration"`
	SubjectID string   `json:"id,omitempty"`
	Subject   *Subject `json:"subject,omitempty"`
}

func GetFirstLessonForDay(userID string, day time.Weekday, conn *dynamodb.DynamoDB) (Lesson, error) {
	todaysShcedule, err := getTodaysSchedule(userID, day, conn)
	if err != nil {
		return Lesson{}, err
	}
	firstLesson, err := getFirstLesson(todaysShcedule)
	if err != nil {
		return Lesson{}, err
	}
	subj, err := getSubjectById(userID, firstLesson.SubjectID, conn)
	if err != nil {
		return Lesson{}, err
	}
	firstLesson.Subject = &subj
	return firstLesson, nil
}

func getTodaysSchedule(userID string, day time.Weekday, conn *dynamodb.DynamoDB) (DailySchedule, error) {
	getInput := dynamodb.GetItemInput{
		TableName: aws.String("plannerix-schedule"),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": &dynamodb.AttributeValue{
				S: &userID,
			},
			"day": &dynamodb.AttributeValue{
				N: aws.String(strconv.Itoa(int(day))),
			},
		},
	}
	out, err := conn.GetItem(&getInput)
	if err != nil {
		return DailySchedule{}, err
	}
	todaysSchedule := DailySchedule{}
	err = dynamodbattribute.UnmarshalMap(out.Item, &todaysSchedule)
	return todaysSchedule, err
}

func getFirstLesson(ds DailySchedule) (Lesson, error) {
	if len(ds.Lessons) == 0 {
		return Lesson{}, errNoLessonsToday
	}
	return ds.Lessons[0], nil
}

func CreateFirstLessonItem(ds DailySchedule, userId string) (notifications.FirstLessonNotificationItem, error) {
	fl, err := getFirstLesson(ds)
	if err != nil {
		return notifications.FirstLessonNotificationItem{}, err
	}
	return notifications.FirstLessonNotificationItem{
		Minutes: fl.Start - timeBeforeFirstLesson,
		UserID:  userId,
	}, nil
}
