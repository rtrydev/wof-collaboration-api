package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rtrydev/wof-collaboration-api/src/application/queries"
	"github.com/rtrydev/wof-collaboration-api/src/service"
)

type CollaborationBody struct {
	Id       string `json:"id"`
	SchemaId string `json:"schema_id"`
	Token    string `json:"token"`
}

type GetCollaborationForSchemaResponse struct {
	Collaboration CollaborationBody `json:"collaboration"`
}

func GetCollaborationForSchema(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	issuerId := event.RequestContext.Authorizer["lambda"].(map[string]interface{})["user_id"].(string)
	schemaId := event.PathParameters["schema_id"]

	app := service.NewApplication(ctx)

	collaboration, err := app.Queries.GetCollaborationForSchema.Handle(ctx, queries.GetCollaborationForSchema{
		SchemaId: schemaId,
		UserId:   issuerId,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	}

	responseBody, _ := json.Marshal(GetCollaborationForSchemaResponse{
		Collaboration: CollaborationBody{
			Id:       collaboration.Id,
			SchemaId: collaboration.SchemaId,
			Token:    collaboration.Token,
		},
	})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}
