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
	"github.com/machearn/galaxy_controller/pb"
	mockpb "github.com/machearn/galaxy_controller/pb/mock"
	"github.com/machearn/galaxy_controller/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRenewToken(t *testing.T) {
	url := "/token/renew"
	refreshToken := util.GetRandomString(32)
	newAccessToken := util.GetRandomString(32)
	expired := time.Now().UTC().Truncate(time.Second)

	data, err := json.Marshal(gin.H{
		"refresh_token": refreshToken,
	})
	require.NoError(t, err)

	grpcReq := pb.RenewAccessTokenRequest{
		RefreshToken: refreshToken,
	}
	grpcRes := pb.RenewAccessTokenResponse{
		AccessToken: newAccessToken,
		ExpiredAt:   timestamppb.New(expired),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().RenewAccessToken(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

	server := NewTestServer(t, grpc)
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	var res RenewAccessTokenResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)
	require.Equal(t, newAccessToken, res.AccessToken)
	require.Equal(t, expired, res.AccessExpiredAt)
}
