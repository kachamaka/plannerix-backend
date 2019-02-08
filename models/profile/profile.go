package profile

import (
	"crypto/rsa"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"gitlab.com/s-org-backend/models/database"

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
		return Profile{}, err
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
		return Profile{}, err
	}
	p := Profile{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &p)
	if err != nil {
		log.Println("line 66 error with unmarshal")
		return Profile{}, nil
	}
	return p, nil
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
