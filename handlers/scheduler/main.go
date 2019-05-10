package main

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/notifications"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	lclient "github.com/aws/aws-sdk-go/service/lambda"
)

var conn *dynamodb.DynamoDB

func handler(ctx context.Context) {
	database.SetConn(&conn)
	tc := notifications.NewTimeConverter()
	firstLessonSlice, err := notifications.GetUsersInRange(tc.GetTimeInMinutes(), conn)
	if err != nil {
		fmt.Println("Error by getting users in range", err)
		return
	}

	fmt.Println("Scheduler gets time", tc.GetTimeInMinutes(), "first lesson notification slice", firstLessonSlice)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := lclient.New(sess, &aws.Config{Region: aws.String("eu-central-1")})

	for _, fl := range firstLessonSlice {
		go InvokeLambda(fl, client)
	}

}

func InvokeLambda(firstLesson notifications.FirstLessonNotificationItem, client *lclient.Lambda) {
	p := notifications.NotificationPayload{
		UserID: firstLesson.UserID,
		Type:   1,
	}
	payload, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error marshalling firstLesson request", err)
		return
	}
	_, err = notifications.InvokeNotificationFunction(client, payload)
	if err != nil {
		fmt.Println("Error marshalling firstLesson request", err)
		return
	}

}

func main() {
	lambda.Start(handler)
}
