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
	"github.com/goware/emailx"
	mailgun "github.com/mailgun/mailgun-go"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"
	"golang.org/x/crypto/bcrypt"
)

var conn *dynamodb.DynamoDB

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	// Subjects []string               `json:"subjects"`
	// Schedule []subjects.ScheduleDay `json:"schedule"`
}

type Response struct {
	Success bool `json:"success"`
}

func (r Request) validate() (error, int) {
	//Validate request and give some feedback
	if !profile.UsernameReg.Match([]byte(r.Username)) {
		return errors.Invalid("потребитлеското име"), 101
	}
	if !profile.PasswordReg.Match([]byte(r.Password)) {
		return errors.Invalid("паролата"), 102
	}
	// if len(r.Subjects) == 0 {
	// 	return errors.Invalid("предметите"), 103
	// }
	// if len(r.Schedule) != 5 {
	// 	return errors.Invalid("програмата"), 104
	// }

	err := emailx.Validate(r.Email)
	if err != nil {
		return errors.Invalid("имейла"), 105
	}

	return nil, 42
}

// func sendVerificationKeyEmail(email string, verificationKey string) error {
// 	profile.Auth = smtp.PlainAuth("", "plannerix.noreply@gmail.com", "kowalskiAnal", "smtp.gmail.com")
// 	templateData := struct {
// 		Name    string
// 		URL     string
// 		From    string
// 		To      string
// 		Subject string
// 	}{
// 		Name:    "Тест",
// 		URL:     "https://plannerix.eu/link?verificationKey=" + verificationKey,
// 		From:    "plannerix.noreply@gmail.com",
// 		To:      email,
// 		Subject: "Създаване на акаунт",
// 	}
// 	r := profile.NewRequest(email, "Plannerix Account", "")
// 	err := r.ParseTemplate(assets, "template.html", templateData)
// 	if err == nil {
// 		ok, _ := r.SendEmail(email)
// 		fmt.Println(ok)
// 		return nil
// 	}
// 	return err

// }

func sendVerificationKeyEmail(email string, verificationKey string) error {
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun("sandboxd15ea31d048e4f92b7225d795260ccb5.mailgun.org", "57d452adb41fbed21081e953187a2de7-3fb021d1-6915ac04")

	sender := "plannerix.noreply@gmail.com"
	subject := "TEST subject!"
	body := "Hello from Mailgun Go!\n" + "You can verify your Plannerix account here:\n " + "https://plannerix.eu/link?verificationKey=" + verificationKey
	recipient := email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return err
	}

	log.Println(id, resp)
	return nil
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
	verificationKey, err := profile.NewUnverifiedProfile(body.Username, body.Email, string(hashed), id, conn)
	if err != nil && strings.Contains(err.Error(), dynamodb.ErrCodeConditionalCheckFailedException) {
		return qs.NewError("Потребителското име е заето!", 108)
	} else if err == errors.MarshalJsonToMapError {
		log.Println("Error with marshaling json to map", err)
		return qs.NewError(err.Error(), 201)
	} else if err == errors.PutItemError {
		return qs.NewError(err.Error(), 300)
	}

	err = sendVerificationKeyEmail(body.Email, verificationKey)
	log.Println("email err", err)
	if err != nil {
		err = profile.DeleteUnverifiedProfile(body.Username, conn)
		if err != nil {
			return qs.NewError(err.Error(), 310)
		}
		return qs.NewError(errors.DoesNotExist("Този имейл").Error(), 106)
	}

	// err = subjects.NewSchedule(body.Username, body.Schedule, conn)

	// switch err {
	// case errors.MarshalMapError:
	// 	return qs.NewError(err.Error(), 200)
	// case errors.PutItemError:
	// 	return qs.NewError(err.Error(), 301)
	// default:
	// }

	// err = subjects.NewSubjects(body.Username, body.Subjects, conn)
	// switch err {
	// case errors.MarshalMapError:
	// 	return qs.NewError(err.Error(), 200)
	// case errors.PutItemError:
	// 	return qs.NewError(err.Error(), 303)
	// default:
	// }

	res := Response{
		Success: true,
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
