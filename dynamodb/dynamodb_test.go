// +build local ci

package dynamodb

import (
	"testing"

	"github.com/Ryanair/goaws"
	"github.com/Ryanair/goaws/docker"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/ory/dockertest"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

var cli *Client

const tableName = "db_integration_test"

type TestStruct struct {
	ID     string `dynamodbav:"id"`
	Artist string `dynamodbav:"artist"`
}

func TestDynamoDBClient_Put_ok(t *testing.T) {
	// given
	item := TestStruct{
		ID:     xid.New().String(),
		Artist: "ABBA",
	}

	// when
	err := cli.Put(item, tableName)

	// then
	assert.Nil(t, err)
}

func TestDynamoDBClient_Put_validationException(t *testing.T) {
	// given
	item := struct {
		ID     string      `dynamodbav:"id"`
		Artist chan string `dynamodbav:"artist"`
	}{
		ID:     xid.New().String(),
		Artist: nil,
	}

	// when
	putErr := cli.Put(item, tableName)

	// then
	isValidationFailed := func(err error) bool {
		type validationFailed interface {
			ValidationFailed() bool
		}
		e, ok := err.(validationFailed)
		return ok && e.ValidationFailed()
	}

	assert.True(t, isValidationFailed(putErr))
	containsErr(t, putErr, errors.New("put item failed: ValidationException: "+
		"Supplied AttributeValue is empty, must contain exactly one of the supported datatypes"))
}

func TestDynamoDBClient_PutWithCondition_ok(t *testing.T) {
	// given
	item := TestStruct{
		ID:     xid.New().String(),
		Artist: "Beatles",
	}
	condition := expression.AttributeNotExists(expression.Name("artist"))

	// when
	err := cli.PutWithCondition(item, condition, tableName)

	// then
	assert.Nil(t, err)
}

func TestDynamoDBClient_PutWithCondition_attributeAlreadyExists(t *testing.T) {
	// given
	item := TestStruct{
		ID:     xid.New().String(),
		Artist: "Beatles",
	}
	condition := expression.AttributeNotExists(expression.Name("artist"))
	if err := cli.Put(item, tableName); err != nil {
		t.Fatalf("test %s failed due to %v", t.Name(), err)
	}

	// when
	putErr := cli.PutWithCondition(item, condition, tableName)

	// then
	isConditionFailed := func(err error) bool {
		type conditionFailed interface {
			ConditionFailed() bool
		}
		e, ok := err.(conditionFailed)
		return ok && e.ConditionFailed()
	}

	assert.True(t, isConditionFailed(putErr))
	containsErr(t, putErr, errors.New("put item with condition failed: ConditionalCheckFailedException: The conditional request failed"))
}

func TestDynamoDBClient_Get_ok(t *testing.T) {
	// given
	id := xid.New().String()
	item := TestStruct{
		ID:     id,
		Artist: "ABBA",
	}
	if err := cli.Put(item, tableName); err != nil {
		t.Fatalf("test %s failed due to %v", t.Name(), err)
	}
	key := NewPartitionKey("id", id)

	// when
	out := &TestStruct{}
	ok, err := cli.Get(key, false, tableName, out)

	// then
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestDynamoDBClient_Get_itemNotExists(t *testing.T) {
	// given
	key := NewPartitionKey("id", xid.New().String())

	// when
	out := &TestStruct{}
	ok, err := cli.Get(key, false, tableName, out)

	// then
	assert.False(t, ok)
	assert.Nil(t, err)
}

func getTableDefinition(name string) dynamodb.CreateTableInput {
	return dynamodb.CreateTableInput{
		TableName: aws.String(name),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}
}

func containsErr(t *testing.T, origErr, want error) bool {
	return assert.Contains(t, origErr.Error(), want.Error())
}

func TestMain(m *testing.M) {
	img := docker.Image{
		Repo: "amazon/dynamodb-local",
		Tag:  "latest",
		Env:  nil,
	}

	setup := func(resource *dockertest.Resource) error {
		config, err := goaws.NewConfig(
			goaws.Region(endpoints.EuWest1RegionID),
			goaws.Credentials("secret_id", "secret_key", "random_token"))
		if err != nil {
			return errors.Wrap(err, "couldn't create config")
		}

		cli = NewClient(config, Endpoint("http://localhost:"+resource.GetPort("8000/tcp")))

		tableDef := getTableDefinition(tableName)
		if _, err = cli.db.CreateTable(&tableDef); err != nil {
			return errors.Wrap(err, "could not create table")
		}

		return nil
	}

	docker.Setup(m, img, setup)
}
