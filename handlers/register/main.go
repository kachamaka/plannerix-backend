package main

import (
	"context"
	"encoding/hex"
	"hash/fnv"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"
	qs "gitlab.com/zapochvam-ei-sq/my-go-service/models/QS"
	"gitlab.com/zapochvam-ei-sq/my-go-service/models/database"
	"gitlab.com/zapochvam-ei-sq/my-go-service/models/profile"
	"golang.org/x/crypto/bcrypt"
)

var conn *dynamodb.DynamoDB

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Response struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func (r Request) validate() error {
	//Validate request and give some feedback
	return nil
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
	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		log.Println("Error with hashing password:", err)
		return qs.NewError("Couldn't register!", 2)
	}
	database.SetConn(&conn)
	id := createID(body.Username)
	p, err := profile.NewProfile(body.Username, body.Email, string(hashed), id, conn)
	if err != nil && strings.Contains(err.Error(), dynamodb.ErrCodeConditionalCheckFailedException) {
		return qs.NewError("Username taken!", 3)
	} else if err != nil {
		log.Println("Error with writing user to database:", err)
		return qs.NewError("Couldn't register user!", 3)
	}
	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")
	if err != nil {
		return qs.NewError("Internal server error! User is registered succesfully!", 8)
	}
	token, err := p.GetToken(&key.PublicKey)
	if err != nil {
		return qs.NewError("Internal server error! User is registered succesfuly!", 9)
	}
	res := Response{
		Success: true,
		Token:   token,
	}
	return qs.NewResponse(200, res)
}
func createID(username string) string {
	h := fnv.New64a()
	t := time.Now().String()
	h.Write([]byte(t))
	h.Write([]byte(username))
	h.Write([]byte("kowalski analysis"))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

func main() {
	lambda.Start(handler)
}
