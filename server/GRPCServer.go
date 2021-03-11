package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/authhandler"
	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/databasehandler"
	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/objectstoragehandler"
	"github.com/ScienceObjectsDB/go-api/services"
	"google.golang.org/grpc"
)

type GenericEndpoints struct {
	AuthHandler           authhandler.AuthHandler
	ObjectStorageHandler  *objectstoragehandler.S3Handler
	ProjectActionHandler  *databasehandler.ProjectActionHandler
	DatasetHandler        *databasehandler.DatasetActionHandler
	DatasetVersionHandler *databasehandler.DatasetVersionActionHandler
	ObjectGroupHandler    *databasehandler.ObjectGroupHandler
}

type GRPCServerHandler struct {
}

// StartGRPCServer Starts the GRPC server
func (server *GRPCServerHandler) StartGRPCServer(port int64) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = server.StartGRPCServerWithListener(listener)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (server *GRPCServerHandler) StartGRPCServerWithListener(listener net.Listener) error {
	grpcServer := grpc.NewServer()

	genericEndpoints, err := server.initGenericEndpoints()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	projectEndpoints, err := NewProjectEndpoint(genericEndpoints)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	datasetEndpoints, err := NewDatasetEndpoints(genericEndpoints)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	objectEndpoints, err := NewObjectEndpoints(genericEndpoints)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	loadEndpoints, err := NewLoadEndpoints(genericEndpoints)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	services.RegisterProjectAPIServer(grpcServer, projectEndpoints)
	services.RegisterDatasetServiceServer(grpcServer, datasetEndpoints)
	services.RegisterDatasetObjectsServiceServer(grpcServer, objectEndpoints)
	services.RegisterObjectLoadServer(grpcServer, loadEndpoints)

	return nil
}

func (server *GRPCServerHandler) initGenericEndpoints() (*GenericEndpoints, error) {
	ctx := context.Background()

	client, err := databasehandler.NewMongoClient(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	dbHandler, err := databasehandler.NewDBUtilsHandler(client, ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	tokenHandler := databasehandler.TokenActionHandler{
		DBUtilsHandler: dbHandler,
	}

	projectHandler := databasehandler.ProjectActionHandler{
		DBUtilsHandler: dbHandler,
	}

	objectstorageHandler, err := objectstoragehandler.NewS3Handler()

	datasetHandler, err := databasehandler.NewDatasetHandler(dbHandler)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	datasetVersionHandler, err := databasehandler.NewDatasetVersionHandler(dbHandler)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	objectGroupHandler, err := databasehandler.NewObjectGroupHandler(dbHandler)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	auth, err := authhandler.InitProjectHandler(&projectHandler, &tokenHandler, datasetHandler, datasetVersionHandler, objectGroupHandler)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	genericEndpoints := GenericEndpoints{
		AuthHandler:           auth,
		ProjectActionHandler:  &projectHandler,
		ObjectStorageHandler:  objectstorageHandler,
		DatasetHandler:        datasetHandler,
		DatasetVersionHandler: datasetVersionHandler,
		ObjectGroupHandler:    objectGroupHandler,
	}

	return &genericEndpoints, nil
}
