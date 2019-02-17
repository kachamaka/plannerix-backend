package main

import (
	"context"
	"log"

	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/events"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/profile"

	"github.com/kinghunter58/jwe"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	qs "gitlab.com/zapochvam-ei-sq/s-org-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
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

	err = events.DeleteEvent(p.Username, body.Timestamp, conn)

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
