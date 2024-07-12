package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration_affiliation"
)

type CollaborationAffiliationModel struct {
	SchemaId string `dynamodbav:"schema_id"`
	UserId   string `dynamodbav:"user_id"`
}

type CollaborationAffiliationDynamoDbRepository struct {
	dynamoDbClient *dynamodb.Client
	tableName      string
}

func NewCollaborationAffiliationDynamoDbRepository(
	dynamoDbClient *dynamodb.Client,
) CollaborationAffiliationDynamoDbRepository {
	return CollaborationAffiliationDynamoDbRepository{
		dynamoDbClient: dynamoDbClient,
		tableName:      "wof-collaboration-affiliations",
	}
}

func (repository CollaborationAffiliationDynamoDbRepository) GetBySchemaId(
	ctx context.Context,
	schemaId string,
) ([]collaboration_affiliation.CollaborationAffiliation, error) {
	var affiliationModels []CollaborationAffiliationModel

	keyCondition := expression.Key("schema_id").Equal(expression.Value(schemaId))
	queryExpression, err := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		Build()

	if err != nil {
		log.Println("Unable to process the provided schemaId")
		return nil, errors.New("unable to build expression")
	}

	response, err := repository.dynamoDbClient.Query(ctx, &dynamodb.QueryInput{
		ExpressionAttributeNames:  queryExpression.Names(),
		ExpressionAttributeValues: queryExpression.Values(),
		KeyConditionExpression:    queryExpression.KeyCondition(),
		TableName:                 &repository.tableName,
	})

	if err != nil {
		log.Println("Unable to fetch affiliation from the database", err)
		panic(err)
	}

	if len(response.Items) == 0 {
		return nil, errors.New("item not found")
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &affiliationModels)

	if err != nil {
		log.Println("Unable to unmarshal query result", err)
		return nil, errors.New("unable to parse")
	}

	result := []collaboration_affiliation.CollaborationAffiliation{}

	for _, model := range affiliationModels {
		result = append(result, collaboration_affiliation.CollaborationAffiliation{
			SchemaId: model.SchemaId,
			UserId:   model.UserId,
		})
	}

	return result, nil
}

func (repository CollaborationAffiliationDynamoDbRepository) GetByUserId(
	ctx context.Context,
	userId string,
) ([]collaboration_affiliation.CollaborationAffiliation, error) {
	var affiliationModels []CollaborationAffiliationModel

	keyCondition := expression.Key("user_id").Equal(expression.Value(userId))
	queryExpression, err := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		Build()

	if err != nil {
		log.Println("Unable to process the provided userId")
		return nil, errors.New("unable to build expression")
	}

	response, err := repository.dynamoDbClient.Query(ctx, &dynamodb.QueryInput{
		IndexName:                 aws.String("user_id-index"),
		ExpressionAttributeNames:  queryExpression.Names(),
		ExpressionAttributeValues: queryExpression.Values(),
		KeyConditionExpression:    queryExpression.KeyCondition(),
		TableName:                 &repository.tableName,
	})

	if err != nil {
		log.Println("Unable to fetch affiliation from the database", err)
		panic(err)
	}

	if len(response.Items) == 0 {
		return nil, errors.New("item not found")
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &affiliationModels)

	if err != nil {
		log.Println("Unable to unmarshal query result", err)
		return nil, errors.New("unable to parse")
	}

	result := []collaboration_affiliation.CollaborationAffiliation{}

	for _, model := range affiliationModels {
		result = append(result, collaboration_affiliation.CollaborationAffiliation{
			SchemaId: model.SchemaId,
			UserId:   model.UserId,
		})
	}

	return result, nil
}

func (repository CollaborationAffiliationDynamoDbRepository) Create(
	ctx context.Context,
	affiliation collaboration_affiliation.CollaborationAffiliation,
) error {
	affiliationModel := CollaborationAffiliationModel{
		SchemaId: affiliation.SchemaId,
		UserId:   affiliation.UserId,
	}

	affiliationData, err := attributevalue.MarshalMap(affiliationModel)

	if err != nil {
		log.Println("Could not create affiliation data!", err)
		return errors.New("unparsable affiliation schema")
	}

	_, err = repository.dynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      affiliationData,
		TableName: &repository.tableName,
	})

	if err != nil {
		log.Println("Could not put item in dynamodb!", err)
		return errors.New("database put failed")
	}

	return nil
}
