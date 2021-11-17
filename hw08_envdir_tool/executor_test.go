package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("executor check with Environment", func(t *testing.T) {
		result := RunCmd([]string{"printenv", "TESTDATA"},
			Environment{"TESTDATA": EnvValue{Value: "testdata/env", NeedRemove: false}})
		require.Equal(t, result, 0)
	})

	t.Run("executor check without Environment", func(t *testing.T) {
		result := RunCmd([]string{"ls", "testdata/env"}, nil)
		require.Equal(t, result, 0)
	})

	t.Run("executor check empty", func(t *testing.T) {
		result := RunCmd([]string{}, nil)
		require.Equal(t, result, 1)
	})
}
