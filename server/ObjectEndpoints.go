package server

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
)

//ObjectEndpoints Handles object related gRPC endpoints
type ObjectEndpoints struct {
	*GenericEndpoints
	services.UnimplementedDatasetObjectsServiceServer
}

//NewObjectEndpoints New endpoints for object and object group handling
func NewObjectEndpoints(genericEndpoints *GenericEndpoints) (*ObjectEndpoints, error) {
	return &ObjectEndpoints{
		GenericEndpoints: genericEndpoints,
	}, nil
}

//CreateObjectHeritage Creates a new object heritage
func (endpoints *ObjectEndpoints) CreateObjectHeritage(ctx context.Context, request *services.CreateObjectHeritageRequest) (*models.ObjectHeritage, error) {
	panic("not implemented") // TODO: Implement
}

//CreateObjectGroup Creates a new object group
func (endpoints *ObjectEndpoints) CreateObjectGroup(ctx context.Context, request *services.CreateObjectGroupRequest) (*models.DatasetObjectGroup, error) {
	authorized, err := endpoints.AuthHandler.Authorize(ctx, models.Resource_Dataset, models.Right_Write, request.GetDatasetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access denied: Can not authorize %v access to %v %v", models.Right_Write, models.Resource_Project, request.GetDatasetID())
		log.Println(err.Error())
		return nil, err

	}

	projectID, err := endpoints.GenericEndpoints.DatasetHandler.GetDatasetProjectID(request.GetDatasetID())

	entry, err := endpoints.GenericEndpoints.ObjectGroupHandler.CreateDatasetObjectGroupObject(request, projectID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return entry, nil
}

//FinishObjectUpload Finishes the upload process for a data
func (endpoints *ObjectEndpoints) FinishObjectUpload(ctx context.Context, id *models.ID) (*models.Empty, error) {
	authorized, err := endpoints.AuthHandler.Authorize(ctx, models.Resource_DatasetObjectGroupResource, models.Right_Write, id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access denied: Can not authorize %v access to %v %v", models.Right_Write, models.Resource_Project, id.GetID())
		log.Println(err.Error())
		return nil, err

	}

	err = endpoints.GenericEndpoints.ObjectGroupHandler.FinishUpload(id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &models.Empty{}, nil
}

//GetObjectGroup Returns an object based on the given ID
func (endpoints *ObjectEndpoints) GetObjectGroup(ctx context.Context, id *models.ID) (*models.DatasetObjectGroup, error) {
	authorized, err := endpoints.AuthHandler.Authorize(ctx, models.Resource_DatasetObjectGroupResource, models.Right_Read, id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access denied: Can not authorize %v access to %v %v", models.Right_Read, models.Resource_DatasetObjectGroupResource, id.GetID())
		log.Println(err.Error())
		return nil, err

	}

	objectGroup, err := endpoints.GenericEndpoints.ObjectGroupHandler.GetObjectGroup(id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return objectGroup, nil
}

func (endpoints *ObjectEndpoints) mustEmbedUnimplementedDatasetObjectsServiceServer() {
	panic("not implemented") // TODO: Implement
}
