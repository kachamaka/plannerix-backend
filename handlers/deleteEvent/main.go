package main

import (
	"context"
	"log"

	"github.com/kinghunter58/jwe"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/events"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Token   string `json:"token"`
	EventID string `json:"event_id"`
}

//Response is the grade input request
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
		return qs.NewError(errors.KeyError.Error(), 109)
	}

	p := profile.Payload{}
	jwe.ParseEncryptedToken(body.Token, key, &p)
	log.Println(p.Username, "username")

	err = events.DeleteEvent(body.EventID, p.ID, conn)

	switch err {
	case errors.DeleteItemError:
		return qs.NewError(err.Error(), 306)
	default:
	}

	res := Response{
		Success: true,
		Message: "event deleted successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
