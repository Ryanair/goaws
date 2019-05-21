package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type UpdateOp struct {
	key       Key
	expr      Expression
	tableName string
}

func NewUpdateOp(key Key, expr Expression, tableName string) UpdateOp {
	op := UpdateOp{
		key:       key,
		expr:      expr,
		tableName: tableName,
	}
	return op
}

func (uo *UpdateOp) build() dynamodb.UpdateItemInput {
	dbKey, _ := marshalKey(uo.key)
	expr, _ := uo.expr.Build()
	input := dynamodb.UpdateItemInput{
		Key:                       dbKey,
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 &uo.tableName,
	}
	return input
}
