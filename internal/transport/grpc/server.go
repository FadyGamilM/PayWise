package grpc

import (
	"log"
	"paywise/config"
	"paywise/internal/transport/grpc/pb"
)

type grpcServer struct {
	// our implementation of the grpc server will compose this method to satisfy the grpc server, this method makes our GrpcServer (sort of) implements all the methods generated in the golang code
	pb.UnimplementedPaywiseServer

	// the port that our server will recieve requests on it
	Address string
}

func NewServer() (*grpcServer, error) {
	configs, err := config.LoadGrpcServerConfig("./config")
	if err != nil {
		log.Printf("error trying to read grpc server config : %v\n", err)
		return nil, err
	}
	return &grpcServer{
		Address: configs.Grpcserver.Port,
	}, nil
}
