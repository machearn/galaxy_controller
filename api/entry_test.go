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
	"github.com/machearn/galaxy_controller/pb"
	mockpb "github.com/machearn/galaxy_controller/pb/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateEntry(t *testing.T) {
	url := "/entry/create"
	data, err := json.Marshal(gin.H{
		"member_id": 1,
		"item_id":   1,
		"quantity":  1,
		"total":     1,
	})
	require.NoError(t, err)

	grpcReq := pb.CreateEntryRequest{
		UserId:   1,
		ItemId:   1,
		Quantity: 1,
		Total:    1,
	}
	createAt := time.Now().UTC().Truncate(time.Second)
	grpcRes := pb.CreateEntryResponse{
		Entry: &pb.Entry{
			ID:        1,
			UserId:    1,
			ItemId:    1,
			Quantity:  1,
			Total:     1,
			CreatedAt: timestamppb.New(createAt),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, nil)
	grpc.EXPECT().GetItem(gomock.Any(), gomock.Any()).Return(nil, nil)
	grpc.EXPECT().CreateEntry(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err = io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var entry Entry
	err = json.Unmarshal(data, &entry)
	require.NoError(t, err)

	require.Equal(t, int32(1), entry.ID)
	require.Equal(t, int32(1), entry.UserID)
	require.Equal(t, int32(1), entry.ItemID)
	require.Equal(t, int32(1), entry.Quantity)
	require.Equal(t, int32(1), entry.Total)
	require.Equal(t, createAt, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	url := "/entry/get/1"

	grpcReq := pb.GetEntryRequest{
		Id: 1,
	}

	createdAt := time.Now().UTC().Truncate(time.Second)
	grpcRes := pb.GetEntryResponse{
		Entry: &pb.Entry{
			ID:        1,
			UserId:    1,
			ItemId:    1,
			Quantity:  1,
			Total:     1,
			CreatedAt: timestamppb.New(createdAt),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().GetEntry(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err := io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var entry Entry
	err = json.Unmarshal(data, &entry)
	require.NoError(t, err)

	require.Equal(t, int32(1), entry.ID)
	require.Equal(t, int32(1), entry.UserID)
	require.Equal(t, int32(1), entry.ItemID)
	require.Equal(t, int32(1), entry.Quantity)
	require.Equal(t, int32(1), entry.Total)
	require.Equal(t, createdAt, entry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	url := "/entry/list"
	data, err := json.Marshal(gin.H{
		"offset": 0,
		"limit":  10,
	})
	require.NoError(t, err)

	grpcReq := pb.ListEntriesRequest{
		Offset: 0,
		Limit:  10,
	}

	grpcEntries := make([]*pb.Entry, 10)
	for i := 0; i < 10; i++ {
		createdAt := time.Now().UTC().Truncate(time.Second)
		grpcEntries[i] = &pb.Entry{
			ID:        int32(i + 1),
			UserId:    1,
			ItemId:    1,
			Quantity:  1,
			Total:     1,
			CreatedAt: timestamppb.New(createdAt),
		}
	}

	grpcRes := pb.ListEntriesResponse{
		Entries: grpcEntries,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().ListEntries(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err = io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var res ListEntriesResponse
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.Len(t, res.Entries, 10)

	for i := 0; i < 10; i++ {
		require.Equal(t, int32(i+1), res.Entries[i].ID)
		require.Equal(t, int32(1), res.Entries[i].UserID)
		require.Equal(t, int32(1), res.Entries[i].ItemID)
		require.Equal(t, int32(1), res.Entries[i].Quantity)
		require.Equal(t, int32(1), res.Entries[i].Total)
		require.Equal(t, grpcEntries[i].CreatedAt.AsTime(), res.Entries[i].CreatedAt)
	}
}

func TestListEntriesByUser(t *testing.T) {
	url := "/entry/list/user"
	data, err := json.Marshal(gin.H{
		"user_id": 1,
		"offset":  0,
		"limit":   10,
	})
	require.NoError(t, err)

	grpcReq := pb.ListEntriesByUserRequest{
		UserId: 1,
		Offset: 0,
		Limit:  10,
	}

	grpcEntries := make([]*pb.Entry, 10)
	for i := 0; i < 10; i++ {
		createdAt := time.Now().UTC().Truncate(time.Second)
		grpcEntries[i] = &pb.Entry{
			ID:        int32(i + 1),
			UserId:    1,
			ItemId:    1,
			Quantity:  1,
			Total:     1,
			CreatedAt: timestamppb.New(createdAt),
		}
	}

	grpcRes := pb.ListEntriesResponse{
		Entries: grpcEntries,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().ListEntriesByUser(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err = io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var res ListEntriesResponse
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.Len(t, res.Entries, 10)

	for i := 0; i < 10; i++ {
		require.Equal(t, int32(i+1), res.Entries[i].ID)
		require.Equal(t, int32(1), res.Entries[i].UserID)
		require.Equal(t, int32(1), res.Entries[i].ItemID)
		require.Equal(t, int32(1), res.Entries[i].Quantity)
		require.Equal(t, int32(1), res.Entries[i].Total)
		require.Equal(t, grpcEntries[i].CreatedAt.AsTime(), res.Entries[i].CreatedAt)
	}
}

func TestListEntriesByItem(t *testing.T) {
	url := "/entry/list/item"
	data, err := json.Marshal(gin.H{
		"item_id": 1,
		"offset":  0,
		"limit":   10,
	})
	require.NoError(t, err)

	grpcReq := pb.ListEntriesByItemRequest{
		ItemId: 1,
		Offset: 0,
		Limit:  10,
	}

	grpcEntries := make([]*pb.Entry, 10)
	for i := 0; i < 10; i++ {
		createdAt := time.Now().UTC().Truncate(time.Second)
		grpcEntries[i] = &pb.Entry{
			ID:        int32(i + 1),
			UserId:    1,
			ItemId:    1,
			Quantity:  1,
			Total:     1,
			CreatedAt: timestamppb.New(createdAt),
		}
	}

	grpcRes := pb.ListEntriesResponse{
		Entries: grpcEntries,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpc := mockpb.NewMockGalaxyClient(ctrl)
	grpc.EXPECT().ListEntriesByItem(gomock.Any(), gomock.Eq(&grpcReq)).Return(&grpcRes, nil)

	server := NewTestServer(t, grpc)
	recoder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recoder, request)
	require.Equal(t, http.StatusOK, recoder.Code)

	data, err = io.ReadAll(recoder.Body)
	require.NoError(t, err)

	var res ListEntriesResponse
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.Len(t, res.Entries, 10)

	for i := 0; i < 10; i++ {
		require.Equal(t, int32(i+1), res.Entries[i].ID)
		require.Equal(t, int32(1), res.Entries[i].UserID)
		require.Equal(t, int32(1), res.Entries[i].ItemID)
		require.Equal(t, int32(1), res.Entries[i].Quantity)
		require.Equal(t, int32(1), res.Entries[i].Total)
		require.Equal(t, grpcEntries[i].CreatedAt.AsTime(), res.Entries[i].CreatedAt)
	}
}
