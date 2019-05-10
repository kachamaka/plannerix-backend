package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/notifications"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/schedule"
)

var conn *dynamodb.DynamoDB
var (
	privateKey = "zXrRUAuy3Z2_2635hUmbhwZGGQKLWXW1eBvggLXUpR4"
	publicKey  = "BEhfpAWKg9VR38PWRMIw0bsviVSqfaK_48gt_LeTll-EbxGM2qv7m9wp3VCK1oyCCmTj7XfD2mi4vhN_J0Lsx_8"
)

func handler(ctx context.Context, req notifications.NotificationPayload) {
	database.SetConn(&conn)
	id := req.UserID
	subscription, err := profile.GetSubscription(id, conn)
	if err != nil {
		fmt.Println("UserID", id, "Error:", err)
		return
	}
	s := &webpush.Subscription{}

	err = json.Unmarshal([]byte(subscription), s)
	if err != nil {
		fmt.Println("UserID", id, "Error by unmarshal:", err)
		return
	}

	switch req.Type {
	case 1:
		// notification for startLesson
		lesson, err := getFirstLessonOfUser(id, conn)
		if err != nil {
			fmt.Println("UserID", id, "Error:", err)
			return
		}

		bytes, err := json.Marshal(lesson)
		if err != nil {
			fmt.Println("UserID", id, "Error:", err)
			return
		}
		err = sendNotification(bytes, s)
		if err != nil {
			fmt.Println("UserID", id, "Error:", err)
			return
		}
	case 2:
		err := sendNotification([]byte(req.Msg), s)
		if err != nil {
			fmt.Println("UserID", id, "Error:", err)
			return
		}
	}

	// return qs.NewResponse(200, "hello")
}

func main() {
	lambda.Start(handler)
}

func sendNotification(bytes []byte, s *webpush.Subscription) error {
	res, err := webpush.SendNotification(bytes, s, &webpush.Options{
		VAPIDPrivateKey: privateKey,
		VAPIDPublicKey:  publicKey,
		TTL:             30,
		Subscriber:      "traqn02@gmail.com",
	})
	b, err := ioutil.ReadAll(res.Body)
	fmt.Println(res.StatusCode, string(b), err)
	return err
}

func getFirstLessonOfUser(id string, conn *dynamodb.DynamoDB) (schedule.Lesson, error) {
	location, _ := time.LoadLocation("Europe/Sofia")
	currentTime := time.Now().In(location)
	lesson, err := schedule.GetFirstLessonForDay(id, currentTime.Weekday(), conn)
	return lesson, err
}
