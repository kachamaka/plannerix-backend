package groups

import (
	"encoding/hex"
	"hash/fnv"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/events"
)

type Group struct {
	ID      string         `json:"group_id"`
	Name    int            `json:"group_name"`
	Owner   string         `json:"owner"`
	Events  []events.Event `json:"group_events"`
	Members []string       `json:"members"`
}

func CreateGroup(groupName string, owner string, conn *dynamodb.DynamoDB) error {
	groupID := generateGroupID()
	events := map[string]interface{}{}
	members := map[string]string{}

	group := map[string]interface{}{
		"group_id":   groupID,
		"group_name": groupName,
		"owner":      owner,
		"events":     events,
		"members":    members,
	}
	body, err := dynamodbattribute.MarshalMap(group)
	if err != nil {
		log.Println("line 45 error with marshal map")
		return errors.MarshalMapError
	}
	input := &dynamodb.PutItemInput{
		TableName:           aws.String("s-org-groups"),
		ConditionExpression: aws.String("attribute_not_exists(group_id)"),
		Item:                body,
	}
	_, err = conn.PutItem(input)
	if err != nil {
		log.Println("line 55 error with put item")
		return errors.PutItemError
	}

	return nil
}

func AddMember(groupID string, owner string, member string, conn *dynamodb.DynamoDB) error {

	filt := expression.Name("group_id").Equal(expression.Value(groupID)).
		And(expression.Name("owner").Equal(expression.Value(owner)))

	proj := expression.NamesList(expression.Name("group_id"), expression.Name("members"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("line 78 couldn't build expression")
		return errors.ExpressionBuilderError
	}

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("s-org-groups"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	output, err := conn.Scan(getItemScanInput)
	if err != nil {
		log.Println("line 92 error with output")
		return errors.OutputError
	}

	group := []Group{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &group)
	if err != nil {
		log.Println("line 99 error with unmarshal")
		return errors.UnmarshalListOfMapsError
	}
	// log.Println(group[0], "group")

	members := group[0].Members
	if contains(members, member) == true {
		return errors.DuplicationError
	}

	newMembers := append(members, member)

	log.Println(newMembers)

	newMembersMarshal, err := dynamodbattribute.MarshalList(newMembers)
	if err != nil {
		log.Println("line 124 err with marshal map", err)
		return errors.MarshalMapError
	}
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("s-org-groups"),
		Key: map[string]*dynamodb.AttributeValue{
			"group_id": {
				S: aws.String(groupID),
			},
			"owner": {
				S: aws.String(owner),
			},
		},
		UpdateExpression: aws.String("set members = :members"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":members": {
				L: newMembersMarshal,
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
	}

	_, err = conn.UpdateItem(updateItemInput)
	if err != nil {
		log.Println("line 148 err with update item", err)
		return errors.UpdateItemError
	}

	return nil
}

func DeleteGroup(groupID string, owner string, conn *dynamodb.DynamoDB) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("s-org-groups"),
		Key: map[string]*dynamodb.AttributeValue{
			"group_id": {
				S: aws.String(groupID),
			},
			"owner": {
				S: aws.String(owner),
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

func generateGroupID() string {
	h := fnv.New64a()
	t := time.Now().String()
	h.Write([]byte(t))
	// h.Write([]byte(username))
	h.Write([]byte("kowalski group analysis"))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
