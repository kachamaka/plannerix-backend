package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/subjects"

	"github.com/aws/aws-lambda-go/lambda"
)

var conn *dynamodb.DynamoDB

//Request is the grade input request
type Request struct {
	Token string `json:"token"`
}

type Response struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	NextPeriod subjects.Period `json:"nextPeriod"`
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

	nextPeriod, err := subjects.GetNextPeriod(p.Username, conn)
	log.Println(nextPeriod, "next per")

	switch err {
	case errors.ExpressionBuilderError:
		return qs.NewError(err.Error(), 206)
	case errors.OutputError:
		return qs.NewError(err.Error(), 205)
	case errors.UnmarshalListOfMapsError:
		return qs.NewError(err.Error(), 204)
	default:
	}

	if nextPeriod == (subjects.Period{}) {
		res := Response{
			Success: true,
			Message: "Няма повече часове за днес",
		}
		return qs.NewResponse(200, res)
	}

	res := Response{
		Success:    true,
		Message:    "next period fetched successfully",
		NextPeriod: nextPeriod,
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
