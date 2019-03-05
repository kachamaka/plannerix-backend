package profile

import (
	"bytes"
	"crypto/rsa"
	"log"
	"net/http"
	"net/smtp"
	"regexp"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/kinghunter58/jwe"
	"golang.org/x/crypto/bcrypt"
)

//Profile is the data structure of a user profile
type Profile struct {
	Username      string          `json:"username"`
	Password      string          `json:"password"`
	Email         string          `json:"email"`
	Notifications map[string]bool `json:"notifications"`
	ID            string
	// Add more fields here
}

type EmailRequest struct {
	from    string
	to      string
	subject string
	body    string
}

var Auth smtp.Auth

func NewRequest(to string, subject, body string) *EmailRequest {
	return &EmailRequest{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *EmailRequest) SendEmail(email string) (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, Auth, "plannerix.noreply@gmail.com", []string{email, "plannerix.noreply@gmail.com"}, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *EmailRequest) ParseTemplate(assets http.FileSystem, templateFileName string, data interface{}) error {
	t, err := vfstemplate.ParseFiles(assets, nil, templateFileName)

	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func NewProfile(un, email, password, id string, conn *dynamodb.DynamoDB) (Profile, error) {
	notifications := map[string]bool{
		"all":         true,
		"events":      true,
		"improvement": true,
		"period":      true,
	}
	profile := map[string]interface{}{
		"username":      un,
		"email":         email,
		"password":      password,
		"id":            id,
		"notifications": notifications,
	}
	body, err := dynamodbattribute.MarshalMap(profile)
	if err != nil {
		log.Println("line 45 error with marshal map")
		return Profile{}, errors.MarshalMapError
	}
	input := &dynamodb.PutItemInput{
		TableName:           aws.String("s-org-users"),
		ConditionExpression: aws.String("attribute_not_exists(username)"),
		Item:                body,
	}
	_, err = conn.PutItem(input)
	if err != nil {
		log.Println("line 55 error with put item")
		return Profile{}, err
	}
	return Profile{Username: un, Password: password, Email: email, ID: id}, nil
}

func DeleteProfile(username string, conn *dynamodb.DynamoDB) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("s-org-users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	}
	_, err := conn.DeleteItem(deleteItemInput)
	if err != nil {
		log.Println("line 169 error with unmarshal")
		return errors.DeleteItemError
	}
	return nil
}

//GetProfile gets the profile of the user from the db
func GetProfile(username string, conn *dynamodb.DynamoDB) (Profile, error) {

	filt := expression.Name("username").Equal(expression.Value(username))

	proj := expression.NamesList(expression.Name("username"), expression.Name("email"), expression.Name("password"), expression.Name("notifications"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("line 78 couldn't build expression")
		return Profile{}, errors.ExpressionBuilderError
	}

	getItemScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("s-org-users"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	output, err := conn.Scan(getItemScanInput)
	if err != nil {
		log.Println("line 92 error with output")
		return Profile{}, errors.OutputError
	}
	p := []Profile{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &p)
	if err != nil {
		log.Println("line 99 error with unmarshal")
		return Profile{}, errors.UnmarshalListOfMapsError
	}
	log.Println(p)
	return p[0], nil
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

func UpdateNotifications(username string, notifications map[string]bool, conn *dynamodb.DynamoDB) error {
	notificationsMarshal, err := dynamodbattribute.MarshalMap(notifications)
	if err != nil {
		log.Println("line 124 err with marshal map", err)
		return errors.MarshalMapError
	}
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("s-org-users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
		UpdateExpression: aws.String("set notifications = :notifications"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":notifications": {
				M: notificationsMarshal,
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
