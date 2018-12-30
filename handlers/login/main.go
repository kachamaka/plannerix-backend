package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"

	"github.com/aws/aws-lambda-go/lambda"
	qs "gitlab.com/zapochvam-ei-sq/my-go-service/models/QS"
	"gitlab.com/zapochvam-ei-sq/my-go-service/models/database"
	"gitlab.com/zapochvam-ei-sq/my-go-service/models/profile"
)

var conn *dynamodb.DynamoDB

//Request is the login request
type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r Request) validate() error {
	//validate request body
	return nil
}

type Response struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)
	if err != nil {
		return qs.NewError("Internal Server Error", -1)
	}
	if err := body.validate(); err != nil {
		return qs.NewError(err.Error(), 1)
	}
	database.SetConn(&conn)
	p, err := profile.GetProfile(body.Username, conn)
	if err != nil {
		log.Println("Error with getting user from data base:", err)
		return qs.NewError("Could not find user profile", 3)
	}
	if p.Username == "" {
		return qs.NewError("Could not find user", 4)
	}
	if ok := p.CheckPassword(body.Password); !ok {
		return qs.NewError("Password not correct", 5)
	}

	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")
	if err != nil {
		return qs.NewError("Internal Server Error", 6)
	}
	token, err := p.GetToken(&key.PublicKey)
	if err != nil {
		return qs.NewError("Internal Server Error", 7)
	}
	res := Response{
		Token:   token,
		Success: true,
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
