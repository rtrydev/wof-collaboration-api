package service

import (
	"context"

	"github.com/rtrydev/wof-collaboration-api/src/adapters/repositories"
	"github.com/rtrydev/wof-collaboration-api/src/application"
	"github.com/rtrydev/wof-collaboration-api/src/application/commands"
	"github.com/rtrydev/wof-collaboration-api/src/application/queries"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration"
	"github.com/rtrydev/wof-collaboration-api/src/domain/collaboration_affiliation"
	"github.com/rtrydev/wof-collaboration-api/src/domain/schema"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewApplication(ctx context.Context) application.Application {
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-central-1"))

	if err != nil {
		panic(err)
	}

	schemaRepository := repositories.NewSchemaDynamoDbRepository(dynamodb.NewFromConfig(awsConfig))
	collaborationRepository := repositories.NewCollaborationDynamoDbRepository(dynamodb.NewFromConfig(awsConfig))
	collaborationAffiliationRepository := repositories.NewCollaborationAffiliationDynamoDbRepository(dynamodb.NewFromConfig(awsConfig))

	return newApplication(
		schemaRepository,
		collaborationRepository,
		collaborationAffiliationRepository,
	)
}

func newApplication(
	schemaRepository schema.SchemaRepository,
	collaborationRepository collaboration.CollaborationRepository,
	collaborationAffiliationRepository collaboration_affiliation.CollaborationAffiliationRepository,
) application.Application {
	return application.Application{
		Commands: application.Commands{
			CreateCollaboration:    commands.NewCreateCollaborationHandler(schemaRepository, collaborationRepository),
			AddUserToCollaboration: commands.NewAddUserToCollaborationHandler(collaborationRepository, collaborationAffiliationRepository),
		},
		Queries: application.Queries{
			GetUserCollaborations:     queries.NewGetUserCollaborationsHandler(collaborationAffiliationRepository),
			GetCollaborationForSchema: queries.NewGetCollaborationForSchemaHandler(collaborationRepository, schemaRepository),
		},
	}
}
