package gqlerror_test

import (
	"reflect"
	"testing"

	"github.com/movio/gqlt/gqlparser/gqlerror"
	"github.com/stretchr/testify/require"
)

type testError struct {
	message string
}

func (e testError) Error() string {
	return e.message
}

var (
	underlyingError = testError{
		"Underlying error",
	}

	error1 = &gqlerror.Error{
		Message: "Some error 1",
	}
	error2 = &gqlerror.Error{
		Err:     underlyingError,
		Message: "Some error 2",
	}
)

func TestErrorFormatting(t *testing.T) {
	t.Run("without filename", func(t *testing.T) {
		err := gqlerror.ErrorLocf("", 66, 2, "kabloom")

		require.Equal(t, `input:66: kabloom`, err.Error())
		require.Equal(t, nil, err.Extensions["file"])
	})

	t.Run("with filename", func(t *testing.T) {
		err := gqlerror.ErrorLocf("schema.graphql", 66, 2, "kabloom")

		require.Equal(t, `schema.graphql:66: kabloom`, err.Error())
		require.Equal(t, "schema.graphql", err.Extensions["file"])
	})
}

func TestList_As(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		errs        gqlerror.List
		target      any
		wantsTarget any
		targetFound bool
	}{
		{
			name: "Empty list",
			errs: gqlerror.List{},
		},
		{
			name:        "List with one error",
			errs:        gqlerror.List{error1},
			target:      new(*gqlerror.Error),
			wantsTarget: &error1,
			targetFound: true,
		},
		{
			name:        "List with multiple errors 1",
			errs:        gqlerror.List{error1, error2},
			target:      new(*gqlerror.Error),
			wantsTarget: &error1,
			targetFound: true,
		},
		{
			name:        "List with multiple errors 2",
			errs:        gqlerror.List{error1, error2},
			target:      new(testError),
			wantsTarget: &underlyingError,
			targetFound: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targetFound := tt.errs.As(tt.target)

			if targetFound != tt.targetFound {
				t.Errorf("List.As() = %v, want %v", targetFound, tt.targetFound)
			}

			if tt.targetFound && !reflect.DeepEqual(tt.target, tt.wantsTarget) {
				t.Errorf("target = %v, want %v", tt.target, tt.wantsTarget)
			}
		})
	}
}

func TestList_Is(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		errs             gqlerror.List
		target           error
		hasMatchingError bool
	}{
		{
			name:             "Empty list",
			errs:             gqlerror.List{},
			target:           new(gqlerror.Error),
			hasMatchingError: false,
		},
		{
			name: "List with one error",
			errs: gqlerror.List{
				error1,
			},
			target:           error1,
			hasMatchingError: true,
		},
		{
			name: "List with multiple errors 1",
			errs: gqlerror.List{
				error1,
				error2,
			},
			target:           error2,
			hasMatchingError: true,
		},
		{
			name: "List with multiple errors 2",
			errs: gqlerror.List{
				error1,
				error2,
			},
			target:           underlyingError,
			hasMatchingError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hasMatchingError := tt.errs.Is(tt.target)
			if hasMatchingError != tt.hasMatchingError {
				t.Errorf("List.Is() = %v, want %v", hasMatchingError, tt.hasMatchingError)
			}
			if hasMatchingError && tt.target == nil {
				t.Errorf("List.Is() returned nil target, wants concrete error")
			}
		})
	}
}
