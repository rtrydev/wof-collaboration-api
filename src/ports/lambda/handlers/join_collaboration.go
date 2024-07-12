package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rtrydev/wof-collaboration-api/src/application/commands"
	"github.com/rtrydev/wof-collaboration-api/src/service"
)

func JoinCollaborationHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	issuerId := event.RequestContext.Authorizer["lambda"].(map[string]interface{})["user_id"].(string)
	schemaId := event.PathParameters["schema_id"]
	token := event.PathParameters["token"]

	app := service.NewApplication(ctx)

	_, err := app.Commands.AddUserToCollaboration.Handle(ctx, commands.AddUserToCollaboration{
		UserId:   issuerId,
		SchemaId: schemaId,
		Token:    token,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
