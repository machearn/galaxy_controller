package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/machearn/galaxy_controller/pb"
)

type Item struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Quantity int32  `json:"quantity"`
	Price    int32  `json:"price"`
}

type CreateItemRequest struct {
	Name     string `json:"name"`
	Quantity int32  `json:"quantity"`
	Price    int32  `json:"price"`
}

func (server *Server) CreateItem(ctx *gin.Context) {
	var req CreateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	grpcReq := pb.CreateItemRequest{
		Name:     req.Name,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	result, err := server.grpc.CreateItem(ctx, &grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	item := result.GetItem()
	ctx.JSON(http.StatusOK, Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		Price:    item.Price,
	})
}

type GetItemRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetItem(ctx *gin.Context) {
	var req GetItemRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	grpcReq := pb.GetItemRequest{
		Id: req.ID,
	}

	result, err := server.grpc.GetItem(ctx, &grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	item := result.GetItem()
	ctx.JSON(http.StatusOK, Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		Price:    item.Price,
	})
}

type ListItemsRequest struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

type ListItemsResponse struct {
	Items []Item `json:"items"`
}

func (server *Server) ListItems(ctx *gin.Context) {
	var req ListItemsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	grpcReq := pb.ListItemsRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	result, err := server.grpc.ListItems(ctx, &grpcReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rows := result.GetItems()
	var res ListItemsResponse
	for _, row := range rows {
		res.Items = append(res.Items, Item{
			ID:       row.ID,
			Name:     row.Name,
			Quantity: row.Quantity,
			Price:    row.Price,
		})
	}

	ctx.JSON(http.StatusOK, res)
}
