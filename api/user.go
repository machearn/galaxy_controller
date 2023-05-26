package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/machearn/galaxy_controller/pb"
	"github.com/machearn/galaxy_controller/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Plan      int32     `json:"plan"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	AutoRenew bool      `json:"auto_renew"`
}

type LoginResponse struct {
	User             User      `json:"user"`
	AccessToken      string    `json:"access_token"`
	AccessExpiredAt  time.Time `json:"access_expired_at"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiredAt time.Time `json:"refresh_expired_at"`
}

func (server *Server) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	grpcGetUserRequest := pb.GetUserByUsernameRequest{
		Username: req.Username,
	}

	userResult, err := server.grpc.GetUserByUsername(ctx, &grpcGetUserRequest)
	if err != nil {
		if apiErr, ok := status.FromError(err); ok {
			if apiErr.Code() == codes.NotFound {
				ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("username or password is incorrect")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(apiErr.Err()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, userResult.GetPassword())
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("username or password is incorrect")))
		return
	}

	grpcCreateSessionReq := pb.CreateSessionRequest{
		UserId:    userResult.GetUser().GetID(),
		ClientIp:  ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
	}

	result, err := server.grpc.CreateSession(ctx, &grpcCreateSessionReq)
	if err != nil {
		if apiErr, ok := status.FromError(err); ok {
			if apiErr.Code() == codes.InvalidArgument {
				ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("username or password is incorrect")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(apiErr.Err()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user := userResult.GetUser()
	res := LoginResponse{
		User: User{
			ID:        user.ID,
			Username:  user.Username,
			Fullname:  user.Fullname,
			Email:     user.Email,
			Plan:      user.Plan,
			CreatedAt: user.CreatedAt.AsTime(),
			ExpiredAt: user.ExpiredAt.AsTime(),
			AutoRenew: user.AutoRenew,
		},
		AccessToken:      result.GetAccessToken(),
		AccessExpiredAt:  result.GetExpiredAt().AsTime(),
		RefreshToken:     result.GetSession().GetRefreshToken(),
		RefreshExpiredAt: result.GetSession().GetExpiredAt().AsTime(),
	}

	ctx.JSON(http.StatusOK, res)
}

type CreateUserRequest struct {
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Plan      int32  `json:"plan"`
	AutoRenew bool   `json:"auto_renew"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.grpc.GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{Username: req.Username})
	if err == nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("username already exists")))
		return
	} else {
		if apiErr, ok := status.FromError(err); ok {
			if apiErr.Code() != codes.NotFound {
				ctx.JSON(http.StatusInternalServerError, errorResponse(apiErr.Err()))
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	grpcReq := pb.CreateUserRequest{
		Username:  req.Username,
		Fullname:  req.Fullname,
		Email:     req.Email,
		Password:  hashedPassword,
		Plan:      req.Plan,
		AutoRenew: req.AutoRenew,
	}

	result, err := server.grpc.CreateUser(ctx, &grpcReq)
	if err != nil {
		if apiErr, ok := status.FromError(err); ok {
			if apiErr.Code() == codes.InvalidArgument {
				ctx.JSON(http.StatusBadRequest, errorResponse(apiErr.Err()))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(apiErr.Err()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user := result.GetUser()
	ctx.JSON(http.StatusOK, User{
		ID:        user.ID,
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		Plan:      user.Plan,
		CreatedAt: user.CreatedAt.AsTime(),
		ExpiredAt: user.ExpiredAt.AsTime(),
		AutoRenew: user.AutoRenew,
	})
}

type UpdateUserRequest struct {
	ID        int32   `json:"id"`
	Username  *string `json:"username"`
	Fullname  *string `json:"fullname"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	Plan      *int32  `json:"plan"`
	AutoRenew *bool   `json:"auto_renew"`
}

func (server *Server) UpdateUser(ctx *gin.Context) {
	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet("auth_payload").(*AuthPayload)
	if authPayload.UserID != req.ID {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("you are not allowed to access this resource")))
		return
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(*req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		req.Password = &hashedPassword
	}

	grpcReq := pb.UpdateUserRequest{
		ID:        req.ID,
		Username:  req.Username,
		Fullname:  req.Fullname,
		Email:     req.Email,
		Password:  req.Password,
		Plan:      req.Plan,
		AutoRenew: req.AutoRenew,
	}

	result, err := server.grpc.UpdateUser(ctx, &grpcReq)
	if err != nil {
		if apiErr, ok := status.FromError(err); ok {
			if apiErr.Code() == codes.InvalidArgument {
				ctx.JSON(http.StatusBadRequest, errorResponse(apiErr.Err()))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(apiErr.Err()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user := result.GetUser()
	ctx.JSON(http.StatusOK, User{
		ID:        user.ID,
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		Plan:      user.Plan,
		CreatedAt: user.CreatedAt.AsTime(),
		ExpiredAt: user.ExpiredAt.AsTime(),
		AutoRenew: user.AutoRenew,
	})
}

type GetUserRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) GetUser(ctx *gin.Context) {
	var req GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet("auth_payload").(*AuthPayload)
	if authPayload.UserID != req.ID {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("you are not allowed to access this resource")))
		return
	}

	grpcReq := pb.GetUserRequest{
		ID: req.ID,
	}

	result, err := server.grpc.GetUser(ctx, &grpcReq)
	if err != nil {
		if apiErr, ok := status.FromError(err); ok {
			if apiErr.Code() == codes.NotFound {
				ctx.JSON(http.StatusNotFound, errorResponse(apiErr.Err()))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(apiErr.Err()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user := result.GetUser()
	ctx.JSON(http.StatusOK, User{
		ID:        user.ID,
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		Plan:      user.Plan,
		CreatedAt: user.CreatedAt.AsTime(),
		ExpiredAt: user.ExpiredAt.AsTime(),
		AutoRenew: user.AutoRenew,
	})
}
