package databasehandler

import (
	"fmt"
	"log"
	"path"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
		Name:               request.Name,
		InitializedObjects: int64(len(request.Objects)),
		Labels:             request.Labels,
		AdditionalMetadata: request.AdditionalMetadata,
		DatasetID:          request.DatasetID,
		UploadedObjects:    0,
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
			Status:             models.Status_Available,
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
