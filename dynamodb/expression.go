package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Expression struct {
	updateBuilder     *expression.UpdateBuilder
	condBuilder       *expression.ConditionBuilder
	projectionBuilder *expression.ProjectionBuilder
	filterBuilder     *expression.ConditionBuilder
}

func NewExpression(options ...func(*Expression)) Expression {
	var expr Expression
	for _, opt := range options {
		opt(&expr)
	}
	return expr
}

func (e *Expression) Build() (expression.Expression, error) {
	return createExpression(e)
}

func UpdateExpression(builder expression.UpdateBuilder) func(*Expression) {
	return func(expr *Expression) {
		expr.updateBuilder = &builder
	}
}

func ConditionExpression(builder expression.ConditionBuilder) func(*Expression) {
	return func(expr *Expression) {
		expr.condBuilder = &builder
	}
}

func ProjectionExpression(builder expression.ProjectionBuilder) func(*Expression) {
	return func(expr *Expression) {
		expr.projectionBuilder = &builder
	}
}

func FilterExpression(builder expression.ConditionBuilder) func(*Expression) {
	return func(expr *Expression) {
		expr.filterBuilder = &builder
	}
}

func createExpression(expr *Expression) (expression.Expression, error) {
	exprBuilder := expression.NewBuilder()
	switch {
	case expr.updateBuilder != nil:
		exprBuilder = exprBuilder.WithUpdate(*expr.updateBuilder)
	case expr.condBuilder != nil:
		exprBuilder = exprBuilder.WithCondition(*expr.condBuilder)
	case expr.projectionBuilder != nil:
		exprBuilder = exprBuilder.WithProjection(*expr.projectionBuilder)
	case expr.filterBuilder != nil:
		exprBuilder = exprBuilder.WithFilter(*expr.filterBuilder)
	}
	return exprBuilder.Build()
}
