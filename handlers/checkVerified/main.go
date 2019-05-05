package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"

	"github.com/aws/aws-lambda-go/lambda"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the verification input request
type Request struct {
	Username string `json:"username"`
}

//Response is
type Response struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Verified bool   `json:"verified"`
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)

	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}

	database.SetConn(&conn)
	// key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	// if err != nil {
	// 	return qs.NewError(errors.KeyError.Error(), 109)
	// }

	verified, err := profile.CheckVerified(body.Username, conn)
	log.Println(verified, err)

	// switch err {
	// case errors.OutputError:
	// 	return qs.NewError(err.Error(), 205)
	// case errors.UnmarshalListOfMapsError:
	// 	return qs.NewError(err.Error(), 204)
	// default:
	// }

	res := Response{
		Success:  true,
		Message:  "checked verification successfully",
		Verified: verified,
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
