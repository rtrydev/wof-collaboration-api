package collaboration

import "context"

type CollaborationRepository interface {
	GetById(ctx context.Context, id string) (*Collaboration, error)
	GetForSchema(ctx context.Context, schemaId string) (*Collaboration, error)
	Create(ctx context.Context, schemaId string) (string, error)
}
