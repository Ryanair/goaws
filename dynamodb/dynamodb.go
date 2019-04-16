package dynamodb

import (
	"github.com/Ryanair/goaws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pkg/errors"
)

type Client struct {
	db *dynamodb.DynamoDB
}

func NewClient(cfg *goaws.Config, options ...func(*dynamodb.DynamoDB)) *Client {
	db := dynamodb.New(cfg.Provider)
	for _, opt := range options {
		opt(db)
	}

	return &Client{db: db}
}

func Endpoint(endpoint string) func(*dynamodb.DynamoDB) {
	return func(db *dynamodb.DynamoDB) {
		db.Endpoint = endpoint
	}
}

func (c *Client) Put(item interface{}, tableName string) error {
	av, err := marshalItem(item)
	if err != nil {
		return newErr(err, MarshalErrCode, "put item marshal failed")
	}

	input := dynamodb.PutItemInput{
		Item:      av,
		TableName: &tableName,
	}
	if _, err := c.db.PutItem(&input); err != nil {
		return newOpsErr(err, "put item failed")
	}

	return nil
}

func (c *Client) PutWithCondition(item interface{}, conditionBuilder expression.ConditionBuilder, tableName string) error {
	exp, err := expression.NewBuilder().WithCondition(conditionBuilder).Build()
	if err != nil {
		return newErr(err, InvalidConditionErrCode, "invalid put condition")
	}

	av, err := marshalItem(item)
	if err != nil {
		return newErr(err, MarshalErrCode, "marshal put item with condition failed")
	}

	input := dynamodb.PutItemInput{
		Item:                      av,
		ConditionExpression:       exp.Condition(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		TableName:                 &tableName,
	}
	if _, err := c.db.PutItem(&input); err != nil {
		return newOpsErr(err, "put item with condition failed")
	}

	return nil
}

type Key struct {
	partitionName  string
	partitionValue string
	sortName       *string
	sortValue      *string
}

func NewPartitionKey(name, value string) Key {
	return Key{
		partitionName:  name,
		partitionValue: value,
	}
}

func NewPartitionAndSortKey(partitionName, partitionValue, sortName, sortValue string) Key {
	return Key{
		partitionName:  partitionName,
		partitionValue: partitionValue,
		sortName:       &sortName,
		sortValue:      &sortValue,
	}
}

func (c *Client) Get(key Key, consistentRead bool, tableName string, out interface{}) (bool, error) {
	dbKey, err := marshalKey(key)
	if err != nil {
		return false, newErr(err, MarshalErrCode, "marshal key failed")
	}

	input := dynamodb.GetItemInput{
		Key:            dbKey,
		TableName:      &tableName,
		ConsistentRead: &consistentRead,
	}
	output, getErr := c.db.GetItem(&input)
	if getErr != nil {
		return false, newOpsErr(err, "get item failed")
	}

	if unmarshalErr := dynamodbattribute.UnmarshalMap(output.Item, &out); unmarshalErr != nil {
		return false, newErr(unmarshalErr, UnmarshalErrCode, "unmarshal GetOutput failed")
	}

	if len(output.Item) == 0 {
		return false, nil
	}

	return true, nil
}

func marshalItem(item interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil, errors.Wrap(err, "marshal item failed")
	}

	return av, nil
}

func marshalKey(key Key) (map[string]*dynamodb.AttributeValue, error) {
	keys := map[string]*dynamodb.AttributeValue{}
	partition, err := dynamodbattribute.Marshal(key.partitionValue)
	if err != nil {
		return nil, errors.Wrap(err, "marshal partition key failed")
	}

	keys[key.partitionName] = partition
	if key.sortName != nil && key.sortValue != nil {
		sort, err := dynamodbattribute.Marshal(key.sortValue)
		if err != nil {
			return nil, errors.Wrap(err, "marshal sort key failed")
		}
		keys[*key.sortName] = sort
	}

	return keys, nil
}
