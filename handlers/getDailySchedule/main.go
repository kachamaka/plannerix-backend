package main

import (
	"context"
	"time"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/schedule"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"

	"github.com/aws/aws-lambda-go/lambda"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Token string `json:"token"`
}

type Response struct {
	Success  bool                   `json:"success"`
	Message  string                 `json:"message"`
	Schedule schedule.DailySchedule `json:"schedule"`
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)
	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}

	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError(errors.KeyError.Error(), 109)
	}

	p := profile.Payload{}
	jwe.ParseEncryptedToken(body.Token, key, &p)
	database.SetConn(&conn)
	location, _ := time.LoadLocation("Europe/Sofia")
	currentTime := time.Now().In(location)
	ds, err := schedule.GetTodaysSchedule(p.ID, currentTime.Weekday(), conn)
	if err != nil {
		return qs.NewError(err.Error(), 110) //fix
	}
	subjects, err := schedule.GetSubejctsFromDB(p.ID, conn)
	if err != nil {
		return qs.NewError(err.Error(), 110)
	}
	ds.MergeSubjects(subjects)

	res := Response{
		Success:  true,
		Message:  "schedule fetched successfully",
		Schedule: ds,
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
