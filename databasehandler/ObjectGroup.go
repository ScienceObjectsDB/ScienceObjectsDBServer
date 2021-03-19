package databasehandler

import (
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SingleObject struct {
	ID      string
	Objects []*models.DatasetObjectEntry
}

//ObjectGroupHandler Handles dataset object group actions
type ObjectGroupHandler struct {
	*DBUtilsHandler
}

func NewObjectGroupHandler(dbUtilsHandler *DBUtilsHandler) (*ObjectGroupHandler, error) {
	handler := ObjectGroupHandler{
		DBUtilsHandler: dbUtilsHandler,
	}

	return &handler, nil
}

//CreateDatasetObjectGroupObject Creates a new dataset object group
func (handler *ObjectGroupHandler) CreateDatasetObjectGroupObject(request *services.CreateObjectGroupRequest, projectID string) (*models.DatasetObjectGroup, error) {
	uuidString := uuid.New().String()

	objectGroup := models.DatasetObjectGroup{
		ID:                 uuidString,
		Name:               request.Name,
		InitializedObjects: int64(len(request.Objects)),
		Labels:             request.Labels,
		AdditionalMetadata: request.AdditionalMetadata,
		DatasetID:          request.DatasetID,
		UploadedObjects:    0,
		Status:             models.Status_Available,
	}

	var objects []*models.DatasetObjectEntry

	for i, requestedObject := range request.GetObjects() {
		objectuuidString := fmt.Sprintf("%v-%v", uuidString, i)
		uploadID := uuid.New().String()

		objectKey := path.Join(projectID, request.DatasetID, uuidString, fmt.Sprintf("%v", i), requestedObject.Filename)

		object := models.DatasetObjectEntry{
			ID:                 objectuuidString,
			Filename:           requestedObject.Filename,
			Filetype:           requestedObject.Filetype,
			Created:            timestamppb.Now(),
			ContentLen:         requestedObject.ContentLen,
			AdditionalMetadata: requestedObject.AdditionalMetadata,
			UploadID:           uploadID,
			Origin:             requestedObject.Origin,
			Location: &models.Location{
				Bucket:       "fgoo",
				Key:          objectKey,
				LocationType: models.LocationType_Object,
			},
		}

		objects = append(objects, &object)

	}

	objectGroup.Objects = objects

	insertedValue := &models.DatasetObjectGroup{}

	err := handler.DBUtilsHandler.Insert(handler.DBUtilsHandler.GetDatasetObjectGroupCollection(), &objectGroup, insertedValue)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return insertedValue, nil
}

func (handler *ObjectGroupHandler) FinishUpload(objectGroupID string) error {
	_, err := handler.DBUtilsHandler.GetDatasetObjectGroupCollection().UpdateOne(handler.MongoDefaultContext,
		bson.M{"ID": objectGroupID},
		bson.M{"$set": bson.M{
			"Status": models.Status_Available,
		}},
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (handler *ObjectGroupHandler) GetObjectGroup(objectGroupID string) (*models.DatasetObjectGroup, error) {
	result := handler.DBUtilsHandler.GetDatasetObjectGroupCollection().FindOne(handler.MongoDefaultContext, bson.M{
		"ID": objectGroupID,
	})

	log.Println(objectGroupID)

	datasetObjectGroup := models.DatasetObjectGroup{}

	if result.Err() != nil {
		log.Println(result.Err().Error())
		return nil, result.Err()
	}

	err := result.Decode(&datasetObjectGroup)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &datasetObjectGroup, nil
}

func (handler *ObjectGroupHandler) GetObjectGroups(objectGroupID []string) ([]*models.DatasetObjectGroup, error) {
	var objectGroups []*models.DatasetObjectGroup

	results, err := handler.DBUtilsHandler.GetDatasetObjectGroupCollection().Find(handler.MongoDefaultContext, bson.M{
		"ID": bson.M{"$in": objectGroupID},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = results.All(handler.MongoDefaultContext, objectGroups)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return objectGroups, nil
}

func (handler *ObjectGroupHandler) GetObject(objectID string) (string, *models.DatasetObjectEntry, error) {
	option := options.FindOneOptions{
		Projection: bson.M{"Objects.$": 1, "ID": 1},
	}

	result := handler.DBUtilsHandler.GetDatasetObjectGroupCollection().FindOne(
		handler.MongoDefaultContext,
		bson.M{"Objects.ID": objectID},
		&option,
	)

	datasetObject := SingleObject{}

	if result.Err() != nil {
		log.Println(result.Err().Error())
		return "", nil, result.Err()
	}

	err := result.Decode(&datasetObject)
	if err != nil {
		log.Println(err.Error())
		return "", nil, err
	}

	return datasetObject.ID, datasetObject.Objects[0], nil
}

//GetDatasetObjects Lists all objectgroups of a dataset
func (handler *ObjectGroupHandler) GetDatasetObjects(datasetID string) ([]*models.DatasetObjectGroup, error) {
	var objectGroups []*models.DatasetObjectGroup

	results, err := handler.DBUtilsHandler.GetDatasetObjectGroupCollection().Find(handler.MongoDefaultContext, bson.M{
		"DatasetID": datasetID,
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = results.All(handler.MongoDefaultContext, &objectGroups)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return objectGroups, nil
}
