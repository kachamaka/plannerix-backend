package main

import (
	"context"

	"gitlab.com/s-org-backend/models/errors"
	"gitlab.com/s-org-backend/models/grades"
	"gitlab.com/s-org-backend/models/profile"

	"github.com/kinghunter58/jwe"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	qs "gitlab.com/s-org-backend/models/QS"
	"gitlab.com/s-org-backend/models/database"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Token   string `json:"token"`
	Subject string `json:"subject"`
	Value   int    `json:"value"`
	Time    int64  `json:"time"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (r Request) validate() error {
	if r.Value < 2 || r.Value > 6 {
		return errors.Invalid("grade value")
	}
	return nil
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)

	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}

	if err := body.validate(); err != nil {
		return qs.NewError(err.Error(), 100)
	}

	database.SetConn(&conn)
	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError(errors.KeyError.Error(), 109)
	}

	p := profile.Payload{}
	jwe.ParseEncryptedToken(body.Token, key, &p)
	err = grades.InputGrade(p.Username, body.Time, body.Value, body.Subject, conn)

	switch err {
	case errors.MarshalJsonToMapError:
		return qs.NewError(err.Error(), 201)
	case errors.PutItemError:
		return qs.NewError(err.Error(), 302)
	default:
	}

	res := Response{
		Success: true,
		Message: "grade inserted successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
