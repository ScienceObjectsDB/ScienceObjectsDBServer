package databasehandler

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/util"
	"github.com/ScienceObjectsDB/go-api/services"
)

var dbHandler *DBUtilsHandler

func TestMain(m *testing.M) {
	err := util.InitTestEnv()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = setup()
	if err != nil {
		log.Fatalln(err.Error())
	}

	code := m.Run()
	os.Exit(code)
}

func setup() error {
	ctx := context.Background()

	client, err := NewMongoClient(ctx)
	if err != nil {
		return err
	}

	dbHandler, err = NewDBUtilsHandler(client, ctx)
	if err != nil {
		return err
	}

	return nil
}

func Test_Dataset(t *testing.T) {
	datasetHandler, err := NewDatasetHandler(dbHandler)
	if err != nil {
		t.Error(err)
	}

	datasetRequest := services.CreateDatasetRequest{
		DatasetName: "test123",
		Datatype:    "txt",
		ProjectID:   "145",
	}

	entry, err := datasetHandler.CreateNewDataset(&datasetRequest)
	if err != nil {
		t.Error(err)
	}

	if entry.Datasetname != datasetRequest.DatasetName {
		t.Errorf("Inserted datasetname does not match")
	}

	if entry.Datasettype != datasetRequest.Datatype {
		t.Errorf("Inserted datasettype does not match")
	}

	if entry.ProjectID != datasetRequest.ProjectID {
		t.Errorf("Inserted projectID does not match")
	}

	panic("foo")

}
