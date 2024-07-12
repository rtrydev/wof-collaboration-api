package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rtrydev/wof-collaboration-api/src/application/queries"
	"github.com/rtrydev/wof-collaboration-api/src/service"
)

type GetUserCollaborationsResponse struct {
	SchemaIds []string `json:"schema_ids"`
}

func GetUserCollaborationsHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	issuerId := event.RequestContext.Authorizer["lambda"].(map[string]interface{})["user_id"].(string)

	app := service.NewApplication(ctx)

	schemaIds, _ := app.Queries.GetUserCollaborations.Handle(ctx, queries.GetUserCollaborations{
		UserId: issuerId,
	})

	responseBody, _ := json.Marshal(GetUserCollaborationsResponse{
		SchemaIds: schemaIds,
	})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}
