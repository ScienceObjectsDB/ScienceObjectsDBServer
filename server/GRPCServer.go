package server

import (
	"fmt"
	"log"
	"net"

	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/authhandler"
	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/databasehandler"
	"github.com/ScienceObjectsDB/ScienceObjectsDBServer/objectstoragehandler"
)

type GenericEndpoints struct {
	AuthHandler          *authhandler.AuthHandler
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

	return nil
}
