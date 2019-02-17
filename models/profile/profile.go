package profile

import (
	"crypto/rsa"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/kinghunter58/jwe"
	"golang.org/x/crypto/bcrypt"
)

//Profile is the data structure of a user profile
type Profile struct {
	Username string
	Password string // hashed password
	Email    string
	ID       string
	// Add more fields here
}

func NewProfile(un, email, password, id string, conn *dynamodb.DynamoDB) (Profile, error) {
	j := fmt.Sprintf(`{"username": "%v", "email": "%v","password": "%v", "id":"%v"}`, un, email, password, id)
	body, err := database.MarshalJSONToDynamoMap(j)
	if err != nil {
		log.Println("line 31 error with marshal json to map")
		return Profile{}, errors.MarshalJsonToMapError
	}
	input := &dynamodb.PutItemInput{
		TableName:           aws.String("s-org-users"),
		ConditionExpression: aws.String("attribute_not_exists(username)"),
		Item:                body,
	}
	_, err = conn.PutItem(input)
	if err != nil {
		log.Println("line 42 error with put item")
		return Profile{}, err
	}
	return Profile{Username: un, Password: password, Email: email, ID: id}, nil
}

//GetProfile gets the profile of the user from the db
func GetProfile(username string, conn *dynamodb.DynamoDB) (Profile, error) {
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("s-org-users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	}
	output, err := conn.GetItem(getItemInput)
	if err != nil {
		log.Println("line 60 error with output")
		return Profile{}, errors.OutputError
	}
	p := Profile{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &p)
	if err != nil {
		log.Println("line 66 error with unmarshal")
		return Profile{}, errors.UnmarshalMapError
	}
	return p, nil
}

func ChangePassword(username string, newPassword string, conn *dynamodb.DynamoDB) error {
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("s-org-users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
		UpdateExpression: aws.String("set password = :password"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":password": {
				S: aws.String(newPassword),
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
	}

	_, err := conn.UpdateItem(updateItemInput)
	if err != nil {
		log.Println("line 107 error with update item")
		return errors.UpdateItemError
	}
	return nil
}

func ChangeEmail(username string, newEmail string, conn *dynamodb.DynamoDB) error {
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("s-org-users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
		UpdateExpression: aws.String("set email = :email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(newEmail),
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
	}

	_, err := conn.UpdateItem(updateItemInput)
	if err != nil {
		log.Println("line 107 error with update item")
		return errors.UpdateItemError
	}
	return nil
}

//CheckPassword is self-explanatory returns true if success
func (p Profile) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}

//Payload is the data that gets encrypted in the token
type Payload struct {
	Username string `json:"username"`
}

func (p Profile) newPayload() Payload {
	return Payload{Username: p.Username}
}

//GetToken creates a jwe token from the profile data
func (p Profile) GetToken(pk *rsa.PublicKey) (string, error) {
	payl := p.newPayload()
	// fmt.Println(payl)
	return jwe.GetEncryptedToken(payl, pk)
}

var UsernameReg = regexp.MustCompile("^\\w.{3,16}$")
var PasswordReg = regexp.MustCompile("^[a-z0-9]{8,35}$")
