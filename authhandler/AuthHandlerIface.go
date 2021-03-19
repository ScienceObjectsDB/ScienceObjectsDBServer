package authhandler

import (
	"context"

	"github.com/ScienceObjectsDB/go-api/models"
)

//AuthHandler Interface for the authentication handler
type AuthHandler interface {
	Authorize(requestContext context.Context, resource models.Resource, requiredRight models.Right, resourceID string) (bool, error)
	UserID(requestContext context.Context) (string, error)
}
