package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"

	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"
)

var conn *dynamodb.DynamoDB

type Request struct {
	Token      string   `json:"token"`
	SubjectIDs []string `json:"subjects"`
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
	for _, v := range body.SubjectIDs {
		input := &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"user_id": {
					S: aws.String(p.ID),
				},
				"id": {
					S: aws.String(v),
				},
			},
			TableName: aws.String("plannerix-subjects"),
		}
		_, err := conn.DeleteItem(input)
		if err != nil {
			return qs.NewError(err.Error(), 304)
		}
	}

	res := Response{
		Success: true,
		Message: "subjects deleted successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
