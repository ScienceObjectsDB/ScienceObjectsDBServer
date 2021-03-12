package server

import (
	"context"
	"fmt"
	"log"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
)

//DatasetEndpoints Handles dataset related gRPC endpoints
type DatasetEndpoints struct {
	*GenericEndpoints
	services.UnimplementedDatasetServiceServer
}

//NewDatasetEndpoints implements the dataset service endpoints
func NewDatasetEndpoints(genericEndpoints *GenericEndpoints) (*DatasetEndpoints, error) {

	return &DatasetEndpoints{
		GenericEndpoints: genericEndpoints,
	}, nil
}

// CreateNewDataset Creates a new dataset and associates it with a dataset
func (datasetEndpoint *DatasetEndpoints) CreateNewDataset(ctx context.Context, request *services.CreateDatasetRequest) (*models.DatasetEntry, error) {
	authorized, err := datasetEndpoint.AuthHandler.Authorize(ctx, models.Resource_Project, models.Right_Write, request.ProjectID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access denied: Can not authorize %v access to %v %v", models.Right_Write, models.Resource_Project, request.ProjectID)
		log.Println(err.Error())
		return nil, err

	}

	entry, err := datasetEndpoint.DatasetHandler.CreateNewDataset(request)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return entry, nil
}

// Dataset Returns a specific dataset
func (datasetEndpoint *DatasetEndpoints) Dataset(ctx context.Context, id *models.ID) (*models.DatasetEntry, error) {
	authorized, err := datasetEndpoint.AuthHandler.Authorize(ctx, models.Resource_Dataset, models.Right_Read, id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access denied: Can not authorize %v access to %v %v", models.Right_Write, models.Resource_Dataset, id.GetID())
		log.Println(err.Error())
		return nil, err

	}

	entry, err := datasetEndpoint.DatasetHandler.GetDataset(id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return entry, nil
}

//DatasetVersions Lists Versions of a dataset
func (datasetEndpoint *DatasetEndpoints) DatasetVersions(ctx context.Context, id *models.ID) (*services.DatasetVersionList, error) {
	authorized, err := datasetEndpoint.AuthHandler.Authorize(ctx, models.Resource_DatasetVersion, models.Right_Read, id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access denied: Can not authorize %v access to %v %v", models.Right_Read, models.Resource_Dataset, id.GetID())
		log.Println(err.Error())
		return nil, err

	}

	entries, err := datasetEndpoint.DatasetHandler.GetDatasetVersions(id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	versionList := services.DatasetVersionList{
		DatasetVersions: entries,
	}

	return &versionList, nil
}

//UpdateDatasetField Updates a field of a dataset
func (datasetEndpoint *DatasetEndpoints) UpdateDatasetField(_ context.Context, _ *models.UpdateFieldsRequest) (*models.DatasetEntry, error) {
	panic("not implemented") // TODO: Implement
}

// DeleteDataset Delete a dataset
func (datasetEndpoint *DatasetEndpoints) DeleteDataset(_ context.Context, _ *models.ID) (*models.Empty, error) {
	panic("not implemented") // TODO: Implement
}

//ReleaseDatasetVersion Release a new dataset version
func (datasetEndpoint *DatasetEndpoints) ReleaseDatasetVersion(ctx context.Context, request *services.ReleaseDatasetVersionRequest) (*models.DatasetVersionEntry, error) {
	authorized, err := datasetEndpoint.AuthHandler.Authorize(ctx, models.Resource_Dataset, models.Right_Write, request.GetDatasetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access denied: Can not authorize %v access to %v %v", models.Right_Write, models.Resource_Dataset, request.GetDatasetID())
		log.Println(err.Error())
		return nil, err

	}

	version, err := datasetEndpoint.DatasetVersionHandler.ReleaseDatasetVersion(request)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return version, nil
}
