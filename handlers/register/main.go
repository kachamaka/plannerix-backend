package main

import (
	"context"
	"encoding/hex"
	"errors"
	"hash/fnv"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/goware/emailx"
	"github.com/kinghunter58/jwe"
	qs "gitlab.com/s-org-backend/models/QS"
	"gitlab.com/s-org-backend/models/database"
	"gitlab.com/s-org-backend/models/profile"
	"gitlab.com/s-org-backend/models/subjects"
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

func (r Request) validate() error {
	//Validate request and give some feedback
	if !profile.UsernameReg.Match([]byte(r.Username)) {
		return errors.New("Invalid Username")
	}
	if !profile.PasswordReg.Match([]byte(r.Password)) {
		return errors.New("Invalid Password")
	}
	if len(r.Subjects) == 0 {
		return errors.New("Please select subjects")
	}
	if len(r.Schedule) != 5 {
		return errors.New("Please input schedule for every day")
	}
	err := emailx.Validate(r.Email)
	if err != nil {
		return errors.New("Invalid email")
	}
	err = sendEmail(r.Email)
	if err != nil {
		return errors.New("Email does not exist")
	}

	return nil
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
	// fmt.Println(body)
	if err != nil {
		log.Println(err)
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

	err = subjects.NewSchedule(body.Username, body.Schedule, conn)
	if err != nil {
		return qs.NewError("Error with schedule", 8)
	}

	err = subjects.NewSubjects(body.Username, body.Subjects, conn)
	if err != nil {
		return qs.NewError("Error with schedule", 8)
	}

	return qs.Response{}, nil

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
