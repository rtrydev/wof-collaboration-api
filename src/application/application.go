package application

import (
	"github.com/rtrydev/wof-collaboration-api/src/application/commands"
	"github.com/rtrydev/wof-collaboration-api/src/application/queries"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateCollaboration    commands.CreateCollaborationHandler
	AddUserToCollaboration commands.AddUserToCollaborationHandler
}

type Queries struct {
	GetUserCollaborations queries.GetUserCollaborationsHandler
}
