package schedule

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

type Subject struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"user_id,omitempty"`
}

func addSubjects(subjects []Subject, userID string, db *dynamodb.DynamoDB) error {
	indexes := getIndexesOfUnregisteredSubjects(subjects)
	serializeUnregistered(subjects, indexes, userID)
	err := putInDatabase(subjects, indexes, db)
	return err
}

func getIndexesOfUnregisteredSubjects(subjects []Subject) []int {
	indexes := []int{}
	for i := range subjects {
		if subjects[i].ID == "" {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func serializeUnregistered(subjects []Subject, indexes []int, userID string) {
	for _, v := range indexes {
		subjects[v].ID = CreateID(subjects[v].Name)
		subjects[v].UserID = userID
	}
}

func CreateID(v string) string {
	h := fnv.New64a()
	t := time.Now().String()
	h.Write([]byte(t))
	h.Write([]byte(v))
	h.Write([]byte("kowalski analysis"))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

func putInDatabase(subjects []Subject, indexes []int, db *dynamodb.DynamoDB) error {
	for _, v := range indexes {
		fmt.Println(v)
		inputBody, err := dynamodbattribute.MarshalMap(subjects[v])
		if err != nil {
			return errors.Wrapf(err, "Could not marshal body of %v, index: %v", subjects[v].Name, v)
		}
		input := &dynamodb.PutItemInput{
			Item:      inputBody,
			TableName: aws.String("plannerix-subjects"),
		}
		_, err = db.PutItem(input)
		if err != nil {
			return errors.Wrapf(err, "Could not put item in database | V: %v ; I: %v; ID: %v", subjects[v].Name, v, subjects[v].ID)
		}
	}
	return nil
}

func GetSubejctsFromDB(userID string, db *dynamodb.DynamoDB) ([]Subject, error) {
	qInput := &dynamodb.QueryInput{
		TableName: aws.String("plannerix-subjects"),
		KeyConditions: map[string]*dynamodb.Condition{
			"user_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(userID)}},
			},
		},
	}

	res, err := db.Query(qInput)
	if err != nil {
		return nil, err
	}

	subjects := []Subject{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &subjects)
	return subjects, err
}

func getSubjectById(userId, subjectId string, conn *dynamodb.DynamoDB) (Subject, error) {
	getInput := dynamodb.GetItemInput{
		TableName: aws.String("plannerix-subjects"),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": &dynamodb.AttributeValue{
				S: aws.String(userId),
			},
			"id": &dynamodb.AttributeValue{
				S: aws.String(subjectId),
			},
		},
	}
	out, err := conn.GetItem(&getInput)
	if err != nil {
		return Subject{}, err
	}
	s := Subject{}
	err = dynamodbattribute.UnmarshalMap(out.Item, &s)
	return s, err
}

func checkForUpdate(fromDb []Subject, fromReq []Subject, bd *dynamodb.DynamoDB) {
	for i := range fromDb {
		id := fromDb[i].ID
		index := searchSubjectsById(id, fromReq)
		if index < 0 {
			continue
		}
		if fromDb[i].Name != fromReq[index].Name {

		}
	}
}

func searchSubjectsById(ID string, s []Subject) int {
	for i, v := range s {
		if v.ID == ID {
			return i
		}
	}
	return -1
}

// add those who have no id to db & update those that need to be updated & delete others

// get subjects based on user id
