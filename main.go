package main

import (
	"log"

	"github.com/machearn/galaxy_controller/api"
	"github.com/machearn/galaxy_controller/pb"
	"github.com/machearn/galaxy_controller/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	conn, err := grpc.Dial(config.GrpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to dial gRPC server: ", err)
	}
	grpc := pb.NewGalaxyClient(conn)

	server, err := api.NewServer(config, grpc)
	if err != nil {
		log.Fatal("Failed to create server: ", err)
	}

	if err := server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
