package main

import (
	"context"
	"log"

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
	Token string `json:"token"`
}

type Response struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	Subjects []string `json:"subjects"`
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

	subjects, err := subjects.GetSubjects(p.Username, conn)
	log.Println(subjects, "all subjects")

	if err != nil {
		log.Println("Error with fetching subjects from database:", err)
		return qs.NewError("Could not get subjects", 3)
	}
	res := Response{
		Success:  true,
		Message:  "subjects fetched successfully",
		Subjects: subjects,
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}