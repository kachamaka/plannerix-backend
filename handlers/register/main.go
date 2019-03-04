package main

import (
	"context"
	"encoding/hex"
	"hash/fnv"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/goware/emailx"
	"github.com/kinghunter58/jwe"
	qs "gitlab.com/zapochvam-ei-sq/s-org-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/profile"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/subjects"
	"golang.org/x/crypto/bcrypt"
)

var conn *dynamodb.DynamoDB

type Request struct {
	Username string                 `json:"username"`
	Password string                 `json:"password"`
	Email    string                 `json:"email"`
	Subjects []string               `json:"subjects"`
	Schedule []subjects.ScheduleDay `json:"schedule"`
}

type Response struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func (r Request) validate() (error, int) {
	//Validate request and give some feedback
	if !profile.UsernameReg.Match([]byte(r.Username)) {
		return errors.Invalid("потребитлеското име"), 101
	}
	if !profile.PasswordReg.Match([]byte(r.Password)) {
		return errors.Invalid("паролата"), 102
	}
	if len(r.Subjects) == 0 {
		return errors.Invalid("предметите"), 103
	}
	if len(r.Schedule) != 5 {
		return errors.Invalid("програмата"), 104
	}
	err := emailx.Validate(r.Email)
	if err != nil {
		return errors.Invalid("имейла"), 105
	}
	err = sendEmail(r.Email)
	if err != nil {
		return errors.DoesNotExist("Този имейл"), 106
	}

	return nil, 42
}

func sendEmail(email string) error {
	from := "s.org.noreply@gmail.com"
	pass := "kowalskiAnal"
	to := email

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Account activation\n\n" +
		"Account created successfully \n kowalski anal"

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print("sent, visit http://foobarbazz.mailinator.com")
	return nil
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)
	log.Println("req - ", req)
	log.Println("body - ", body)
	log.Println("err - ", err)
	if err != nil {
		log.Println(err)
		return qs.NewError(errors.LambdaError.Error(), -1)
	}
	if err, code := body.validate(); err != nil {
		return qs.NewError(err.Error(), code)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		log.Println("Error with hashing password:", err)
		return qs.NewError(errors.ErrorWith("hashing password").Error(), 107)
	}
	// log.Println("hash", string(hashed))
	// return qs.Response{}, nil
	database.SetConn(&conn)
	id := createID(body.Username)
	p, err := profile.NewProfile(body.Username, body.Email, string(hashed), id, conn)
	if err != nil && strings.Contains(err.Error(), dynamodb.ErrCodeConditionalCheckFailedException) {
		return qs.NewError("Потребителското име е заето!", 108)
	} else if err == errors.MarshalJsonToMapError {
		log.Println("Error with marshaling json to map", err)
		return qs.NewError(err.Error(), 201)
	} else if err == errors.PutItemError {
		return qs.NewError(err.Error(), 300)
	}

	err = subjects.NewSchedule(body.Username, body.Schedule, conn)

	switch err {
	case errors.MarshalMapError:
		return qs.NewError(err.Error(), 200)
	case errors.PutItemError:
		return qs.NewError(err.Error(), 301)
	default:
	}

	err = subjects.NewSubjects(body.Username, body.Subjects, conn)
	switch err {
	case errors.MarshalMapError:
		return qs.NewError(err.Error(), 200)
	case errors.PutItemError:
		return qs.NewError(err.Error(), 303)
	default:
	}

	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")
	if err != nil {
		return qs.NewError("Internal server error! User is registered succesfully!", 109)
	}
	token, err := p.GetToken(&key.PublicKey)
	if err != nil {
		return qs.NewError("Internal server error! User is registered succesfuly!", 110)
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
