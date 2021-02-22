package server

import (
	"context"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
)

//ProjectEndpoints Handles project related gRPC endpoints
type ProjectEndpoints struct {
	*GenericEndpoints
}

//NewProjectEndpoint A new object that handles project endpoints
func NewProjectEndpoint(genericEndpoints *GenericEndpoints) (*ProjectEndpoints, error) {
	projects := ProjectEndpoints{
		GenericEndpoints: genericEndpoints,
	}

	return &projects, nil
}

//CreateProject creates a new projects
func (endpoint *ProjectEndpoints) CreateProject(_ context.Context, _ *services.CreateProjectRequest) (*models.ProjectEntry, error) {
	panic("not implemented") // TODO: Implement
}

//AddUserToProject Adds a new user to a given project
func (endpoint *ProjectEndpoints) AddUserToProject(_ context.Context, _ *services.AddUserToProjectRequest) (*models.ProjectEntry, error) {
	panic("not implemented") // TODO: Implement
}

//GetProjectDatasets Returns all datasets that belong to a certain project
func (endpoint *ProjectEndpoints) GetProjectDatasets(_ context.Context, _ *models.ID) (*services.DatasetList, error) {
	panic("not implemented") // TODO: Implement
}

//GetUserProjects Returns all projects that a specified user has access to
func (endpoint *ProjectEndpoints) GetUserProjects(_ context.Context, _ *models.Empty) (*services.ProjectEntryList, error) {
	panic("not implemented") // TODO: Implement
}

//DeleteProject Deletes a specific project
//Will also delete all associated resources (Datasets/Objects/etc...) both from objects storage and the database
func (endpoint *ProjectEndpoints) DeleteProject(_ context.Context, _ *models.ID) (*models.Empty, error) {
	panic("not implemented") // TODO: Implement
}

func (endpoint *ProjectEndpoints) mustEmbedUnimplementedProjectAPIServer() {
	panic("not implemented") // TODO: Implement
}
