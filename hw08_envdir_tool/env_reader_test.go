package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("env reader check testdata/env", func(t *testing.T) {
		data := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		result, err := ReadDir("testdata/env")

		require.NoError(t, err)
		require.Equal(t, data, result)
	})

	t.Run("env reader check not dir", func(t *testing.T) {
		_, err := ReadDir("main.go")
		require.Error(t, err)
	})

	t.Run("env reader check dir not exist", func(t *testing.T) {
		_, err := ReadDir("123")
		require.Error(t, err)
	})
}
