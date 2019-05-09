package main

import (
	"context"

	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/groups"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"

	"github.com/kinghunter58/jwe"

	"github.com/aws/aws-sdk-go/service/dynamodb"

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

	err = groups.LeaveGroup(body.GroupID, p.Username, conn)

	//rewrite
	switch err {
	case errors.DuplicationError:
		return qs.NewError(err.Error(), 111)
	case errors.ExpressionBuilderError:
		return qs.NewError(err.Error(), 206)
	case errors.UnmarshalListOfMapsError:
		return qs.NewError(err.Error(), 204)
	case errors.OutputError:
		return qs.NewError(err.Error(), 205)
	case errors.MarshalMapError:
		return qs.NewError(err.Error(), 200)
	case errors.UpdateItemError:
		return qs.NewError(err.Error(), 311)
	default:
	}

	res := Response{
		Success: true,
		Message: "group left successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
