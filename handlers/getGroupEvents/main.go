package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/events"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/groups"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"

	"github.com/aws/aws-lambda-go/lambda"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Token   string `json:"token"`
	GroupID string `json:"group_id"`
}

type Response struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Events  []events.Event `json:"events"`
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

	owner, err := groups.GetGroupOwner(body.GroupID, conn)
	if err != nil {
		return qs.NewError(err.Error(), 0)
	}
	log.Println(owner)

	ownerP, err := profile.GetProfile(owner, conn)
	if err != nil {
		return qs.NewError(err.Error(), 0)
	}

	e, err := events.GetAllEvents(body.GroupID, ownerP.ID, conn)
	log.Println(e, "all events")
	// log.Println(time.Unix(e[0].Timestamp, 0))
	// return qs.Response{}, nil

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
		Success: true,
		Message: "events fetched successfully",
		Events:  e,
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
