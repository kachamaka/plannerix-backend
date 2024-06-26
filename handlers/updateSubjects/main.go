package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	db "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"

	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/schedule"
)

var conn *dynamodb.DynamoDB

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

	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}

	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError("Internal Server Error", 6)
	}

	p := profile.Payload{}
	err = jwe.ParseEncryptedToken(body.Token, key, &p)
	if err != nil {
		return qs.NewError("Internal Server Error", 6) // Fix output
	}
	database.SetConn(&conn)
	for i := range body.Subjects {
		input := &db.UpdateItemInput{
			TableName: aws.String("plannerix-subjects"),
			Key: map[string]*dynamodb.AttributeValue{
				"id":      {S: aws.String(body.Subjects[i].ID)},
				"user_id": {S: aws.String(p.ID)},
			},
			ExpressionAttributeNames: map[string]*string{
				"#name": aws.String("name"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":newName": {S: aws.String(body.Subjects[i].Name)},
			},
			UpdateExpression: aws.String("set #name = :newName"),
		}
		_, err := conn.UpdateItem(input)
		if err != nil {
			return qs.NewError(err.Error(), 304)
		}
	}

	res := Response{
		Success: true,
		Message: "subjects updated successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
