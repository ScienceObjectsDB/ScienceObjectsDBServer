package authhandler

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/databasehandler"
	"github.com/ScienceObjectsDB/go-api/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const tokenLen = 64

// TokenType Type of the provided token
type TokenType int

const (
	//OAuth2Token Api token from oauth2 provider
	OAuth2Token TokenType = 0
	//UserAPIToken API token
	UserAPIToken TokenType = 1
)

//ExtractedToken An extracted token
type ExtractedToken struct {
	Token     string
	TokenType TokenType
}

//ProjectAuthHandler Simple project based authentication handler
type ProjectAuthHandler struct {
	OAuth2Handler         *OAuth2Handler
	DatabaseTokenHandler  *databasehandler.TokenActionHandler
	ProjectHandler        *databasehandler.ProjectActionHandler
	DatasetHandler        *databasehandler.DatasetActionHandler
	DatasetVersionHandler *databasehandler.DatasetVersionActionHandler
	ObjectGroupHandler    *databasehandler.ObjectGroupHandler
}

func InitProjectHandler(
	projectHandler *databasehandler.ProjectActionHandler,
	tokenHandler *databasehandler.TokenActionHandler,
	datasetHandler *databasehandler.DatasetActionHandler,
	datasetVersionHandler *databasehandler.DatasetVersionActionHandler,
	objectGroupHandler *databasehandler.ObjectGroupHandler) (*ProjectAuthHandler, error) {

	oauth2handler, err := InitOauth2()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &ProjectAuthHandler{
		OAuth2Handler:         oauth2handler,
		DatabaseTokenHandler:  tokenHandler,
		ProjectHandler:        projectHandler,
		DatasetHandler:        datasetHandler,
		DatasetVersionHandler: datasetVersionHandler,
		ObjectGroupHandler:    objectGroupHandler,
	}, nil

}

func (handler *ProjectAuthHandler) UserID(requestContext context.Context) (string, error) {
	requestToken, err := getToken(requestContext)
	if err != nil {
		log.Println(err.Error())
		return "false", err
	}

	userID, err := handler.OAuth2Handler.getUserIDFromOAuth2(requestToken.Token)
	if err != nil {
		log.Println(err.Error())
		return "false", err
	}

	return userID, nil
}

//Authorize Authorizes the request for a resource based on project scoped rights
func (handler *ProjectAuthHandler) Authorize(
	requestContext context.Context,
	resource models.Resource,
	requiredRight models.Right,
	resourceID string) (bool, error) {

	var err error
	var projectID string

	switch resource {
	case models.Resource_Project:
		projectID = resourceID
	case models.Resource_Dataset:
		projectID, err = handler.DatasetHandler.GetDatasetProjectID(resourceID)
	case models.Resource_DatasetVersion:
		projectID, err = handler.getDatasetVersionProjectID(resourceID)
	case models.Resource_DatasetObjectGroupResource:
		projectID, err = handler.getObjectGroupProjectID(resourceID)
	case models.Resource_DatasetObject:
		projectID, err = handler.getObjectProjectID(resourceID)
	default:
		err = fmt.Errorf("Can not process resource type: %v", resource)
		return false, err
	}

	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	requestToken, err := getToken(requestContext)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	userID, err := handler.OAuth2Handler.getUserIDFromOAuth2(requestToken.Token)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	var authorized bool

	switch requestToken.TokenType {
	case OAuth2Token:
		authorized, err = handler.ProjectHandler.UserCanAccessProject(requiredRight, userID, projectID)
	case UserAPIToken:
		authorized, err = handler.DatabaseTokenHandler.ValidateTokenForResourceAction(requestToken.Token, projectID, models.Resource_Project, requiredRight)
	default:
		authorized, err = false, fmt.Errorf("Could not process tokentype")
	}

	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	return authorized, nil
}

func getToken(ctx context.Context) (*ExtractedToken, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	accessToken := meta.Get("AccessToken")
	apiToken := meta.Get("UserAPIToken")

	extractedToken := ExtractedToken{}

	if len(accessToken) > 0 {
		extractedToken.Token = accessToken[0]
		extractedToken.TokenType = OAuth2Token
	} else if len(apiToken) > 0 {
		extractedToken.Token = apiToken[0]
		extractedToken.TokenType = UserAPIToken
	} else {
		return nil, fmt.Errorf("Could not extract auth token, please specify access_token or user_api_token")
	}

	return &extractedToken, nil
}

func (handler *ProjectAuthHandler) getDatasetVersionProjectID(id string) (string, error) {
	datasetID, err := handler.DatasetVersionHandler.GetDatasetVersionDatasetID(id)
	if err != nil {
		log.Println(err.Error())
		return "", nil
	}
	return handler.DatasetHandler.GetDatasetProjectID(datasetID)
}

func (handler *ProjectAuthHandler) getObjectGroupProjectID(id string) (string, error) {
	objectGroup, err := handler.ObjectGroupHandler.GetObjectGroup(id)
	if err != nil {
		log.Println(err.Error())
		return "", nil
	}
	datasetID := objectGroup.GetDatasetID()

	return handler.DatasetHandler.GetDatasetProjectID(datasetID)
}

func (handler *ProjectAuthHandler) getObjectProjectID(id string) (string, error) {
	groupID, _, err := handler.ObjectGroupHandler.GetObject(id)
	if err != nil {
		log.Println(err.Error())
		return "", nil
	}

	objectGroup, err := handler.ObjectGroupHandler.GetObjectGroup(groupID)
	if err != nil {
		log.Println(err.Error())
		return "", nil
	}
	datasetID := objectGroup.GetDatasetID()

	return handler.DatasetHandler.GetDatasetProjectID(datasetID)
}
