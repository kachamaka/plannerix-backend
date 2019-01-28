package main

import (
	"context"
	"log"

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
		return qs.NewError("Internal Server Error", -1)
	}

	database.SetConn(&conn)
	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError("Internal Server Error", 6)
	}

	p := profile.Payload{}
	jwe.ParseEncryptedToken(body.Token, key, &p)
	log.Println(p.ID, "id")

	err = grades.DeleteGrade(p.ID, body.Timestamp, conn)

	if err != nil {
		log.Println("Error with deleting a grade from database:", err)
		return qs.NewError("Could not delete grade", 3)
	}
	res := Response{
		Success: true,
		Message: "grade deleted successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
