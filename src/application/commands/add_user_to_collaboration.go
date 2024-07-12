package commands

import (
	"context"
	"errors"
	"log"

	"github.com/rtrydev/wof-collaboration-api/src/application/interfaces"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration_affiliation"
	"github.com/rtrydev/wof-collaboration-api/src/domain/schema"
)

type AddUserToCollaboration struct {
	UserId   string
	SchemaId string
	Token    string
}

type AddUserToCollaborationHandler interfaces.CommandHandler[AddUserToCollaboration, any]

type addUserToCollaborationHandler struct {
	collaborationRepository            collaboration.CollaborationRepository
	collaborationAffiliationRepository collaboration_affiliation.CollaborationAffiliationRepository
	schemaRepository                   schema.SchemaRepository
}

func NewAddUserToCollaborationHandler(
	collaborationRepository collaboration.CollaborationRepository,
	collaborationAffiliationRepository collaboration_affiliation.CollaborationAffiliationRepository,
	schemaRepository schema.SchemaRepository,
) AddUserToCollaborationHandler {
	return addUserToCollaborationHandler{
		collaborationRepository:            collaborationRepository,
		collaborationAffiliationRepository: collaborationAffiliationRepository,
		schemaRepository:                   schemaRepository,
	}
}

func (handler addUserToCollaborationHandler) Handle(ctx context.Context, command AddUserToCollaboration) (any, error) {
	collaboration, err := handler.collaborationRepository.GetForSchema(ctx, command.SchemaId)

	if err != nil {
		log.Println("Could not find the collaboration!")
		return nil, errors.New("collaboration not found")
	}

	if collaboration.Token != command.Token {
		log.Println("Invalid token!")
		return nil, errors.New("invalid token")
	}

	schema, err := handler.schemaRepository.GetSchema(ctx, collaboration.SchemaId)

	if err != nil {
		log.Println("Could not find the schema!")
		return nil, errors.New("schema not found")
	}

	if schema.OwnerId == command.UserId {
		log.Println("User cannot join their own collaboration.")
		return nil, errors.New("issuer is owner")
	}

	affiliation := collaboration_affiliation.CollaborationAffiliation{
		SchemaId: collaboration.SchemaId,
		UserId:   command.UserId,
	}
	err = handler.collaborationAffiliationRepository.Create(ctx, affiliation)

	return nil, err
}
