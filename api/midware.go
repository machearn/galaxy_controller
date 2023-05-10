package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	api_error "github.com/machearn/galaxy_controller/api_errors"
	"github.com/machearn/galaxy_controller/pb"
)

type AuthPayload struct {
	ID        string    `json:"id"`
	UserID    int32     `json:"user_id"`
	CreateAt  time.Time `json:"create_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func authMiddleware(server *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if len(authHeader) == 0 {
			err := errors.New("authorization header is required")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("authorization header is invalid")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != "bearer" {
			err := errors.New("authorization header is invalid")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		if len(accessToken) == 0 {
			err := errors.New("access token is required")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		grpcReq := pb.AuthRequest{
			Token: accessToken,
		}

		result, err := server.grpc.Authorize(ctx, &grpcReq)
		if err != nil {
			apiErr := err.(*api_error.APIError)
			ctx.JSON(int(apiErr.Code), errorResponse(apiErr))
			return
		}

		ctx.Set("auth_payload", &AuthPayload{
			ID:        result.ID,
			UserID:    result.UserId,
			CreateAt:  result.CreatedAt.AsTime(),
			ExpiredAt: result.ExpiredAt.AsTime(),
		})
		ctx.Next()
	}
}
