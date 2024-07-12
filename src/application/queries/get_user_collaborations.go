package queries

import (
	"context"
	"errors"
	"log"

	"github.com/rtrydev/wof-collaboration-api/src/application/interfaces"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration_affiliation"
)

type GetUserCollaborations struct {
	UserId string
}

type GetUserCollaborationsHandler interfaces.QueryHandler[GetUserCollaborations, []string]

type getUserCollaborationsHandler struct {
	collaborationAffiliationRepository collaboration_affiliation.CollaborationAffiliationRepository
}

func NewGetUserCollaborationsHandler(
	collaborationAffiliationRepository collaboration_affiliation.CollaborationAffiliationRepository,
) GetUserCollaborationsHandler {
	return getUserCollaborationsHandler{
		collaborationAffiliationRepository: collaborationAffiliationRepository,
	}
}

func (handler getUserCollaborationsHandler) Handle(
	ctx context.Context,
	query GetUserCollaborations,
) ([]string, error) {
	affiliations, err := handler.collaborationAffiliationRepository.GetByUserId(ctx, query.UserId)

	if err != nil {
		log.Println("Could not get affiliations for the user")
		return []string{}, errors.New("could not get affiliations")
	}

	schemaIds := []string{}
	for _, affiliation := range affiliations {
		schemaIds = append(schemaIds, affiliation.SchemaId)
	}

	return schemaIds, nil
}
