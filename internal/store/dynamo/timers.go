package dynamo

import (
	"context"
	"fmt"
	"poker"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TimerRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewTimerRepository(client *dynamodb.Client, tableName string) *TimerRepository {
	return &TimerRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *TimerRepository) Timer(ctx context.Context, id string) (*poker.Timer, error) {

	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var timer = new(poker.Timer)

	err = attributevalue.UnmarshalMap(result.Item, timer)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ddb record: %w", err)
	}

	return timer, nil

}

func (r *TimerRepository) TimersByUserID(ctx context.Context, userID string) ([]*poker.Timer, error) {

	emailExpr := expression.Key("UserID").Equal(expression.Value(userID))
	expr, err := expression.NewBuilder().WithKeyCondition(emailExpr).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression for user by email query: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		IndexName:                 aws.String("user-id-index"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var timers = make([]*poker.Timer, 0, len(result.Items))

	err = attributevalue.UnmarshalListOfMaps(result.Items, &timers)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ddb record: %w", err)
	}

	return timers, nil

}

func (r *TimerRepository) SaveTimer(ctx context.Context, timer *poker.Timer) error {

	timer.CreatedAt = time.Now()
	timer.UpdatedAt = time.Now()

	item, err := attributevalue.MarshalMap(timer)
	if err != nil {
		return fmt.Errorf("failed to marshal timer: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	return err

}
func (r *TimerRepository) DeleteTimer(ctx context.Context, id string) error {

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})

	return err

}
