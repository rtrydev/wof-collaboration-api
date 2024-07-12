package commands

import (
	"context"
	"errors"
	"log"

	"github.com/rtrydev/wof-collaboration-api/src/application/interfaces"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration"
	"github.com/rtrydev/wof-collaboration-api/src/domain/schema"
)

type CreateCollaboration struct {
	SchemaId string
	IssuerId string
}

type CreateCollaborationHandler interfaces.CommandHandler[CreateCollaboration, string]

type createCollaborationHandler struct {
	schemaRepository        schema.SchemaRepository
	collaborationRepository collaboration.CollaborationRepository
}

func NewCreateCollaborationHandler(
	schemaRepository schema.SchemaRepository,
	collaborationRepository collaboration.CollaborationRepository,
) CreateCollaborationHandler {
	return createCollaborationHandler{
		schemaRepository:        schemaRepository,
		collaborationRepository: collaborationRepository,
	}
}

func (handler createCollaborationHandler) Handle(ctx context.Context, command CreateCollaboration) (string, error) {
	schema, err := handler.schemaRepository.GetSchema(ctx, command.SchemaId)

	if err != nil {
		return "", errors.New("failed to get schema")
	}

	if schema.OwnerId != command.IssuerId {
		log.Println("The issuer is not owner of the schema!")
		return "", errors.New("issuer not owner")
	}

	_, err = handler.collaborationRepository.GetForSchema(ctx, command.SchemaId)

	if err == nil {
		return "", errors.New("a collaboration for the provided schema already exists")
	}

	collaborationToken, err := handler.collaborationRepository.Create(ctx, command.SchemaId)

	if err != nil {
		log.Println("Could not create collaboration!")
		return "", errors.New("could not create collaboration")
	}

	return collaborationToken, nil
}
