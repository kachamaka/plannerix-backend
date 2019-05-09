package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
)

var conn *dynamodb.DynamoDB

type Request struct {
	Email string `json:"email"`
}

type Response struct {
	Success bool `json:"success"`
}

// func SendSimpleMessageHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := appengine.NewContext(r)
// 	httpc := urlfetch.Client(ctx)

// 	mg := mailgun.NewMailgun(
// 		"sandboxd15ea31d048e4f92b7225d795260ccb5.mailgun.org", // Domain name
// 		"57d452adb41fbed21081e953187a2de7-3fb021d1-6915ac04",  // API Key
// 	)
// 	mg.SetClient(httpc)

// 	msg, id, err := mg.Send(ctx, mg.NewMessage(
// 		/* From */ "Excited User <mailgun@sandboxd15ea31d048e4f92b7225d795260ccb5.mailgun.org>",
// 		/* Subject */ "Hello",
// 		/* Body */ "Testing some Mailgun awesomness!",
// 		/* To */ "bar@example.com", "martilevski1@abv.bg",
// 	))
// 	if err != nil {
// 		msg := fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg)
// 		http.Error(w, msg, http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write([]byte("Message sent!"))
// }

const (
	apiKey       = "e983b6d93edcd0946bc13b41ece79c37-us20"
	mandrusanKey = "D9IshA7QNDnbpNzgOaywSA"
)

func sendEmail(email string) error {

	from := mail.NewEmail("Plannerix Support", "support@plannerix.eu")
	subject := "Plannerix Account Activation"
	to := mail.NewEmail("", email)
	plainTextContent := "Hello world"
	// htmlContent := "<strong>You can verify your Plannerix account by clicking the following link: https://plannerix.eu/</strong>"
	htmlContent := "<strong>You can verify your Plannerix account by clicking the following link: https://plannerix.eu/</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient("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")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return nil

	// Create an instance of the Mailgun Client
	// mg := mailgun.NewMailgun("sandboxd15ea31d048e4f92b7225d795260ccb5.mailgun.org", "57d452adb41fbed21081e953187a2de7-3fb021d1-6915ac04")

	// sender := "plannerix.noreply@gmail.com"
	// subject := "TEST subject!"
	// body := "Hello from Mailgun Go!"
	// recipient := email

	// message := mg.NewMessage(sender, subject, body, recipient)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	// resp, id, err := mg.Send(ctx, message)

	// if err != nil {
	// 	return err
	// }

	// log.Println(id, resp)
	// return nil
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

	err = sendEmail(body.Email)
	log.Println("email err", err)
	// if err != nil {
	// 	err = profile.DeleteProfile(body.Username, conn)
	// 	if err != nil {
	// 		return qs.NewError(err.Error(), 99)
	// 	}
	// 	return qs.NewError(errors.DoesNotExist("Този имейл").Error(), 106)
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
