package api

import (
	"net/http"
	"os"
	"testing"

	"github.com/machearn/galaxy_controller/pb"
	"github.com/machearn/galaxy_controller/util"
	"github.com/stretchr/testify/require"
)

func addAuthHeader(req *http.Request, token string) {
	req.Header.Set("authorization", "bearer "+token)
}

func NewTestServer(t *testing.T, grpc pb.GalaxyClient) *Server {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)
	require.NotEmpty(t, config)

	server, err := NewServer(config, grpc)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
