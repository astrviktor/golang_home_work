package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Slices struct {
		Int    []int    `validate:"in:1,2,3"`
		String []string `validate:"in:a,b,c"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          int64(5),
			expectedErr: ErrProgramError,
		},
		{
			in: User{
				ID:     "12345",
				Name:   "Nick",
				Age:    10,
				Email:  "1@1",
				Role:   "stuff",
				Phones: []string{"79000000000", "79000000000"},
				meta:   []byte{},
			},
			expectedErr: ErrValidationError,
		},
		{
			in:          App{Version: "1.2.3"},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte{1},
				Payload:   []byte{2},
				Signature: []byte{3},
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 400,
				Body: "",
			},
			expectedErr: ErrValidationError,
		},
		{
			in: Slices{
				Int:    []int{1, 2, 3},
				String: []string{"a", "b", "c"},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.True(t, errors.Is(err, tt.expectedErr))
		})
	}
}
