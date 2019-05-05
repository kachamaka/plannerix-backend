package main

import (
	"context"
	"log"

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
	Success  bool                              `json:"success"`
	Message  string                            `json:"message"`
	Schedule map[string]schedule.DailySchedule `json:"schedule"`
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)
	log.Println(body)
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
	// log.Println(p.Username, "username")
	sch := schedule.Schedule{
		UserID: p.ID,
		Conn:   conn,
	}
	err = sch.GetSchedule()
	if err != nil {
		//handle
		return qs.NewError(err.Error(), 6)
	}
	log.Println("schedule", sch)
	log.Println("=======================================")
	subjects, err := schedule.GetSubejctsFromDB(p.ID, conn)
	if err != nil {
		//handle
		return qs.NewError(err.Error(), 6)
	}
	sch.MergeSubjects(subjects)

	res := Response{
		Success:  true,
		Message:  "schedule fetched successfully",
		Schedule: sch.Schedule,
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
