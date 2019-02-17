package main

import (
	"context"
	"log"

	"github.com/goware/emailx"
	"github.com/kinghunter58/jwe"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	qs "gitlab.com/zapochvam-ei-sq/s-org-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/profile"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (r Request) validate() (error, int) {
	//Validate request and give some feedback
	err := emailx.Validate(r.Email)
	if err != nil {
		return errors.Invalid("Email"), 105
	}
	// err = sendEmail(r.Email)
	// if err != nil {
	// 	return errors.DoesNotExist("email"), 106
	// }
	return nil, 42
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)

	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}

	if err, code := body.validate(); err != nil {
		return qs.NewError(err.Error(), code)
	}

	database.SetConn(&conn)
	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError(errors.KeyError.Error(), 109)
	}

	p := profile.Payload{}
	jwe.ParseEncryptedToken(body.Token, key, &p)
	log.Println(p.Username, "username")

	err = profile.ChangeEmail(p.Username, body.Email, conn)

	switch err {
	case errors.UpdateItemError:
		return qs.NewError(err.Error(), 309)
	default:
	}

	//MAYBE SEND EMAIL

	res := Response{
		Success: true,
		Message: "email changed successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
