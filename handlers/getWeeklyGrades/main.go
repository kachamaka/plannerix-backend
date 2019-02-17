package main

import (
	"context"
	"log"

	"github.com/kinghunter58/jwe"
	qs "gitlab.com/zapochvam-ei-sq/s-org-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/grades"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/profile"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
)

var conn *dynamodb.DynamoDB

//Request is the grade input request
type Request struct {
	Token string `json:"token"`
}

type Response struct {
	Success      bool           `json:"success"`
	Message      string         `json:"message"`
	WeeklyGrades []grades.Grade `json:"weeklyGrades"`
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

	weeklyGrades, err := grades.GetWeeklyGrades(p.Username, conn)
	log.Println(weeklyGrades, "weekly grades")

	switch err {
	case errors.ExpressionBuilderError:
		return qs.NewError(err.Error(), 206)
	case errors.OutputError:
		return qs.NewError(err.Error(), 205)
	case errors.UnmarshalListOfMapsError:
		return qs.NewError(err.Error(), 204)
	default:
	}

	res := Response{
		Success:      true,
		Message:      "grades fetched successfully",
		WeeklyGrades: weeklyGrades,
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
