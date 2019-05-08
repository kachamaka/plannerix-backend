package database

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gitlab.com/zapochvam-ei-sq/s-org-backend/models/database"
)

func TestUpdateUsersTable(t *testing.T) {
	updateTable := dynamodb.UpdateTableInput{
		TableName: aws.String("s-org-users"),
		GlobalSecondaryIndexUpdates: []*dynamodb.GlobalSecondaryIndexUpdate{
			&dynamodb.GlobalSecondaryIndexUpdate{
				Create: &dynamodb.CreateGlobalSecondaryIndexAction{
					IndexName: aws.String("idIndex"),
					KeySchema: []*dynamodb.KeySchemaElement{
						&dynamodb.KeySchemaElement{
							AttributeName: aws.String("id"),
							KeyType:       aws.String("HASH"),
						},
					},
					ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(1),
						WriteCapacityUnits: aws.Int64(1),
					},
					Projection: &dynamodb.Projection{
						ProjectionType: aws.String("ALL"),
					},
				},
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("username"),
				AttributeType: aws.String("S"),
			},
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
	}
	var conn *dynamodb.DynamoDB
	database.SetConn(&conn)
	out, err := conn.UpdateTable(&updateTable)
	fmt.Println(err, out)
}
