package databasehandler

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/util"
	"github.com/ScienceObjectsDB/go-api/models"
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

}

func TestDatasetVersions(t *testing.T) {
	datasetHandler, err := NewDatasetHandler(dbHandler)
	if err != nil {
		t.Error(err)
	}

	datasetVersionHandler, err := NewDatasetVersionHandler(dbHandler)
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

	datasetVersionRequest1 := services.ReleaseDatasetVersionRequest{
		Name:      "foo",
		DatasetID: entry.GetID(),
		Version: &models.Version{
			Major:    0,
			Minor:    0,
			Patch:    1,
			Revision: 1,
			Stage:    models.Version_Stable,
		},
		ObjectGroupIDs: make([]string, 0),
	}

	_, err = datasetVersionHandler.ReleaseDatasetVersion(&datasetVersionRequest1)
	if err != nil {
		t.Error(err)
	}

	datasetVersionRequest2 := services.ReleaseDatasetVersionRequest{
		Name:      "bar",
		DatasetID: entry.GetID(),
		Version: &models.Version{
			Major:    0,
			Minor:    0,
			Patch:    1,
			Revision: 1,
			Stage:    models.Version_Stable,
		},
		ObjectGroupIDs: make([]string, 0),
	}

	versionEntry2, err := datasetVersionHandler.ReleaseDatasetVersion(&datasetVersionRequest2)
	if err != nil {
		t.Error(err)
	}

	if versionEntry2.GetVersion().GetRevision() != 2 {
		t.Errorf("Wrong revision number")
	}

	datasetVersionRequest3 := services.ReleaseDatasetVersionRequest{
		Name:      "bar",
		DatasetID: entry.GetID(),
		Version: &models.Version{
			Major:    0,
			Minor:    0,
			Patch:    2,
			Revision: 1,
			Stage:    models.Version_Stable,
		},
		ObjectGroupIDs: make([]string, 0),
	}

	_, err = datasetVersionHandler.ReleaseDatasetVersion(&datasetVersionRequest3)
	if err != nil {
		t.Error(err)
	}

	versionEntries, err := datasetHandler.GetDatasetVersions(entry.GetID())
	if err != nil {
		t.Error(err)
	}

	if len(versionEntries) != 3 {
		t.Errorf("Wrong number of version entries found, expected: %v, found: %v", 3, len(versionEntries))
	}

}
