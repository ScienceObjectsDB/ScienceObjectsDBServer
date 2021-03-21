package objectstoragehandler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/spf13/viper"
)

func TestS3Handler_CreatePresignedLinks(t *testing.T) {
	key := path.Join("foo", "baa")

	os.Setenv("AWS_SECRET_ACCESS_KEY", "minioadmin")
	os.Setenv("AWS_ACCESS_KEY_ID", "minioadmin")

	viper.Set("Config.S3.Bucketname", "testbucket")

	object := models.DatasetObjectEntry{
		ID:       "test",
		Filename: "foo",
		Filetype: "txt",
		Location: &models.Location{
			Bucket:       "testbucket",
			Key:          key,
			LocationType: models.LocationType_Object,
		},
	}

	handler, err := NewS3Handler()
	if err != nil {
		t.Fatalf(err.Error())
	}

	uploadLink, err := handler.CreatePresignedUploadLink(&object)
	if err != nil {
		t.Fatalf(err.Error())
	}

	data := "data"

	req, err := http.NewRequest(http.MethodPut, uploadLink, bytes.NewBuffer([]byte(data)))
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fatalf("%v", resp)
	}

	downloadLink, err := handler.CreatePresignedDownloadLink(&object)
	if err != nil {
		t.Fatalf(err.Error())
	}

	downloadResp, err := http.Get(downloadLink)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if downloadResp.StatusCode != 200 {
		t.Fatalf("%v", resp)
	}

	respData, err := ioutil.ReadAll(downloadResp.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}

	respDataString := string(respData)

	if data != respDataString {
		t.Fatalf("Data in download response did not match original string: %v : %v", data, respDataString)
	}

}
