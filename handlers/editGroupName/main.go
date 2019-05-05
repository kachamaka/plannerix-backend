package main

import (
	"context"
	"log"

	"github.com/kinghunter58/jwe"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/groups"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Token     string `json:"token"`
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
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
	log.Println(p.Username, "username")

	err = groups.EditGroupName(body.GroupID, p.Username, body.GroupName, conn)

	log.Println(err)
	switch err {
	case errors.UpdateItemError:
		return qs.NewError(err.Error(), 307)
	default:
	}

	if err != nil {
		return qs.NewError(err.Error(), 404)
	}

	res := Response{
		Success: true,
		Message: "group name updated successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}