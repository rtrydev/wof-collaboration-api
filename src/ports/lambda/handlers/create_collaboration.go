package handlers

import (
	"context"
	"encoding/json"

	"github.com/rtrydev/wof-collaboration-api/src/application/commands"
	"github.com/rtrydev/wof-collaboration-api/src/service"

	"github.com/aws/aws-lambda-go/events"
)

type CreateCollaborationBody struct {
	SchemaId string `json:"schema_id"`
}

type CreateCollaborationResponse struct {
	Id                 string `json:"id"`
	CollaborationToken string `json:"collaboration_token"`
}

func CreateCollaborationHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body CreateCollaborationBody
	json.Unmarshal([]byte(event.Body), &body)
	issuerId := event.RequestContext.Authorizer["lambda"].(map[string]interface{})["user_id"].(string)

	app := service.NewApplication(ctx)

	collaboration, err := app.Commands.CreateCollaboration.Handle(ctx, commands.CreateCollaboration{
		SchemaId: body.SchemaId,
		IssuerId: issuerId,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
		}, nil
	}

	responseBody, err := json.Marshal(CreateCollaborationResponse{
		Id:                 collaboration.Id,
		CollaborationToken: collaboration.Token,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}
