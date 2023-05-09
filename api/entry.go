package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api_error "github.com/machearn/galaxy_controller/api_errors"
	"github.com/machearn/galaxy_controller/pb"
)

type CreateEntryRequest struct {
	UserID   int32 `json:"member_id"`
	ItemID   int32 `json:"item_id"`
	Quantity int32 `json:"quantity"`
	Total    int32 `json:"total"`
}

type Entry struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"member_id"`
	ItemID    int32     `json:"item_id"`
	Quantity  int32     `json:"quantity"`
	Total     int32     `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) CreateEntry(ctx *gin.Context) {
	var req CreateEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.grpc.GetUser(ctx, &pb.GetUserRequest{ID: req.UserID})
	if err != nil {
		apiErr := err.(*api_error.APIError)
		if apiErr.Code == http.StatusNotFound {
			ctx.JSON(http.StatusBadRequest, errorResponse(apiErr))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(apiErr))
		return
	}

	_, err = server.grpc.GetItem(ctx, &pb.GetItemRequest{Id: req.ItemID})
	if err != nil {
		apiErr := err.(*api_error.APIError)
		if apiErr.Code == http.StatusNotFound {
			ctx.JSON(http.StatusBadRequest, errorResponse(apiErr))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(apiErr))
		return
	}

	grpcReq := pb.CreateEntryRequest{
		UserId:   req.UserID,
		ItemId:   req.ItemID,
		Quantity: req.Quantity,
		Total:    req.Total,
	}

	result, err := server.grpc.CreateEntry(ctx, &grpcReq)
	if err != nil {
		apiErr := err.(*api_error.APIError)
		ctx.JSON(int(apiErr.Code), errorResponse(apiErr))
		return
	}

	entry := result.GetEntry()
	ctx.JSON(http.StatusOK, Entry{
		ID:        entry.ID,
		UserID:    entry.UserId,
		ItemID:    entry.ItemId,
		Quantity:  entry.Quantity,
		Total:     entry.Total,
		CreatedAt: entry.CreatedAt.AsTime(),
	})
}

type GetEntryRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetEntry(ctx *gin.Context) {
	var req GetEntryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.grpc.GetEntry(ctx, &pb.GetEntryRequest{Id: req.ID})
	if err != nil {
		apiErr := err.(*api_error.APIError)
		ctx.JSON(int(apiErr.Code), errorResponse(apiErr))
		return
	}

	entry := result.GetEntry()
	ctx.JSON(http.StatusOK, Entry{
		ID:        entry.ID,
		UserID:    entry.UserId,
		ItemID:    entry.ItemId,
		Quantity:  entry.Quantity,
		Total:     entry.Total,
		CreatedAt: entry.CreatedAt.AsTime(),
	})
}

type ListEntriesRequest struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

type ListEntriesResponse struct {
	Entries []Entry `json:"entries"`
}

func (server *Server) ListEntries(ctx *gin.Context) {
	var req ListEntriesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	grpcReq := pb.ListEntriesRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	result, err := server.grpc.ListEntries(ctx, &grpcReq)
	if err != nil {
		apiErr := err.(*api_error.APIError)
		ctx.JSON(int(apiErr.Code), errorResponse(apiErr))
		return
	}

	rows := result.GetEntries()
	entries := make([]Entry, len(rows))
	for i, row := range rows {
		entries[i] = Entry{
			ID:        row.ID,
			UserID:    row.UserId,
			ItemID:    row.ItemId,
			Quantity:  row.Quantity,
			Total:     row.Total,
			CreatedAt: row.CreatedAt.AsTime(),
		}
	}

	ctx.JSON(http.StatusOK, ListEntriesResponse{
		Entries: entries,
	})
}

type ListEntriesByUserRequest struct {
	UserID int32 `json:"user_id"`
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (server *Server) ListEntriesByUser(ctx *gin.Context) {
	var req ListEntriesByUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	grpcReq := pb.ListEntriesByUserRequest{
		UserId: req.UserID,
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	result, err := server.grpc.ListEntriesByUser(ctx, &grpcReq)
	if err != nil {
		apiErr := err.(*api_error.APIError)
		ctx.JSON(int(apiErr.Code), errorResponse(apiErr))
		return
	}

	rows := result.GetEntries()
	entries := make([]Entry, len(rows))
	for i, row := range rows {
		entries[i] = Entry{
			ID:        row.ID,
			UserID:    row.UserId,
			ItemID:    row.ItemId,
			Quantity:  row.Quantity,
			Total:     row.Total,
			CreatedAt: row.CreatedAt.AsTime(),
		}
	}

	ctx.JSON(http.StatusOK, ListEntriesResponse{
		Entries: entries,
	})
}

type ListEntriesByItemRequest struct {
	ItemID int32 `json:"item_id"`
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (server *Server) ListEntriesByItem(ctx *gin.Context) {
	var req ListEntriesByItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	grpcReq := pb.ListEntriesByItemRequest{
		ItemId: req.ItemID,
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	result, err := server.grpc.ListEntriesByItem(ctx, &grpcReq)
	if err != nil {
		apiErr := err.(*api_error.APIError)
		ctx.JSON(int(apiErr.Code), errorResponse(apiErr))
		return
	}

	rows := result.GetEntries()
	entries := make([]Entry, len(rows))
	for i, row := range rows {
		entries[i] = Entry{
			ID:        row.ID,
			UserID:    row.UserId,
			ItemID:    row.ItemId,
			Quantity:  row.Quantity,
			Total:     row.Total,
			CreatedAt: row.CreatedAt.AsTime(),
		}
	}

	ctx.JSON(http.StatusOK, ListEntriesResponse{
		Entries: entries,
	})
}
