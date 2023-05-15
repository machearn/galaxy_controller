package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/machearn/galaxy_controller/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if apiErr, ok := status.FromError(err); ok {
			if apiErr.Code() == codes.Unauthenticated {
				ctx.JSON(http.StatusUnauthorized, errorResponse(apiErr.Err()))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(apiErr.Err()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, RenewAccessTokenResponse{
		AccessToken:     result.AccessToken,
		AccessExpiredAt: result.ExpiredAt.AsTime(),
	})
}
