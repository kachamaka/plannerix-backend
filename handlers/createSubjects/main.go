package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kinghunter58/jwe"
	errorsWrap "github.com/pkg/errors"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/schedule"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Token    string             `json:"token"`
	Subjects []schedule.Subject `json:"subjects"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)

	fmt.Println(req)
	fmt.Println(body)
	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}

	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError(errors.KeyError.Error(), 109)
	}
	fmt.Println(body)

	p := profile.Payload{}
	jwe.ParseEncryptedToken(body.Token, key, &p)
	database.SetConn(&conn)
	for i := range body.Subjects {
		body.Subjects[i].ID = schedule.CreateID(p.ID)
		body.Subjects[i].UserID = p.ID
		inputBody, err := dynamodbattribute.MarshalMap(body.Subjects[i])
		if err != nil {
			return qs.NewError(errorsWrap.Wrapf(err, "Could not marshal body of %v, index: %v", body.Subjects[i].Name, i).Error(), 300)
		}
		input := &dynamodb.PutItemInput{
			Item:                inputBody,
			TableName:           aws.String("plannerix-subjects"),
			ConditionExpression: aws.String("attribute_not_exists(#name)"),
			ExpressionAttributeNames: map[string]*string{
				"#name": aws.String("name"),
			},
		}
		_, err = conn.PutItem(input)
		if err != nil {
			return qs.NewError(errorsWrap.Wrapf(err, "Could not put item in database | V: %v ; I: %v; ID: %v", body.Subjects[i].Name, i, body.Subjects[i].ID).Error(), 304)
		}
	}
	res := Response{
		Success: true,
		Message: "subjects created successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
