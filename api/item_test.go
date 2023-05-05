package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/machearn/galaxy_controller/pb"
	mockpb "github.com/machearn/galaxy_controller/pb/mock"
	"github.com/stretchr/testify/require"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().CreateItem(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().GetItem(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	url := "/item/get/1"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().ListItems(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

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

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err = io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var res ListItemsResponse
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)

	require.Equal(t, result, res)
}
