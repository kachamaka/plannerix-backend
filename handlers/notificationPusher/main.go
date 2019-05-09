package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func handler(ctx context.Context, req notifications.NotificationPayload) {
	database.SetConn(&conn)

	switch req.Type {
	case 1:
		// notification for startLesson
		id := req.UserID
		location, _ := time.LoadLocation("Europe/Sofia")
		currentTime := time.Now().In(location)
		lesson, err := schedule.GetFirstLessonForDay(id, currentTime.Weekday(), conn)
		if err != nil {
			fmt.Println("UserID", id, "Error:", err)
			return
		}
		subscription, err := profile.GetSubscription(id, conn)
		if err != nil {
			fmt.Println("UserID", id, "Error:", err)
			return
		}

		s := &webpush.Subscription{}
		privateKey := "zXrRUAuy3Z2_2635hUmbhwZGGQKLWXW1eBvggLXUpR4"
		publicKey := "BEhfpAWKg9VR38PWRMIw0bsviVSqfaK_48gt_LeTll-EbxGM2qv7m9wp3VCK1oyCCmTj7XfD2mi4vhN_J0Lsx_8"
		// _ = privateKey
		json.Unmarshal([]byte(subscription), s)
		bytes, err := json.Marshal(lesson)
		if err != nil {
			fmt.Println("UserID", id, "Error:", err)
			return
		}
		_, err = webpush.SendNotification(bytes, s, &webpush.Options{
			VAPIDPrivateKey: privateKey,
			VAPIDPublicKey:  publicKey,
			TTL:             30,
			Subscriber:      "traqn02@gmail.com",
		})
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
