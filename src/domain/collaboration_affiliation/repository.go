package collaboration_affiliation

import "context"

type CollaborationAffiliationRepository interface {
	GetBySchemaId(ctx context.Context, schemaId string) ([]CollaborationAffiliation, error)
	GetByUserId(ctx context.Context, userId string) ([]CollaborationAffiliation, error)
	Create(ctx context.Context, affiliation CollaborationAffiliation) error
}
