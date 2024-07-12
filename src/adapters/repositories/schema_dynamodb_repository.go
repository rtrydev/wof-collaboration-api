package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/rtrydev/wof-collaboration-api/src/domain/schema"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SchemaModel struct {
	Id      string `dynamodbav:"id"`
	Name    string `dynamodbav:"name"`
	OwnerId string `dynamodbav:"owner_id"`
}

type SchemaDynamoDbRepository struct {
	dynamoDbClient *dynamodb.Client
	tableName      string
}

func NewSchemaDynamoDbRepository(dynamoDbClient *dynamodb.Client) SchemaDynamoDbRepository {
	return SchemaDynamoDbRepository{
		dynamoDbClient: dynamoDbClient,
		tableName:      "wof-schemas",
	}
}

func (repository SchemaDynamoDbRepository) GetSchema(ctx context.Context, schemaId string) (*schema.Schema, error) {
	var schemaModel *SchemaModel = nil
	id, err := attributevalue.Marshal(schemaId)

	if err != nil {
		return nil, err
	}

	key := map[string]types.AttributeValue{
		"id": id,
	}

	response, err := repository.dynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: key, TableName: &repository.tableName,
	})

	if err != nil {
		log.Printf("Could not get schema with id %v. Error: %v\n", schemaId, err)
		return nil, err
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &schemaModel)

		if err != nil {
			log.Printf("Could not unmarshal response. Error: %v\n", err)
			return nil, err
		}
	}

	if schemaModel == nil {
		return nil, errors.New("schema model is nil")
	}

	schema := schema.Schema{
		Id:      schemaModel.Id,
		Name:    schemaModel.Name,
		OwnerId: schemaModel.OwnerId,
	}

	return &schema, err
}
