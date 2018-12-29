package main

import (
	"context"
	"log"

	"github.com/kinghunter58/jwe"

	"gitlab.com/zapochvam-ei-sq/my-go-service/models/database"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gitlab.com/zapochvam-ei-sq/my-go-service/models/profile"

	"github.com/aws/aws-lambda-go/lambda"
	qs "gitlab.com/zapochvam-ei-sq/my-go-service/models/QS"
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

func handler(ctx context.Context, req qs.Request) (qs.Response, error) {
	var body = Request{}
	err := req.ReadTo(&body)
	if err != nil {
		log.Println("Error by decoding request to struct:", err)
		return qs.NewError("Couldn't decode body from request", 0)
	}
	if err := body.validate(); err != nil {
		return qs.NewError(err.Error(), 1)
	}
	database.SetConn(&conn)
	p, err := profile.GetProfile(body.Username, conn)
	if err != nil {
		log.Println("Error with getting user from data base")
		return qs.NewError("Could not find user profile", 3)
	}
	if p.Username == "" {
		return qs.NewError("Could not find user", 4)
	}
	if ok := p.CheckPassword(body.Password); !ok {
		return qs.NewError("Password not correct", 5)
	}

	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")
	log.Println(key)
	if err != nil {
		log.Println(err)
		return qs.NewError("Internal Server Error", 6)
	}
	token, err := p.GetToken(&key.PublicKey)
	if err != nil {
		log.Println(err)
		return qs.NewError("Internal Server Error", 7)
	}
	res := Response{
		Token:   token,
		Success: true,
	}
	return qs.NewResponse(200, map[string]string{"Content-Type": "application/json"}, res)
}

func main() {
	lambda.Start(handler)
}
