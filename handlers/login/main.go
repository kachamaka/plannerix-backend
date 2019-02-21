package main

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"
	qs "gitlab.com/zapochvam-ei-sq/s-org-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/profile"

	"github.com/aws/aws-lambda-go/lambda"
)

var conn *dynamodb.DynamoDB

//Request is the login request
type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r Request) validate() error {
	if !profile.UsernameReg.Match([]byte(r.Username)) {
		return errors.Invalid("Username")
	}
	// if !profile.PasswordReg.Match([]byte(r.Password)) {
	// 	return errors.New("Invalid Password")
	// }
	//^^^ correct
	return nil
}

type Response struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)
	// log.Println(req, body)
	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}
	if err := body.validate(); err != nil {
		return qs.NewError(err.Error(), 100)
	}
	database.SetConn(&conn)
	p, err := profile.GetProfile(body.Username, conn)
	// log.Printf("Profile %+v:  %+v\n", p, body, req)
	switch err {
	case errors.OutputError:
		return qs.NewError(err.Error(), 205)
	case errors.UnmarshalMapError:
		return qs.NewError(err.Error(), 203)
	default:
	}
	if p.Username == "" {
		return qs.NewError(errors.NotFound("User").Error(), 404)
	}
	if ok := p.CheckPassword(body.Password); !ok {
		return qs.NewError("Password not correct", 102)
	}

	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")
	if err != nil {
		return qs.NewError(errors.KeyError.Error(), 109)
	}
	token, err := p.GetToken(&key.PublicKey)
	if err != nil {
		return qs.NewError(errors.TokenError.Error(), 110)
	}
<<<<<<< HEAD
	// log.Println(token, "token")
	// log.Println("???")
=======

>>>>>>> 44d95e0db972e2de39c03746a560f7de83e3c3f0
	res := Response{
		Token:   token,
		Success: true,
	}
<<<<<<< HEAD
	r, err := qs.NewResponse(200, res)
	// log.Println(r)
=======

	r, err := qs.NewResponse(200, res)
	log.Println(r)

>>>>>>> 44d95e0db972e2de39c03746a560f7de83e3c3f0
	return r, err
}

func main() {
	lambda.Start(handler)
}
