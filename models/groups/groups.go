package groups

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
)

type Group struct {
	ID      string   `json:"group_id"`
	Name    string   `json:"group_name"`
	Owner   string   `json:"owner"`
	Members []string `json:"members"`
}

func CreateGroup(groupName string, owner string, conn *dynamodb.DynamoDB) error {
	groupID := generateGroupID()
	members := map[string]string{}

	group := map[string]interface{}{
		"group_id":   groupID,
		"group_name": groupName,
		"owner":      owner,
		"members":    members,
	}
	body, err := dynamodbattribute.MarshalMap(group)
	if err != nil {
		log.Println("line 45 error with marshal map")
		return errors.MarshalMapError
	}
	input := &dynamodb.PutItemInput{
		// TableName:           aws.String("plannerix-groups"),
		TableName:           aws.String("plannerix-groups"),
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

func DeleteGroup(groupID string, owner string, conn *dynamodb.DynamoDB) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("plannerix-groups"),
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

func getGroup(groupID string, owner string, conn *dynamodb.DynamoDB) (Group, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("plannerix-groups"),
		// IndexName: aws.String("ownerIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"group_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(groupID)}},
			},
			"owner": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(owner)}},
			},
		},
	}

	res, err := conn.Query(queryInput)
	log.Println(res)
	if err != nil {
		log.Println(err)
		return Group{}, err
	}
	// log.Println(res, "test")

	groups := []Group{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &groups)

	if err != nil {
		log.Println("line 99 error with unmarshal")
		log.Println(err)
		return Group{}, errors.UnmarshalListOfMapsError
	}
	if len(groups) == 0 {
		return Group{}, nil
	}

	return groups[0], nil

}

func GetGroupOwner(groupID string, conn *dynamodb.DynamoDB) (string, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("plannerix-groups"),
		// IndexName: aws.String("ownerIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"group_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(groupID)}},
			},
		},
	}

	res, err := conn.Query(queryInput)
	// log.Println(res)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// log.Println(res, "test")

	groups := []Group{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &groups)

	if err != nil {
		log.Println("line 99 error with unmarshal")
		log.Println(err)
		return "", errors.UnmarshalListOfMapsError
	}
	if len(groups) == 0 {
		return "", nil
	}

	return groups[0].Owner, nil

}

func GetOwnedGroups(owner string, conn *dynamodb.DynamoDB) ([]Group, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("plannerix-groups"),
		IndexName: aws.String("ownerIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"owner": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(owner)}},
			},
		},
	}

	res, err := conn.Query(queryInput)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	// log.Println(res, "test")

	groups := []Group{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &groups)

	if err != nil {
		log.Println("line 99 error with unmarshal")
		return nil, errors.UnmarshalListOfMapsError
	}

	return groups, nil
}

func GetMyGroups(user string, conn *dynamodb.DynamoDB) ([]Group, error) {

	filt := expression.Name("owner").NotEqual(expression.Value(user))

	proj := expression.NamesList(expression.Name("group_id"), expression.Name("owner"), expression.Name("group_name"), expression.Name("members"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("line 78 couldn't build expression")
		return nil, errors.ExpressionBuilderError
	}

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("plannerix-groups"),
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
	// log.Println(output)

	allGroups := []Group{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &allGroups)
	if err != nil {
		log.Println("line 99 error with unmarshal")
		log.Println(err)
		return nil, errors.UnmarshalListOfMapsError
	}
	groups := []Group{}
	for _, g := range allGroups {
		if contains(g.Members, user) {
			groups = append(groups, g)
		}
	}

	return groups, nil

}

func AddMember(groupID string, owner string, member string, conn *dynamodb.DynamoDB) error {
	group, err := getGroup(groupID, owner, conn)

	if err != nil {
		return err
	}

	members := group.Members
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
		TableName: aws.String("plannerix-groups"),
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

func DeleteMember(groupID string, owner string, member string, conn *dynamodb.DynamoDB) error {

	group, err := getGroup(groupID, owner, conn)
	if err != nil {
		return err
	}
	members := group.Members

	for i, m := range members {
		if m == member {
			members = append(members[:i], members[i+1:]...)
		}
	}

	membersMarshal, err := dynamodbattribute.MarshalList(members)

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("plannerix-groups"),
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
				L: membersMarshal,
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

func LeaveGroup(groupID string, user string, conn *dynamodb.DynamoDB) error {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("plannerix-groups"),
		// IndexName: aws.String("ownerIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"group_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{{S: aws.String(groupID)}},
			},
		},
	}

	res, err := conn.Query(queryInput)
	log.Println(res)
	if err != nil {
		log.Println(err)
		return err
	}
	// log.Println(res, "test")

	groups := []Group{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &groups)

	if err != nil {
		log.Println("line 99 error with unmarshal")
		log.Println(err)
		return errors.UnmarshalListOfMapsError
	}
	if len(groups) == 0 {
		return errors.DoesNotExist("Group")
	}
	group := groups[0]
	newMembers := []string{}
	for _, member := range group.Members {
		if member != user {
			newMembers = append(newMembers, member)
		}
	}

	newMembersMarshal, err := dynamodbattribute.MarshalList(newMembers)

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("plannerix-groups"),
		Key: map[string]*dynamodb.AttributeValue{
			"group_id": {
				S: aws.String(groupID),
			},
			"owner": {
				S: aws.String(group.Owner),
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

func EditGroupName(groupID string, owner string, groupName string, conn *dynamodb.DynamoDB) error {

	group, err := getGroup(groupID, owner, conn)
	if err != nil {
		return err
	}

	if group.Owner == "" {
		return errors.DoesNotExist("group")
	}

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("plannerix-groups"),
		Key: map[string]*dynamodb.AttributeValue{
			"group_id": {
				S: aws.String(groupID),
			},
			"owner": {
				S: aws.String(owner),
			},
		},
		UpdateExpression: aws.String("set group_name = :group_name"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":group_name": {
				S: aws.String(groupName),
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

func CreateTable(conn *dynamodb.DynamoDB) error {
	params := &dynamodb.CreateTableInput{
		TableName: aws.String("plannerix-events"),
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			&dynamodb.GlobalSecondaryIndex{
				IndexName: aws.String("queryIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					&dynamodb.KeySchemaElement{
						AttributeName: aws.String("id"),
						KeyType:       aws.String("HASH"),
					},
					&dynamodb.KeySchemaElement{
						AttributeName: aws.String("event_id"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("event_id"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("id"), KeyType: aws.String("RANGE")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("event_id"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("id"), AttributeType: aws.String("S")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(3),
			WriteCapacityUnits: aws.Int64(3),
		},
	}

	// create the table
	resp, err := conn.CreateTable(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// print the response data
	fmt.Println(resp)
	return nil
}

func CheckGroupUser(user string, groupID string, conn *dynamodb.DynamoDB) (bool, error) {
	g, err := getGroup(groupID, user, conn)
	if err != nil {
		return false, err
	}
	if g.ID != "" {
		return true, nil
	}

	return false, nil

}

func generateGroupID() string {
	h := fnv.New64a()
	t := time.Now().String()
	h.Write([]byte(t))
	// h.Write([]byte(username))
	h.Write([]byte("kowalski group anal"))
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
