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

	router.POST("/user/login", server.Login)
	router.POST("/user/create", server.CreateUser)
	router.POST("/token/renew", server.RenewAccessToken)

	authRouter := router.Group("/").Use(authMiddleware(server))

	authRouter.GET("/user/get/:id", server.GetUser)
	authRouter.POST("/user/update", server.UpdateUser)
	authRouter.POST("/item/create", server.CreateItem)
	authRouter.GET("/item/get/:id", server.GetItem)
	authRouter.POST("/item/list", server.ListItems)
	authRouter.POST("/item/update", server.UpdateItem)
	authRouter.DELETE("/item/delete/:id", server.DeleteItem)
	authRouter.POST("/entry/create", server.CreateEntry)
	authRouter.GET("/entry/get/:id", server.GetEntry)
	authRouter.POST("/entry/list", server.ListEntries)
	authRouter.POST("/entry/list/user", server.ListEntriesByUser)
	authRouter.POST("/entry/list/item", server.ListEntriesByItem)

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
