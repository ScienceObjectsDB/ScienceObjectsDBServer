package databasehandler

import (
	"errors"
	"log"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//DatasetVersionActionHandler Handler for dataset version related database actions
type DatasetVersionActionHandler struct {
	DBUtilsHandler
}

func NewDatasetVersionHandler(dbHandler *DBUtilsHandler) (*DatasetVersionActionHandler, error) {
	handler := DatasetVersionActionHandler{
		DBUtilsHandler: *dbHandler,
	}

	return &handler, nil
}

//ReleaseDatasetVersion Releases a new dataset version
func (handler *DatasetVersionActionHandler) ReleaseDatasetVersion(request *services.ReleaseDatasetVersionRequest) (*models.DatasetVersionEntry, error) {
	csr, err := handler.GetDatasetVersionCollection().Find(handler.MongoDefaultContext, bson.M{
		"DatasetID":     request.GetDatasetID(),
		"Version.Major": request.GetVersion().GetMajor(),
		"Version.Minor": request.GetVersion().GetMinor(),
		"Version.Patch": request.GetVersion().GetPatch(),
		"Version.Stage": request.GetVersion().GetStage(),
	})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var revisionNumber int32
	revisionNumber = 1

	for csr.Next(handler.MongoDefaultContext) {
		revisionNumber++
	}

	actualVersion := request.GetVersion()
	actualVersion.Revision = int32(revisionNumber)

	uuidString := uuid.New().String()

	datasetversionEntry := models.DatasetVersionEntry{
		ID:                           uuidString,
		AdditionalMetadata:           request.GetAdditionalMetadata(),
		AdditionalMetadataMessageRef: request.AdditionalMetadataMessageRef,
		Created:                      timestamppb.Now(),
		DatasetID:                    request.GetDatasetID(),
		Status:                       models.Status_Initiating,
		ObjectCount:                  int64(len(request.GetObjectGroupIDs())),
		Version:                      actualVersion,
		ObjectIDs:                    request.ObjectGroupIDs,
	}

	result, err := handler.GetDatasetVersionCollection().InsertOne(handler.MongoDefaultContext, &datasetversionEntry)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var oid primitive.ObjectID
	var ok bool

	if oid, ok = result.InsertedID.(primitive.ObjectID); !ok {
		return nil, errors.New("Error decoding result id")
	}

	insertedDatasetVersion := models.DatasetVersionEntry{}
	insertedDatasetVersionResults := handler.GetDatasetVersionCollection().FindOne(handler.MongoDefaultContext, bson.M{
		"_id": oid,
	})

	err = insertedDatasetVersionResults.Err()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = insertedDatasetVersionResults.Decode(&insertedDatasetVersion)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &insertedDatasetVersion, nil

}

func (handler *DatasetVersionActionHandler) GetDatasetVersion(id string) (*models.DatasetVersionEntry, error) {
	result := handler.GetDatasetVersionCollection().FindOne(handler.MongoDefaultContext, bson.M{
		"ID": id,
	})

	datasetVersionEntry := models.DatasetVersionEntry{}

	if result.Err() != nil {
		log.Println(result.Err().Error())
		return nil, result.Err()
	}

	err := result.Decode(&datasetVersionEntry)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &datasetVersionEntry, nil
}

func (handler *DatasetVersionActionHandler) GetDatasetVersionDatasetID(id string) (string, error) {
	entry, err := handler.GetDatasetVersion(id)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return entry.GetDatasetID(), nil
}
