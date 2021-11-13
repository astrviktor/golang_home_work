package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("test open testdata/input.txt", func(t *testing.T) {
		f, err := os.Open("testdata/input.txt")
		if err != nil {
			t.Fail()
		}
		defer f.Close()
		require.NoError(t, err)
	})

	tests := []struct {
		fromPath string
		toPath   string
		offset   int64
		limit    int64
	}{
		{fromPath: "input.txt", toPath: "out_offset0_limit0.txt", offset: 0, limit: 0},
		{fromPath: "input.txt", toPath: "out_offset0_limit10.txt", offset: 0, limit: 10},
		{fromPath: "input.txt", toPath: "out_offset0_limit1000.txt", offset: 0, limit: 1000},
		{fromPath: "input.txt", toPath: "out_offset0_limit10000.txt", offset: 0, limit: 10000},
		{fromPath: "input.txt", toPath: "out_offset100_limit1000.txt", offset: 100, limit: 1000},
		{fromPath: "input.txt", toPath: "out_offset6000_limit1000.txt", offset: 6000, limit: 1000},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.toPath, func(t *testing.T) {
			err := Copy(filepath.Join("testdata/", tc.fromPath), filepath.Join("/tmp/", tc.toPath), tc.offset, tc.limit)
			require.NoError(t, err)

			testData, err := ioutil.ReadFile(filepath.Join("testdata/", tc.toPath))
			require.NoError(t, err)

			resultData, err := ioutil.ReadFile(filepath.Join("/tmp/", tc.toPath))
			require.NoError(t, err)

			require.Equal(t, testData, resultData)

			err = os.Remove(filepath.Join("/tmp/", tc.toPath))
			require.NoError(t, err)
		})
	}

	t.Run("test copy /dev/random", func(t *testing.T) {
		err := Copy("/dev/random", "/tmp/random.txt", 0, 10)

		require.Error(t, err)
		require.True(t, errors.Is(err, ErrEmptySize))
	})

	t.Run("test empty paths", func(t *testing.T) {
		err := Copy("", "", 0, 0)

		require.Error(t, err)
		require.True(t, errors.Is(err, ErrEmptyPath))

		err = Copy("", "/tmp/123", 0, 0)

		require.Error(t, err)
		require.True(t, errors.Is(err, ErrEmptyPath))

		err = Copy("/tmp/123", "", 0, 0)

		require.Error(t, err)
		require.True(t, errors.Is(err, ErrEmptyPath))
	})

	t.Run("test offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out_offset10000_limit0.txt", 10000, 0)

		require.Error(t, err)
		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize))
	})

	t.Run("test dir instead of file", func(t *testing.T) {
		err := Copy("testdata", "/tmp/out_offset0_limit0.txt", 0, 0)

		require.Error(t, err)
		require.True(t, errors.Is(err, ErrUnsupportedFile))
	})
}
