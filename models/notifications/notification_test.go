package notifications

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"
)

func TestGetDay(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Sofia")
	currentTime := time.Now().In(location)
	fmt.Println(currentTime.Weekday(), time.Wednesday, int(currentTime.Weekday()))
}

func TestAddUserNotificationItem(t *testing.T) {
	tc := NewTimeConverter()
	minutes := tc.GetTimeInMinutes()
	body := FirstLessonNotificationItem{
		UserID:  "e955028dc4da4d6f",
		Minutes: minutes,
	}
	fmt.Println(body)
	var conn *dynamodb.DynamoDB
	database.SetConn(&conn)
	err := AddUserNotificationItem(body, conn)
	if err != nil {
		t.Error(err)
	}
}

func TestGetUsersInRange(t *testing.T) {
	tc := NewTimeConverter()
	minutes := tc.GetTimeInMinutes()
	var conn *dynamodb.DynamoDB
	database.SetConn(&conn)
	items, err := GetUsersInRange(minutes, conn)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(items)
}

func TestDeleteOldItemForDay(t *testing.T) {
	userID := "e955028dc4da4d6f"
	day := time.Wednesday
	var conn *dynamodb.DynamoDB
	database.SetConn(&conn)
	err := DeleteOldItemForDay(day, userID, conn)
	fmt.Println(err)
}

func TestUpdateSubscriptionOfUser(t *testing.T) {
	sub := `{"endpoint":"https://fcm.googleapis.com/fcm/send/fkyq87l516E:APA91bGtSoErGzpvWCbpIlJ7rYRYW6StF_lWm5ngRytNQj7S1O3vn054Pw5PsLLEajDDvy_vdXOJWk-t1nMDs4H2cBpBhdinvRVM5WmNw_OeLfg4bptP28IQCgGmhzyP48Uy28ZeSDfL","expirationTime":null,"keys":{"p256dh":"BOagEacJ9sid0N67W9MXRSdpW4ibsjc8zETDsW_V_pyEpbMT22JIjSrNAdipxYwMkg8w7XMyJM2nMVRl5_20EKk","auth":"mnGQWyFMkh9reNxi3fVYuw"}}`
	user := profile.Payload{
		Username: "testingWE4",
	}
	var conn *dynamodb.DynamoDB
	database.SetConn(&conn)
	err := UpdateSubscriptionOfUser(user, sub, conn)
	fmt.Println(err)
}
