package databasehandler

import (
	"testing"

	"github.com/ScienceObjectsDB/go-api/services"
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

}
