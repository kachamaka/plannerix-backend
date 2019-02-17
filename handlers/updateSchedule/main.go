package main

import (
	"context"
	"log"

	"gitlab.com/s-org-backend/models/errors"
	"gitlab.com/s-org-backend/models/profile"
	"gitlab.com/s-org-backend/models/subjects"

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
	Token    string                  `json:"token"`
	Schedule []subjects.ScheduleData `json:"schedule"`
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

	database.SetConn(&conn)
	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError("Internal Server Error", 6)
	}

	p := profile.Payload{}
	jwe.ParseEncryptedToken(body.Token, key, &p)
	log.Println(p.Username, "username")

	err = subjects.UpdateSchedule(p.Username, body.Schedule, conn)
	switch err {
	case errors.MarshalJsonToMapError:
		return qs.NewError(err.Error(), 201)
	case errors.UpdateItemError:
		return qs.NewError(err.Error(), 308)
	}

	res := Response{
		Success: true,
		Message: "schedule updated successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
