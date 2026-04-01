package storage

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStorageServer(t *testing.T) {
	t.Setenv("SPACES_BUCKET_NAME", "bucket-1")
	t.Setenv("SPACES_REGION", "sgp1")
	t.Setenv("SPACES_ENDPOINT", "https://objects.example.test")
	t.Setenv("SPACES_CDN_ENDPOINT", "https://cdn.example.test")
	t.Setenv("SPACES_ROOT_FOLDER", "root")
	t.Setenv("SPACES_ACCESS_KEY_ID", "key")
	t.Setenv("SPACES_SECRET_ACCESS_KEY", "secret")

	server, err := NewStorageServer(&sql.DB{})
	require.NoError(t, err)
	require.NotNil(t, server)

	server, err = NewStorageServerWithProducer(&sql.DB{}, nil)
	require.NoError(t, err)
	require.NotNil(t, server)
}
