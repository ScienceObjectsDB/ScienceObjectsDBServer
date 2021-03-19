package server

import (
	"context"
	"fmt"
	"log"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
)

//ProjectEndpoints Handles project related gRPC endpoints
type ProjectEndpoints struct {
	*GenericEndpoints
	services.UnimplementedProjectAPIServer
}

//NewProjectEndpoint A new object that handles project endpoints
func NewProjectEndpoint(genericEndpoints *GenericEndpoints) (*ProjectEndpoints, error) {
	projects := ProjectEndpoints{
		GenericEndpoints: genericEndpoints,
	}

	return &projects, nil
}

//CreateProject creates a new projects
func (endpoint *ProjectEndpoints) CreateProject(ctx context.Context, request *services.CreateProjectRequest) (*models.ProjectEntry, error) {
	userID, err := endpoint.AuthHandler.UserID(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return endpoint.ProjectActionHandler.CreateProject(userID, request)
}

//AddUserToProject Adds a new user to a given project
func (endpoint *ProjectEndpoints) AddUserToProject(_ context.Context, _ *services.AddUserToProjectRequest) (*models.ProjectEntry, error) {
	panic("not implemented") // TODO: Implement
}

//GetProjectDatasets Returns all datasets that belong to a certain project
func (endpoint *ProjectEndpoints) GetProjectDatasets(ctx context.Context, id *models.ID) (*services.DatasetList, error) {
	authorized, err := endpoint.AuthHandler.Authorize(ctx, models.Resource_Project, models.Right_Read, id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access on dataset with ID %v denied", id.GetID())
		log.Println(err.Error())
		return nil, err
	}

	datasets, err := endpoint.ProjectActionHandler.GetProjectDatasets(id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	datasetList := services.DatasetList{
		Datasets: datasets,
	}

	return &datasetList, nil
}

//GetUserProjects Returns all projects that a specified user has access to
func (endpoint *ProjectEndpoints) GetUserProjects(ctx context.Context, _ *models.Empty) (*services.ProjectEntryList, error) {
	userID, err := endpoint.AuthHandler.UserID(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	projects, err := endpoint.ProjectActionHandler.GetUserProjects(userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	projectList := services.ProjectEntryList{
		Projects: projects,
	}

	return &projectList, nil
}

//DeleteProject Deletes a specific project
//Will also delete all associated resources (Datasets/Objects/etc...) both from objects storage and the database
func (endpoint *ProjectEndpoints) DeleteProject(_ context.Context, _ *models.ID) (*models.Empty, error) {
	panic("not implemented") // TODO: Implement
}
