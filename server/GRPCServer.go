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
	AuthHandler          authhandler.AuthHandler
	ObjectStorageHandler *objectstoragehandler.S3Handler
	ProjectActionHandler *databasehandler.ProjectActionHandler
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

	services.RegisterProjectAPIServer(grpcServer, projectEndpoints)
	services.RegisterDatasetServiceServer(grpcServer, datasetEndpoints)

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

	auth, err := authhandler.InitProjectHandler(&projectHandler, &tokenHandler)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	genericEndpoints := GenericEndpoints{
		AuthHandler:          auth,
		ProjectActionHandler: &projectHandler,
		ObjectStorageHandler: objectstorageHandler,
	}

	return &genericEndpoints, nil
}
