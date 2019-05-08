package profile

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
)

func TestGetSubscription(t *testing.T) {
	var conn *dynamodb.DynamoDB
	database.SetConn(&conn)
	userId := "94633b7a8c014b57"
	sub, err := GetSubscription(userId, conn)
	fmt.Println(sub, err)
}
