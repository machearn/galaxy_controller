package api

import (
	"github.com/gin-gonic/gin"
	"github.com/machearn/galaxy_controller/pb"
	"github.com/machearn/galaxy_controller/util"
)

type Server struct {
	config util.Config
	router *gin.Engine
	grpc   pb.GalaxyClient
}

func NewServer(config util.Config, grpc pb.GalaxyClient) (*Server, error) {
	server := Server{
		config: config,
		grpc:   grpc,
	}

	server.SetupRouter()

	return &server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()

	router.POST("/item/create", server.CreateItem)
	router.GET("/item/get/:id", server.GetItem)
	router.POST("/item/list", server.ListItems)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) *gin.H {
	return &gin.H{
		"error": err.Error(),
	}
}
