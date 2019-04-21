package main

import (
	"context"
	"encoding/hex"
	"hash/fnv"
	"log"
	"time"

	qs "gitlab.com/zapochvam-ei-sq/s-org-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/profile"
	"golang.org/x/crypto/bcrypt"

	mailgun "github.com/mailgun/mailgun-go"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/lambda"
)

var conn *dynamodb.DynamoDB

//todo grade struct

//Request is the grade input request
type Request struct {
	Username string `json:"username"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func sendResetPasswordEmail(email string, newPass string) error {
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun("sandboxd15ea31d048e4f92b7225d795260ccb5.mailgun.org", "57d452adb41fbed21081e953187a2de7-3fb021d1-6915ac04")

	sender := "plannerix.noreply@gmail.com"
	subject := "TEST subject!"
	body := "Hello from Mailgun Go!\n" + "Your dummy password is:\n " + newPass
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

	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}

	database.SetConn(&conn)

	newPass := createNewPassword()
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), 12)
	if err != nil {
		log.Println("Error with hashing password:", err)
		return qs.NewError(errors.ErrorWith("hashing password").Error(), 107)
	}
	log.Println(string(hashed), "hash")

	email, err := profile.ChangePasswordReset(body.Username, string(hashed), conn)

	switch err {
	case errors.UpdateItemError:
		return qs.NewError(err.Error(), 309)
	default:
	}

	err = sendResetPasswordEmail(email, newPass)
	if err != nil {
		return qs.NewError(err.Error(), 109)
	}
	//TODO
	//SEND EMAIL

	res := Response{
		Success: true,
		Message: "password reset successfully",
	}
	return qs.NewResponse(200, res)
}

func createNewPassword() string {
	h := fnv.New64a()
	t := time.Now().String()
	h.Write([]byte(t))
	// h.Write([]byte(username))
	h.Write([]byte("kowalski analysis"))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

func main() {
	lambda.Start(handler)
}
