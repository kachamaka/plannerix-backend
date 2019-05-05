package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"

	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/schedule"
)

var conn *dynamodb.DynamoDB

type Request struct {
	Token string                            `json:"token"`
	Data  map[string]schedule.DailySchedule `json:"schedule"`
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
	sch := schedule.Schedule{
		Schedule: body.Data,
		UserID:   p.ID,
		Conn:     conn,
	}
	sch.LoadSubjectIDs()
	log.Println(sch)
	err = sch.CheckIfIDsExist()
	if err != nil {
		return qs.NewError(err.Error(), 7)
	}
	err = sch.InsertItem()
	if err != nil {
		return qs.NewError(err.Error(), 8)
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
