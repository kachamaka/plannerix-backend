package notifications

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/database"
)

var conn *dynamodb.DynamoDB

func TestCreateDataBase(t *testing.T) {
	database.SetConn(&conn)
	out, err := conn.CreateTable(&createInput)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", out)
}

var createInput = dynamodb.CreateTableInput{
	TableName: aws.String("plannerix-first-lesson-notifications"),
	AttributeDefinitions: []*dynamodb.AttributeDefinition{
		&dynamodb.AttributeDefinition{
			AttributeName: aws.String("time"),
			AttributeType: aws.String("N"),
		},
		&dynamodb.AttributeDefinition{
			AttributeName: aws.String("user_id"),
			AttributeType: aws.String("S"),
		},
	},
	KeySchema: []*dynamodb.KeySchemaElement{
		&dynamodb.KeySchemaElement{

			AttributeName: aws.String("time"),
			KeyType:       aws.String("HASH"),
		},
		&dynamodb.KeySchemaElement{
			AttributeName: aws.String("user_id"),
			KeyType:       aws.String("RANGE"),
		},
	},
	ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(1),
		WriteCapacityUnits: aws.Int64(1),
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		&dynamodb.GlobalSecondaryIndex{
			IndexName: aws.String("userIDIndex"),
			KeySchema: []*dynamodb.KeySchemaElement{
				&dynamodb.KeySchemaElement{

					AttributeName: aws.String("user_id"),
					KeyType:       aws.String("HASH"),
				},
				&dynamodb.KeySchemaElement{
					AttributeName: aws.String("time"),
					KeyType:       aws.String("RANGE"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("ALL"),
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		},
	},
}
