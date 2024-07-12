package commands

import (
	"context"
	"errors"
	"log"

	"github.com/rtrydev/wof-collaboration-api/src/application/interfaces"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration_affiliation"
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
}

func NewAddUserToCollaborationHandler(
	collaborationRepository collaboration.CollaborationRepository,
	collaborationAffiliationRepository collaboration_affiliation.CollaborationAffiliationRepository,
) AddUserToCollaborationHandler {
	return addUserToCollaborationHandler{
		collaborationRepository:            collaborationRepository,
		collaborationAffiliationRepository: collaborationAffiliationRepository,
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

	affiliation := collaboration_affiliation.CollaborationAffiliation{
		SchemaId: collaboration.SchemaId,
		UserId:   command.UserId,
	}
	err = handler.collaborationAffiliationRepository.Create(ctx, affiliation)

	return nil, err
}
