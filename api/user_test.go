package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	api_error "github.com/machearn/galaxy_controller/api_errors"
	"github.com/machearn/galaxy_controller/pb"
	mockpb "github.com/machearn/galaxy_controller/pb/mock"
	"github.com/machearn/galaxy_controller/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestLoginAPI(t *testing.T) {
	url := "/user/login"
	data, err := json.Marshal(gin.H{
		"username": "test",
		"password": "test",
	})
	require.NoError(t, err)

	created := time.Now().UTC().Truncate(time.Second)
	expired := created.Add(time.Hour * 24 * 30)

	refreshToken := util.GetRandomString(32)
	accessToken := util.GetRandomString(32)

	hashedPassword, err := util.HashPassword("test")
	require.NoError(t, err)

	grpcGetUserReq := pb.GetUserByUsernameRequest{
		Username: "test",
	}
	grpcGetUserRes := pb.GetUserResponse{
		User: &pb.User{
			ID:        1,
			Username:  "test",
			Fullname:  "test",
			Email:     "test",
			Plan:      1,
			CreatedAt: timestamppb.New(created),
			ExpiredAt: timestamppb.New(expired),
			AutoRenew: true,
		},
		Password: hashedPassword,
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	grpcCreateSessionReq := pb.CreateSessionRequest{
		UserId:    1,
		ClientIp:  request.RemoteAddr,
		UserAgent: request.UserAgent(),
	}
	grpcCreateSessionRes := pb.CreateSessionResponse{
		AccessToken: accessToken,
		ExpiredAt:   timestamppb.New(expired),
		Session: &pb.Session{
			ID:           uuid.New().String(),
			UserId:       1,
			ClientIp:     request.RemoteAddr,
			UserAgent:    request.UserAgent(),
			RefreshToken: refreshToken,
			CreatedAt:    timestamppb.New(created),
			ExpiredAt:    timestamppb.New(expired),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(&grpcGetUserReq)).Return(&grpcGetUserRes, nil)
	grpc.EXPECT().CreateSession(gomock.Any(), gomock.Eq(&grpcCreateSessionReq)).Return(&grpcCreateSessionRes, nil)

	server := NewTestServer(t, grpc)
	recorder := httptest.NewRecorder()

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	var res LoginResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, accessToken, res.AccessToken)
	require.Equal(t, refreshToken, res.RefreshToken)
	require.Equal(t, expired, res.AccessExpiredAt)
	require.Equal(t, expired, res.RefreshExpiredAt)
	require.Equal(t, int32(1), res.User.ID)
	require.Equal(t, "test", res.User.Username)
	require.Equal(t, "test", res.User.Fullname)
	require.Equal(t, "test", res.User.Email)
	require.Equal(t, int32(1), res.User.Plan)
	require.Equal(t, created, res.User.CreatedAt)
	require.Equal(t, expired, res.User.ExpiredAt)
}

