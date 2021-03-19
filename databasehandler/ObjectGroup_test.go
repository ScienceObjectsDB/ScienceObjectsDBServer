package databasehandler

import (
	"fmt"
	"log"
	"testing"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
	"google.golang.org/protobuf/proto"
)

func TestObjectGroupHandler_CreateDatasetObjectGroupObject(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	datasetHandler, err := NewObjectGroupHandler(dbHandler)
	if err != nil {
		t.Error(err)
	}

	datasetRequest := services.CreateObjectGroupRequest{
		Name: "foo",
		Labels: []*models.Label{
			{
				Key:   "key1",
				Value: "value1",
			},
		},
		Objects: []*services.CreateObjectRequest{
			{
				Filename:   "testfile",
				Filetype:   "txt",
				ContentLen: 9,
				Labels: []*models.Label{
					{
						Key:   "key1",
						Value: "value1",
					},
				},
			},
		},
	}

	entry, err := datasetHandler.CreateDatasetObjectGroupObject(&datasetRequest, "testproject")
	if err != nil {
		t.Error(err)
	}

	if entry.DatasetID != datasetRequest.DatasetID {
		t.Errorf("Inserted dataset id does not match")
	}

	err = datasetHandler.FinishUpload(entry.GetID())
	if err != nil {
		t.Error(err)
	}

	objectGroup, err := datasetHandler.GetObjectGroup(entry.GetID())
	if err != nil {
		t.Error(err)
	}

	entry.Status = models.Status_Available

	isEqual := proto.Equal(objectGroup, entry)
	if !isEqual {
		err := fmt.Errorf("Requrned object are not equal")
		t.Error(err)
	}

	for _, object := range objectGroup.GetObjects() {
		_, _, err := datasetHandler.GetObject(object.GetID())
		if err != nil {
			t.Error(err)
		}
	}

}
