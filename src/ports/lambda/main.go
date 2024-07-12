package main

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rtrydev/wof-collaboration-api/src/ports/lambda/handlers"
)

func main() {
	handlerPathElements := strings.Split(os.Getenv("_HANDLER"), "/")

	if len(handlerPathElements) == 0 {
		panic("invalid handler path")
	}

	handlerName := handlerPathElements[len(handlerPathElements)-1]

	switch handlerName {
	case "createCollaboration":
		lambda.Start(handlers.CreateCollaborationHandler)
	case "joinCollaboration":
		lambda.Start(handlers.JoinCollaborationHandler)
	case "getUserCollaborations":
		lambda.Start(handlers.GetUserCollaborationsHandler)
	case "getCollaborationForSchema":
		lambda.Start(handlers.GetCollaborationForSchema)
	default:
		panic("unsupported handler")
	}
}