func TestGetUserAPI(t *testing.T) {
	url := "/user/get/1"

	created := time.Now().UTC().Truncate(time.Second)
	expired := created.Add(time.Hour * 24 * 30)

	grpcAuthReq := pb.AuthRequest{
		Token: util.GetRandomString(32),
	}
	grpcAuthRes := pb.AuthResponse{
		ID:        uuid.New().String(),
		UserId:    1,
		CreatedAt: timestamppb.New(created),
		ExpiredAt: timestamppb.New(expired),
	}

	grpcReq := pb.GetUserRequest{
		ID: 1,
	}
	grpcRes := pb.GetUserResponse{
		User: &pb.User{
			ID:        1,
			Username:  "test",
			Fullname:  "test",
			Email:     "test",
			Plan:      1,
			CreatedAt: timestamppb.New(created),
			ExpiredAt: timestamppb.New(expired),
			AutoRenew: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().GetUser(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)
	grpc.EXPECT().Authorize(gomock.Any(), gomock.Eq(&grpcAuthReq)).Return(&grpcAuthRes, nil)

	server := NewTestServer(t, grpc)
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	addAuthHeader(request, grpcAuthReq.Token)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	var res User
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)

	require.Equal(t, int32(1), res.ID)
	require.Equal(t, "test", res.Username)
	require.Equal(t, "test", res.Fullname)
	require.Equal(t, "test", res.Email)
	require.Equal(t, int32(1), res.Plan)
	require.Equal(t, created, res.CreatedAt)
	require.Equal(t, expired, res.ExpiredAt)
	require.True(t, res.AutoRenew)
}

type eqCreateUserRequestMatcher struct {
	req *pb.CreateUserRequest
}

func (m eqCreateUserRequestMatcher) Matches(x interface{}) bool {
	req, ok := x.(*pb.CreateUserRequest)
	if !ok {
		return false
	}
	if m.req.GetUsername() != req.GetUsername() {
		return false
	}
	if m.req.GetFullname() != req.GetFullname() {
		return false
	}
	if m.req.GetEmail() != req.GetEmail() {
		return false
	}
	if m.req.GetPlan() != req.GetPlan() {
		return false
	}
	if m.req.GetAutoRenew() != req.GetAutoRenew() {
		return false
	}
	return true
}

func (m eqCreateUserRequestMatcher) String() string {
	return "matches create user request"
}

func EqCreateUserRequest(req *pb.CreateUserRequest) gomock.Matcher {
	return eqCreateUserRequestMatcher{req: req}
}

func TestCreateUserAPI(t *testing.T) {
	url := "/user/create"
	data, err := json.Marshal(gin.H{
		"username":   "test",
		"password":   "test",
		"fullname":   "test",
		"email":      "test",
		"plan":       1,
		"auto_renew": true,
	})
	require.NoError(t, err)

	created := time.Now().UTC().Truncate(time.Second)
	expired := created.Add(time.Hour * 24 * 30)

	grpcReq := pb.CreateUserRequest{
		Username:  "test",
		Fullname:  "test",
		Email:     "test",
		Plan:      1,
		AutoRenew: true,
	}
	grpcRes := pb.CreateUserResponse{
		User: &pb.User{
			ID:        1,
			Username:  "test",
			Fullname:  "test",
			Email:     "test",
			Plan:      1,
			CreatedAt: timestamppb.New(created),
			ExpiredAt: timestamppb.New(expired),
			AutoRenew: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(&pb.GetUserByUsernameRequest{Username: "test"})).Return(
		nil, api_error.NewAPIError(http.StatusNotFound, "user not found"),
	)
	grpc.EXPECT().CreateUser(gomock.Any(), EqCreateUserRequest(&grpcReq)).Return(&grpcRes, nil)

	server := NewTestServer(t, grpc)
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	var res User
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, int32(1), res.ID)
	require.Equal(t, "test", res.Username)
	require.Equal(t, "test", res.Fullname)
	require.Equal(t, "test", res.Email)
	require.Equal(t, int32(1), res.Plan)
	require.Equal(t, created, res.CreatedAt)
	require.Equal(t, expired, res.ExpiredAt)
	require.True(t, res.AutoRenew)
}

func TestUpdateUserAPI(t *testing.T) {
	url := "/user/update"
	data, err := json.Marshal(gin.H{
		"id":       1,
		"username": "test",
	})
	require.NoError(t, err)

	username := "test"
	created := time.Now().UTC().Truncate(time.Second)
	expired := created.Add(time.Hour * 24 * 30)
	grpcAuthReq := pb.AuthRequest{
		Token: util.GetRandomString(32),
	}
	grpcAuthRes := pb.AuthResponse{
		ID:        uuid.New().String(),
		UserId:    1,
		CreatedAt: timestamppb.New(created),
		ExpiredAt: timestamppb.New(expired),
	}
	grpcReq := pb.UpdateUserRequest{
		ID:       1,
		Username: &username,
	}
	grpcRes := pb.UpdateUserResponse{
		User: &pb.User{
			ID:        1,
			Username:  "test",
			Fullname:  "test",
			Email:     "test",
			Plan:      1,
			CreatedAt: timestamppb.New(created),
			ExpiredAt: timestamppb.New(expired),
			AutoRenew: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().Authorize(gomock.Any(), gomock.Eq(&grpcAuthReq)).Return(&grpcAuthRes, nil)
	grpc.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

	server := NewTestServer(t, grpc)
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	addAuthHeader(request, grpcAuthReq.Token)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	var res User
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, int32(1), res.ID)
	require.Equal(t, "test", res.Username)
	require.Equal(t, "test", res.Fullname)
	require.Equal(t, "test", res.Email)
	require.Equal(t, int32(1), res.Plan)
	require.Equal(t, created, res.CreatedAt)
	require.Equal(t, expired, res.ExpiredAt)
	require.True(t, res.AutoRenew)
}
