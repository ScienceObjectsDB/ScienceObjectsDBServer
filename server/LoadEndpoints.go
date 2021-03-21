package server

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/ScienceObjectsDB/go-api/services"
)

//LoadEndpoints for the ObjectLoad service of the API
type LoadEndpoints struct {
	*GenericEndpoints
	services.UnimplementedObjectLoadServer
}

//NewLoadEndpoints returns a new object to handle the ObjectLoad service
//
func NewLoadEndpoints(genericEndpoints *GenericEndpoints) (*LoadEndpoints, error) {
	return &LoadEndpoints{
		GenericEndpoints: genericEndpoints,
	}, nil
}

//CreateUploadLink Returns an upload link for an individual object
func (endpoint *LoadEndpoints) CreateUploadLink(ctx context.Context, id *models.ID) (*services.CreateUploadLinkResponse, error) {
	authorized, err := endpoint.AuthHandler.Authorize(ctx, models.Resource_DatasetObject, models.Right_Write, id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access denied: Can not authorize %v access to %v %v", models.Right_Write, models.Resource_DatasetObject, id.GetID())
		log.Println(err.Error())
		return nil, err
	}

	_, object, err := endpoint.GenericEndpoints.ObjectGroupHandler.GetObject(id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	link, err := endpoint.GenericEndpoints.ObjectStorageHandler.CreatePresignedUploadLink(object)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	uploadLinkResponse := services.CreateUploadLinkResponse{
		UploadLink: link,
		Object:     object,
	}

	return &uploadLinkResponse, nil
}

//CreateDownloadLink Returns an download link for an individual object
func (endpoint *LoadEndpoints) CreateDownloadLink(ctx context.Context, id *models.ID) (*services.CreateUploadLinkResponse, error) {
	authorized, err := endpoint.AuthHandler.Authorize(ctx, models.Resource_DatasetObject, models.Right_Read, id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !authorized {
		err := fmt.Errorf("Access denied: Can not authorize %v access to %v %v", models.Right_Read, models.Resource_DatasetObject, id.GetID())
		log.Println(err.Error())
		return nil, err
	}

	_, object, err := endpoint.GenericEndpoints.ObjectGroupHandler.GetObject(id.GetID())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	link, err := endpoint.GenericEndpoints.ObjectStorageHandler.CreatePresignedDownloadLink(object)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	uploadLinkResponse := services.CreateUploadLinkResponse{
		UploadLink: link,
		Object:     object,
	}

	return &uploadLinkResponse, nil
}

func (endpoint *LoadEndpoints) mustEmbedUnimplementedObjectLoadServer() {
	panic("not implemented") // TODO: Implement
}
