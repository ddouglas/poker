package dynamo

import (
	"context"
	"fmt"
	"poker"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserRepository(client *dynamodb.Client, tableName string) *UserRepository {
	return &UserRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *UserRepository) User(ctx context.Context, id string) (*poker.User, error) {

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

	var user = new(poker.User)

	err = attributevalue.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ddb record: %w", err)
	}

	return user, nil

}

func (r *UserRepository) UserByEmail(ctx context.Context, email string) (*poker.User, error) {

	emailExpr := expression.Key("Email").Equal(expression.Value(email))
	expr, err := expression.NewBuilder().WithKeyCondition(emailExpr).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression for user by email query: %w", err)
	}

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		Limit:                     aws.Int32(1),
		IndexName:                 aws.String("email-index"),
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

	var user = new(poker.User)

	err = attributevalue.UnmarshalMap(result.Items[0], user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ddb record: %w", err)
	}

	return user, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *poker.User) error {

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	return err

}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})

	return err

}
