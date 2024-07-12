package queries

import (
	"context"
	"errors"
	"log"

	"github.com/rtrydev/wof-collaboration-api/src/application/interfaces"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration"
	"github.com/rtrydev/wof-collaboration-api/src/domain/schema"
)

type GetCollaborationForSchema struct {
	SchemaId string
	UserId   string
}

type GetCollaborationForSchemaHandler interfaces.QueryHandler[GetCollaborationForSchema, *collaboration.Collaboration]

type getCollaborationForSchemaHandler struct {
	collaborationRepository collaboration.CollaborationRepository
	schemaRepository        schema.SchemaRepository
}

func NewGetCollaborationForSchemaHandler(
	collaborationRepository collaboration.CollaborationRepository,
	schemaRepository schema.SchemaRepository,
) GetCollaborationForSchemaHandler {
	return getCollaborationForSchemaHandler{
		collaborationRepository: collaborationRepository,
		schemaRepository:        schemaRepository,
	}
}

func (handler getCollaborationForSchemaHandler) Handle(ctx context.Context, query GetCollaborationForSchema) (*collaboration.Collaboration, error) {
	collaboration, err := handler.collaborationRepository.GetForSchema(ctx, query.SchemaId)

	if err != nil {
		log.Println("Could not find collaboration for schema.")
		return nil, errors.New("collaboration not found")
	}

	schema, err := handler.schemaRepository.GetSchema(ctx, query.SchemaId)

	if err != nil {
		log.Println("Could not find schema!")
		return nil, errors.New("schema not found")
	}

	if schema.OwnerId != query.UserId {
		log.Println("The issuer is not owner of the schema!")
		return nil, errors.New("issuer not owner")
	}

	return collaboration, nil
}
