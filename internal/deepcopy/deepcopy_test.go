package deepcopy

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type test struct {
	name string
	Age  int
}

func TestCopy(t *testing.T) {
	type args struct {
		src any
	}

	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "int",
			args: args{
				src: 1,
			},
			want: 1,
		},
		{
			name: "string",
			args: args{
				src: "foo",
			},
			want: "foo",
		},
		{
			name: "bool",
			args: args{
				src: true,
			},
			want: true,
		},
		{
			name: "map",
			args: args{
				src: map[string]string{"foo": "bar"},
			},
			want: map[string]string{"foo": "bar"},
		},
		{
			name: "nil",
			args: args{
				src: nil,
			},
			want: nil,
		},
		{
			name: "pointer",
			args: args{
				src: &test{name: "foo"},
			},
			want: &test{name: "foo"},
		},
		{
			name: "struct",
			args: args{
				src: test{name: ""},
			},
			want: test{name: ""},
		},
		{
			name: "string slice",
			args: args{
				src: []string{"foo", "bar"},
			},
			want: []string{"foo", "bar"},
		},
		{
			name: "object slice",
			args: args{
				src: []*test{{name: "foo", Age: 10}, {name: "bar", Age: 100}},
			},
			want: []*test{{name: "foo", Age: 10}, {name: "bar", Age: 100}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Copy(tt.args.src)

			require.EqualValues(t, tt.want, got)
			assert.Equal(t, unsafe.Pointer(&got), unsafe.Pointer(&got))
			assert.Equal(t, unsafe.Pointer(&tt.args.src), unsafe.Pointer(&tt.args.src))

			if unsafe.Pointer(&tt.args.src) == unsafe.Pointer(&got) {
				t.Fatalf("expected the pointers to be different; they weren't: src: %v\t copy: %v", unsafe.Pointer(&tt.args.src), unsafe.Pointer(&got))
			}
		})
	}
}
