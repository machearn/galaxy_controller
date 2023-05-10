package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/machearn/galaxy_controller/pb"
	mockpb "github.com/machearn/galaxy_controller/pb/mock"
	"github.com/machearn/galaxy_controller/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateItem(t *testing.T) {
	req := CreateItemRequest{
		Name:     "test",
		Quantity: 1,
		Price:    1,
	}

	item := Item{
		ID:       1,
		Name:     req.Name,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	grpcReq := pb.CreateItemRequest{
		Name:     req.Name,
		Quantity: req.Quantity,
		Price:    req.Price,
	}
	grpcRes := pb.CreateItemResponse{
		Item: &pb.Item{
			ID:       1,
			Name:     req.Name,
			Quantity: req.Quantity,
			Price:    req.Price,
		},
	}

	createdAt := time.Now().UTC().Truncate(time.Second)
	grpcAuthReq := pb.AuthRequest{
		Token: util.GetRandomString(32),
	}
	grpcAuthRes := pb.AuthResponse{
		ID:        uuid.New().String(),
		UserId:    1,
		CreatedAt: timestamppb.New(createdAt),
		ExpiredAt: timestamppb.New(createdAt),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().CreateItem(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)
	grpc.EXPECT().Authorize(gomock.Any(), gomock.Eq(&grpcAuthReq)).Return(&grpcAuthRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	data, err := json.Marshal(gin.H{
		"name":     req.Name,
		"quantity": req.Quantity,
		"price":    req.Price,
	})
	require.NoError(t, err)

	url := "/item/create"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	addAuthHeader(request, grpcAuthReq.Token)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err = io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var res Item
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.Equal(t, item, res)
}

func TestGetItem(t *testing.T) {
	req := GetItemRequest{
		ID: 1,
	}
	item := Item{
		ID:       1,
		Name:     "test",
		Quantity: 1,
		Price:    1,
	}

	grpcReq := pb.GetItemRequest{
		Id: req.ID,
	}
	grpcRes := pb.GetItemResponse{
		Item: &pb.Item{
			ID:       1,
			Name:     "test",
			Quantity: 1,
			Price:    1,
		},
	}

	createdAt := time.Now().UTC().Truncate(time.Second)
	grpcAuthReq := pb.AuthRequest{
		Token: util.GetRandomString(32),
	}
	grpcAuthRes := pb.AuthResponse{
		ID:        uuid.New().String(),
		UserId:    1,
		CreatedAt: timestamppb.New(createdAt),
		ExpiredAt: timestamppb.New(createdAt),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().GetItem(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)
	grpc.EXPECT().Authorize(gomock.Any(), gomock.Eq(&grpcAuthReq)).Return(&grpcAuthRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	url := "/item/get/1"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	addAuthHeader(request, grpcAuthReq.Token)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err := io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var res Item
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.Equal(t, item, res)
}

func TestListItems(t *testing.T) {
	req := ListItemsRequest{
		Offset: 0,
		Limit:  5,
	}

	var items []Item
	for i := 1; i <= 5; i++ {
		items = append(items, Item{
			ID:       int32(i),
			Name:     "test",
			Quantity: 1,
			Price:    1,
		})
	}
	result := ListItemsResponse{
		Items: items,
	}

	grpcReq := pb.ListItemsRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	var grpcItems []*pb.Item
	for i := 1; i <= 5; i++ {
		grpcItems = append(grpcItems, &pb.Item{
			ID:       int32(i),
			Name:     "test",
			Quantity: 1,
			Price:    1,
		})
	}

	grpcRes := pb.ListItemsResponse{
		Items: grpcItems,
	}

	createdAt := time.Now().UTC().Truncate(time.Second)
	grpcAuthReq := pb.AuthRequest{
		Token: util.GetRandomString(32),
	}
	grpcAuthRes := pb.AuthResponse{
		ID:        uuid.New().String(),
		UserId:    1,
		CreatedAt: timestamppb.New(createdAt),
		ExpiredAt: timestamppb.New(createdAt),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().ListItems(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)
	grpc.EXPECT().Authorize(gomock.Any(), gomock.Eq(&grpcAuthReq)).Return(&grpcAuthRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	url := "/item/list"
	data, err := json.Marshal(gin.H{
		"offset": req.Offset,
		"limit":  req.Limit,
	})
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	addAuthHeader(request, grpcAuthReq.Token)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err = io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var res ListItemsResponse
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)

	require.Equal(t, result, res)
}

func TestUpdateItem(t *testing.T) {
	url := "/item/update"

	data, err := json.Marshal(gin.H{
		"id":       1,
		"name":     "test",
		"quantity": nil,
	})
	require.NoError(t, err)

	name := "test"
	grpcReq := pb.UpdateItemRequest{
		Id:       1,
		Name:     &name,
		Quantity: nil,
	}

	grpcRes := pb.UpdateItemResponse{
		Item: &pb.Item{
			ID:       1,
			Name:     name,
			Quantity: 1,
			Price:    1,
		},
	}

	createdAt := time.Now().UTC().Truncate(time.Second)
	grpcAuthReq := pb.AuthRequest{
		Token: util.GetRandomString(32),
	}
	grpcAuthRes := pb.AuthResponse{
		ID:        uuid.New().String(),
		UserId:    1,
		CreatedAt: timestamppb.New(createdAt),
		ExpiredAt: timestamppb.New(createdAt),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().UpdateItem(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)
	grpc.EXPECT().Authorize(gomock.Any(), gomock.Eq(&grpcAuthReq)).Return(&grpcAuthRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	addAuthHeader(request, grpcAuthReq.Token)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err = io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var res Item
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)

	require.Equal(t, grpcRes.Item.ID, res.ID)
	require.Equal(t, grpcRes.Item.Name, res.Name)
	require.Equal(t, grpcRes.Item.Quantity, res.Quantity)
	require.Equal(t, grpcRes.Item.Price, res.Price)
}

func TestDeleteItem(t *testing.T) {
	url := "/item/delete/1"

	grpcReq := pb.DeleteItemRequest{
		Id: 1,
	}

	createdAt := time.Now().UTC().Truncate(time.Second)
	grpcAuthReq := pb.AuthRequest{
		Token: util.GetRandomString(32),
	}
	grpcAuthRes := pb.AuthResponse{
		ID:        uuid.New().String(),
		UserId:    1,
		CreatedAt: timestamppb.New(createdAt),
		ExpiredAt: timestamppb.New(createdAt),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().DeleteItem(gomock.Any(), gomock.Eq(&grpcReq)).Return(&pb.Empty{}, nil)
	grpc.EXPECT().Authorize(gomock.Any(), gomock.Eq(&grpcAuthReq)).Return(&grpcAuthRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	addAuthHeader(request, grpcAuthReq.Token)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)
}
