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
	Token string `json:"token"`
}

type Response struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Grades  []grades.Grade `json:"grades"`
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

	g, err := grades.GetYearGrades(p.ID, conn)
	log.Println(g, "all grades")

	if err != nil {
		log.Println("Error with fetching grades from database:", err)
		return qs.NewError("Could not get grades", 3)
	}
	res := Response{
		Success: true,
		Message: "grades fetched successfully",
		Grades:  g,
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
