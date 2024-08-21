package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNone(t *testing.T) {
	data := None[int]()
	assert.False(t, data.IsSome())
}

func TestSome(t *testing.T) {
	data := Some(1)
	assert.True(t, data.IsSome())
}

func TestUnwrap(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		data := Some(1)
		assert.EqualValues(t, 1, data.Unwrap())
	})

	t.Run("panic", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()

		data := None[int]()
		_ = data.Unwrap()
	})
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		opt       Option[any]
		predicate func(any) bool
		want      Option[any]
	}{
		{
			opt: Some(any(1)),
			predicate: func(a any) bool {
				return true
			},
			want: Some(any(1)),
		},
		{
			opt: Some(any(1)),
			predicate: func(a any) bool {
				return false
			},
			want: None[any](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opt.Filter(tt.predicate)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
