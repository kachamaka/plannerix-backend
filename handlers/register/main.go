package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"html/template"
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

type EmailRequest struct {
	from    string
	to      []string
	subject string
	body    string
}

var auth smtp.Auth

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

func NewRequest(to []string, subject, body string) *EmailRequest {
	return &EmailRequest{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *EmailRequest) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "plannerix.noreply@gmail.com", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *EmailRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func sendEmail(email string) error {
	// from := "s.org.noreply@gmail.com"
	// from := "plannerix.noreply@gmail.com"
	// pass := "kowalskiAnal"
	// to := email

	// mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// subject := "Subject: Създаване на акаунт\n\n"
	// msg := []byte(subject + mime + "<html><body><p>Акаунтът Ви беше създаден успешно!</p></body></html>")

	// payload := "From: " + from + "\n" +
	// 	"To: " + to + "\n" +
	// 	"Subject: Създаване на акаунт\n\n" + string(msg)

	// err := smtp.SendMail("smtp.gmail.com:587",
	// 	smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
	// 	from, []string{to}, []byte(payload))

	// if err != nil {
	// 	return err
	// }

	auth = smtp.PlainAuth("", "plannerix.noreply@gmail.com", "kowalskiAnal", "smtp.gmail.com")
	templateData := struct {
		Name    string
		URL     string
		From    string
		To      string
		Subject string
	}{
		Name:    "Тест",
		URL:     "https://plannerix.eu",
		From:    "plannerix.noreply@gmail.com",
		To:      email,
		Subject: "Деба",
	}
	r := NewRequest([]string{email}, "Plannerix Account", "")
	err := r.ParseTemplate("template.html", templateData)
	if err := r.ParseTemplate("template.html", templateData); err == nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
		return nil
	}
	return err

}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)
	// log.Println("req - ", req)
	// log.Println("body - ", body)
	// log.Println("err - ", err)
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
