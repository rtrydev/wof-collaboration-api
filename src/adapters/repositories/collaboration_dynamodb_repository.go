package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/beevik/guid"
)

type CollaborationModel struct {
	Id       string `dynamodbav:"id"`
	SchemaId string `dynamodbav:"schema_id"`
	Token    string `dynamodbav:"token"`
}

type CollaborationDynamoDbRepository struct {
	dynamoDbClient *dynamodb.Client
	tableName      string
}

func NewCollaborationDynamoDbRepository(
	dynamoDbClient *dynamodb.Client,
) CollaborationDynamoDbRepository {
	return CollaborationDynamoDbRepository{
		dynamoDbClient: dynamoDbClient,
		tableName:      "wof-collaborations",
	}
}

func (repository CollaborationDynamoDbRepository) GetById(ctx context.Context, id string) (*collaboration.Collaboration, error) {
	var collaborationModel *CollaborationModel = nil

	collaborationId, err := attributevalue.Marshal(id)

	if err != nil {
		return nil, err
	}

	key := map[string]types.AttributeValue{
		"id": collaborationId,
	}

	response, err := repository.dynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: key, TableName: &repository.tableName,
	})

	if err != nil {
		log.Printf("Could not get collaboration with id %v. Error: %v\n", id, err)
		return nil, err
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &collaborationModel)

		if err != nil {
			log.Printf("Could not unmarshal response. Error: %v\n", err)
			return nil, err
		}
	}

	if collaborationModel == nil {
		return nil, errors.New("schema model is nil")
	}

	collaboration := collaboration.Collaboration{
		Id:       collaborationModel.Id,
		SchemaId: collaborationModel.SchemaId,
		Token:    collaborationModel.Token,
	}

	return &collaboration, err
}

func (repository CollaborationDynamoDbRepository) GetForSchema(ctx context.Context, schemaId string) (*collaboration.Collaboration, error) {
	var collaborationModel *CollaborationModel = nil

	keyCondition := expression.Key("schema_id").Equal(expression.Value(schemaId))
	queryExpression, err := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		Build()

	if err != nil {
		log.Println("Unable to process the provided schemaId")
		return nil, errors.New("unable to build expression")
	}

	response, err := repository.dynamoDbClient.Query(ctx, &dynamodb.QueryInput{
		IndexName:                 aws.String("schema_id-index"),
		ExpressionAttributeNames:  queryExpression.Names(),
		ExpressionAttributeValues: queryExpression.Values(),
		KeyConditionExpression:    queryExpression.KeyCondition(),
		TableName:                 &repository.tableName,
	})

	if err != nil {
		log.Println("Unable to fetch schema from the database", err)
		panic(err)
	}

	if len(response.Items) == 0 {
		return nil, errors.New("item not found")
	}

	err = attributevalue.UnmarshalMap(response.Items[0], &collaborationModel)

	if err != nil {
		log.Println("Unable to unmarshal query result", err)
		return nil, errors.New("unable to parse")
	}

	collaboration := collaboration.Collaboration{
		Id:       collaborationModel.Id,
		SchemaId: collaborationModel.SchemaId,
		Token:    collaborationModel.Token,
	}

	return &collaboration, err
}

func (repository CollaborationDynamoDbRepository) Create(ctx context.Context, schemaId string) (*collaboration.Collaboration, error) {
	collaborationModel := CollaborationModel{
		Id:       guid.NewString(),
		SchemaId: schemaId,
		Token:    guid.NewString(),
	}

	collaborationData, err := attributevalue.MarshalMap(collaborationModel)

	if err != nil {
		log.Println("Could not create collaboration data!", err)
		return nil, errors.New("unparsable collaboration schema")
	}

	_, err = repository.dynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      collaborationData,
		TableName: &repository.tableName,
	})

	if err != nil {
		log.Println("Could not put item in dynamodb!", err)
		return nil, errors.New("database put failed")
	}

	return &collaboration.Collaboration{
		Id:       collaborationModel.Id,
		SchemaId: collaborationModel.SchemaId,
		Token:    collaborationModel.Token,
	}, nil
}
