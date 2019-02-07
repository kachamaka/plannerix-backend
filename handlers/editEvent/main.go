package main

import (
	"context"
	"log"

	"gitlab.com/s-org-backend/models/events"
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
	Token       string `json:"token"`
	Subject     string `json:"subject"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Timestamp   int64  `json:"timestamp"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)

	if err != nil {
		return qs.NewError("Internal Server Error", -1)
	}

	database.SetConn(&conn)
	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError("Internal Server Error", 6)
	}

	p := profile.Payload{}
	jwe.ParseEncryptedToken(body.Token, key, &p)
	log.Println(p.Username, "username")

	err = events.EditEvent(p.Username, body.Subject, body.Title, body.Description, body.Timestamp, conn)

	if err != nil {
		log.Println("Error with editing event from database:", err)
		return qs.NewError("Could not edit event", 3)
	}

	return qs.Response{}, nil

	res := Response{
		Success: true,
		Message: "schedule updated successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
