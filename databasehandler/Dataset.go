package databasehandler

import (
	"context"
	"errors"
	"log"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//DatasetActionHandler Handler for dataset related database actions
type DatasetActionHandler struct {
	*DBUtilsHandler
}

//NewDatasetHandler Initializes a new dataset action handler
func NewDatasetHandler(dbUtilsHandler *DBUtilsHandler) (*DatasetActionHandler, error) {
	handler := DatasetActionHandler{
		DBUtilsHandler: dbUtilsHandler,
	}

	return &handler, nil
}

// CreateNewDataset Creates and inserts a new dataset
func (handler *DatasetActionHandler) CreateNewDataset(request *services.CreateDatasetRequest) (*models.DatasetEntry, error) {
	uuidString := uuid.New().String()

	datasetEntry := models.DatasetEntry{
		ID:          uuidString,
		ProjectID:   request.ProjectID,
		Status:      models.Status_Available,
		IsPublic:    false,
		Datasetname: request.DatasetName,
		Datasettype: request.Datatype,
		Created:     timestamppb.Now(),
	}

	insertedResult, err := handler.GetDatasetCollection().InsertOne(handler.MongoDefaultContext, &datasetEntry)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var oid primitive.ObjectID
	var ok bool

	if oid, ok = insertedResult.InsertedID.(primitive.ObjectID); !ok {
		return nil, errors.New("Error decoding result id")
	}

	result := handler.GetDatasetCollection().FindOne(handler.MongoDefaultContext, bson.M{
		"_id": oid,
	})

	if result.Err() != nil {
		log.Println(result.Err().Error())
		return nil, err
	}

	err = result.Decode(&datasetEntry)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &datasetEntry, nil
}

// GetDataset Returns the dataset with the provided ID
func (handler *DatasetActionHandler) GetDataset(datasetID string) (*models.DatasetEntry, error) {
	queryResult := handler.GetDatasetCollection().FindOne(handler.MongoDefaultContext, bson.M{
		"ID": datasetID,
	})

	if queryResult.Err() != nil && queryResult.Err() != mongo.ErrNoDocuments {
		log.Println(queryResult.Err().Error())
		return nil, queryResult.Err()
	}

	if queryResult.Err() == mongo.ErrNoDocuments {
		return &models.DatasetEntry{}, nil
	}

	var datasetEntry models.DatasetEntry

	err := queryResult.Decode(&datasetEntry)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &datasetEntry, nil
}

// DeleteDataset Deletes a given dataset
// Datasets should only be deleted if no data objects and dataset versions are associated with this project
func (handler *DatasetActionHandler) DeleteDataset(datasetid string) error {
	_, err := handler.GetDatasetCollection().DeleteOne(handler.MongoDefaultContext, bson.M{
		"ID": datasetid,
	})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

//GetDatasetProjectID Returns the project id of the dataset with the provided id
func (handler *DatasetActionHandler) GetDatasetProjectID(datasetid string) (string, error) {
	entry, err := handler.GetDataset(datasetid)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return entry.GetProjectID(), nil
}

func (handler *DatasetActionHandler) GetDatasetVersions(datasetid string) ([]*models.DatasetVersionEntry, error) {
	var entries []*models.DatasetVersionEntry

	queryResult, err := handler.GetDatasetVersionCollection().Find(handler.MongoDefaultContext, bson.M{
		"DatasetID": datasetid,
	})

	if err != nil {
		log.Println(err.Error())
		return entries, err
	}

	err = queryResult.All(context.Background(), &entries)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return entries, nil
}
