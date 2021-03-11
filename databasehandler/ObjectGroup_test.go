package databasehandler

import (
	"fmt"
	"testing"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
	"google.golang.org/protobuf/proto"
)

func TestObjectGroupHandler_CreateDatasetObjectGroupObject(t *testing.T) {
	datasetHandler, err := NewObjectGroupHandler(dbHandler)
	if err != nil {
		t.Error(err)
	}

	datasetRequest := services.CreateObjectGroupRequest{
		Name: "foo",
		Labels: map[string]string{
			"test1": "foo",
			"test2": "baa",
		},
		Objects: []*services.CreateObjectRequest{
			{
				Filename:   "testfile",
				Filetype:   "txt",
				ContentLen: 9,
				Labels: map[string]string{
					"test3": "bar",
					"test4": "baz",
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

	for key, _ := range objectGroup.GetObjects() {
		_, _, err := datasetHandler.GetObject(key)
		if err != nil {
			t.Error(err)
		}
	}

}
