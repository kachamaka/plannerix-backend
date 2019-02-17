package main

import (
	"context"
	"log"

	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/profile"
	"golang.org/x/crypto/bcrypt"

	"github.com/kinghunter58/jwe"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
	qs "gitlab.com/zapochvam-ei-sq/s-org-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Token    string `json:"token"`
	Password string `json:"password"`
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

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		log.Println("Error with hashing password:", err)
		return qs.NewError(errors.ErrorWith("hashing password").Error(), 107)
	}
	log.Println(hashed, "hash")

	err = profile.ChangePassword(p.Username, string(hashed), conn)

	switch err {
	case errors.UpdateItemError:
		return qs.NewError(err.Error(), 309)
	default:
	}

	//TODO
	//SEND EMAIL

	res := Response{
		Success: true,
		Message: "password changed successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
