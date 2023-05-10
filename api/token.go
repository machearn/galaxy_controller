package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api_error "github.com/machearn/galaxy_controller/api_errors"
	"github.com/machearn/galaxy_controller/pb"
)

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenResponse struct {
	AccessToken     string    `json:"access_token"`
	AccessExpiredAt time.Time `json:"access_expired_at"`
}

func (server *Server) RenewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	grpcReq := pb.RenewAccessTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	result, err := server.grpc.RenewAccessToken(ctx, &grpcReq)
	if err != nil {
		apiErr := err.(*api_error.APIError)
		ctx.JSON(int(apiErr.Code), errorResponse(apiErr))
		return
	}

	ctx.JSON(http.StatusOK, RenewAccessTokenResponse{
		AccessToken:     result.AccessToken,
		AccessExpiredAt: result.ExpiredAt.AsTime(),
	})
}
