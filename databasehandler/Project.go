package databasehandler

import (
	"errors"
	"fmt"
	"log"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//ProjectActionHandler Handler for project related database functions
type ProjectActionHandler struct {
	*DBUtilsHandler
}

// CreateProject Creates a new project and returns the inserted entry
func (handler *ProjectActionHandler) CreateProject(userid string, request *services.CreateProjectRequest) (*models.ProjectEntry, error) {
	projectUser := models.User{
		UserID:   userid,
		Resource: models.Resource_Project,
		Rights:   []models.Right{models.Right_Read, models.Right_Write},
	}

	uuidString := uuid.New().String()

	project := models.ProjectEntry{
		ID:          uuidString,
		Description: request.GetDescription(),
		ProjectName: request.GetName(),
		Users:       []*models.User{&projectUser},
	}

	insertResults, err := handler.GetProjectCollection().InsertOne(handler.MongoDefaultContext, &project)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var oid primitive.ObjectID
	var ok bool

	if oid, ok = insertResults.InsertedID.(primitive.ObjectID); !ok {
		return nil, errors.New("Error decoding result id")
	}

	result := handler.GetProjectCollection().FindOne(handler.MongoDefaultContext, bson.M{
		"_id": oid,
	})

	if result.Err() != nil {
		log.Println(result.Err().Error())
		return nil, err
	}

	insertedProject := models.ProjectEntry{
		Users: []*models.User{},
	}

	err = result.Decode(&insertedProject)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &insertedProject, nil
}

// AddUserToProject Add a user to a project
func (handler *ProjectActionHandler) AddUserToProject(userID string, projectID string, rights []models.Right) (*models.ProjectEntry, error) {
	user := models.User{
		UserID:   userID,
		Resource: models.Resource_Project,
		Rights:   rights,
	}

	_, err := handler.GetProjectCollection().UpdateOne(handler.MongoDefaultContext,
		bson.M{"ID": projectID},
		bson.M{"$addToSet": bson.M{"Users": &user}},
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	project, err := handler.GetProject(projectID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return project, nil
}

// GetUserProjects Returns all projects of a user
func (handler *ProjectActionHandler) GetUserProjects(userID string) ([]*models.ProjectEntry, error) {
	queryResults, err := handler.GetProjectCollection().Find(handler.MongoDefaultContext, bson.M{"Users.UserID": userID})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var projects []*models.ProjectEntry

	for queryResults.Next(handler.MongoDefaultContext) {
		project := models.ProjectEntry{}

		err := queryResults.Decode(&project)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		projects = append(projects, &project)
	}

	return projects, nil
}

// GetProject Returns the project with the given project ID
func (handler *ProjectActionHandler) GetProject(projectID string) (*models.ProjectEntry, error) {
	projectQueryResult := handler.GetProjectCollection().FindOne(handler.MongoDefaultContext, bson.M{"ID": projectID})
	if projectQueryResult.Err() != nil {
		log.Println(projectQueryResult.Err().Error())
		return nil, projectQueryResult.Err()
	}

	project := models.ProjectEntry{}

	err := projectQueryResult.Decode(&project)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &project, nil
}

// GetProjectDatasets Returns all datasets of a project
func (handler *ProjectActionHandler) GetProjectDatasets(projectID string) ([]*models.DatasetEntry, error) {
	csr, err := handler.GetDatasetCollection().Find(handler.MongoDefaultContext, bson.M{
		"ProjectID": projectID,
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var projectDatasets []*models.DatasetEntry

	for csr.Next(handler.MongoDefaultContext) {
		dataset := models.DatasetEntry{}
		err := csr.Decode(&dataset)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		projectDatasets = append(projectDatasets, &dataset)
	}

	return projectDatasets, nil
}

// UserCanAccessProject Checks if a user can access a specific project
func (handler *ProjectActionHandler) UserCanAccessProject(
	requiredRight models.Right,
	userID string,
	projectID string) (bool, error) {

	queryResults := handler.GetProjectCollection().FindOne(handler.MongoDefaultContext, bson.M{
		"ID":           projectID,
		"Users.UserID": userID,
		"Users.Rights": requiredRight,
	})

	if queryResults.Err() != nil && queryResults.Err() != mongo.ErrNoDocuments {
		log.Println(queryResults.Err().Error())
		return false, queryResults.Err()
	}

	if queryResults.Err() == mongo.ErrNoDocuments {
		return false, nil
	}

	return true, nil
}

// DeleteProject Deletes a project
func (handler *ProjectActionHandler) DeleteProject(projectID string) error {
	datasets, err := handler.GetProjectDatasets(projectID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if len(datasets) != 0 {
		return fmt.Errorf("Project %v still has %v datasets associated with", projectID, len(datasets))
	}

	_, err = handler.GetProjectCollection().DeleteOne(handler.MongoDefaultContext, bson.M{"ID": projectID})

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
