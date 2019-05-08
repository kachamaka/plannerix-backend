package main

import (
	"context"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/notifications"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/schedule"
)

var conn *dynamodb.DynamoDB

type Request struct {
	Token    string             `json:"token"`
	Subjects []schedule.Subject `json:"subjects"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func handler(ctx context.Context, req notifications.NotificationPayload) {
	// 	switch req.Type:
	// case 1:
	switch req.Type {
	case 1:
		// notification for startLesson
		id := req.UserID
	}
	// return qs.NewResponse(200, "hello")
}

func main() {
	lambda.Start(handler)
}
