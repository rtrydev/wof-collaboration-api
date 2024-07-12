package schema

import "context"

type SchemaRepository interface {
	GetSchema(ctx context.Context, schemaId string) (*Schema, error)
}
