package bindown

import (
	"crypto/sha256"
	"hash/fnv"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/willabides/bindown/v2/internal/testutil"
	"github.com/willabides/bindown/v2/internal/util"
)

func Test_fileExistsWithChecksum(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		file := filepath.Join(testutil.TmpDir(t), "myfile")
		require.NoError(t, util.CopyFile(testutil.DownloadablesPath("foo.tar.gz"), file, nil))
		got, err := fileExistsWithChecksum(file, testutil.FooChecksum)
		require.NoError(t, err)
		require.True(t, got)
	})

	t.Run("wrong checksum", func(t *testing.T) {
		file := filepath.Join(testutil.TmpDir(t), "myfile")
		checksum := "0000000000000000000000000000000000000000000000000000000000000000"
		require.NoError(t, util.CopyFile(testutil.DownloadablesPath("foo.tar.gz"), file, nil))
		got, err := fileExistsWithChecksum(file, checksum)
		require.NoError(t, err)
		require.False(t, got)
	})

	t.Run("doesn't exist", func(t *testing.T) {
		file := filepath.Join(testutil.TmpDir(t), "myfile")
		got, err := fileExistsWithChecksum(file, testutil.FooChecksum)
		require.NoError(t, err)
		require.False(t, got)
	})
}

func Test_fileChecksum(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		file := filepath.Join(testutil.TmpDir(t), "myfile")
		require.NoError(t, util.CopyFile(testutil.DownloadablesPath("foo.tar.gz"), file, nil))
		got, err := fileChecksum(file)
		require.NoError(t, err)
		require.Equal(t, testutil.FooChecksum, got)
	})

	t.Run("doesn't exist", func(t *testing.T) {
		file := filepath.Join(testutil.TmpDir(t), "myfile")
		got, err := fileChecksum(file)
		require.Error(t, err)
		require.Empty(t, got)
	})
}

func Test_hexHash(t *testing.T) {
	got, err := hexHash(fnv.New64a(), []byte("foo"))
	require.NoError(t, err)
	require.Equal(t, "dcb27518fed9d577", got)
	got, err = hexHash(fnv.New64a(), []byte("foo"), []byte("bar"))
	require.NoError(t, err)
	require.Equal(t, "85944171f73967e8", got)
	content := testutil.MustReadFile(t, testutil.DownloadablesPath("foo.tar.gz"))
	got, err = hexHash(sha256.New(), content)
	require.NoError(t, err)
	require.Equal(t, testutil.FooChecksum, got)
}
