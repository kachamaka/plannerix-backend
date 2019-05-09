package main

import (
	"context"
	"encoding/json"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/notifications"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kinghunter58/jwe"

	lclient "github.com/aws/aws-sdk-go/service/lambda"
	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"
)

var conn *dynamodb.DynamoDB

type Request struct {
	Token        string `json:"token"`
	Subscription string `json:"subscription"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)

	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}

	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError("Internal Server Error", 6)
	}

	p := profile.Payload{}
	err = jwe.ParseEncryptedToken(body.Token, key, &p)
	if err != nil {
		return qs.NewError("Internal Server Error", 6) // Fix output
	}
	database.SetConn(&conn)
	err = notifications.UpdateSubscriptionOfUser(p, body.Subscription, conn)
	if err != nil {
		return qs.NewError("Internal Server Error", 6) // Fix output
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := lclient.New(sess, &aws.Config{Region: aws.String("eu-central-1")})
	payload := notifications.NotificationPayload{
		Type:   2,
		Msg:    "Подобни известия ще ти напомнят за началото на учебния ден",
		UserID: p.ID,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return qs.NewError("Internal Server Error", 6) // Fix output
	}
	notifications.InvokeNotificationFunction(client, bytes)
	res := Response{
		Success: true,
		Message: "subscription updated successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
