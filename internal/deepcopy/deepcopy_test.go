package deepcopy

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Namer interface {
	Name() string
}

type test struct {
	name string
	Age  int
}

func (t *test) Name() string {
	return t.name
}

func TestIface(t *testing.T) {
	got := Iface(1)
	assert.EqualValues(t, 1, got)
}

func TestCopy(t *testing.T) {
	type args[T any] struct {
		src T
	}

	type TestCase[T any] struct {
		args args[T]
		want T
	}

	t.Run("int", func(t *testing.T) {
		tt := TestCase[int]{
			args: args[int]{
				src: 1,
			},
			want: 1,
		}

		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got = 2
		require.NotEqualValues(t, tt.args.src, got)
	})

	t.Run("float", func(t *testing.T) {
		tt := TestCase[float64]{
			args: args[float64]{
				src: 1.11,
			},
			want: 1.11,
		}

		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got = 2.22
		require.NotEqualValues(t, tt.args.src, got)
	})

	t.Run("bool", func(t *testing.T) {
		tt := TestCase[bool]{
			args: args[bool]{
				src: true,
			},
			want: true,
		}

		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got = false
		require.NotEqualValues(t, tt.args.src, got)
	})

	t.Run("string", func(t *testing.T) {
		tt := TestCase[string]{
			args: args[string]{
				src: "foo",
			},
			want: "foo",
		}

		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got = "zoo"
		require.NotEqualValues(t, tt.args.src, got)
	})

	t.Run("nil", func(t *testing.T) {
		tt := TestCase[any]{
			args: args[any]{
				src: nil,
			},
			want: nil,
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)
	})

	t.Run("time.Time", func(t *testing.T) {
		tt := TestCase[time.Time]{
			args: args[time.Time]{
				src: time.Unix(0, 0),
			},
			want: time.Unix(0, 0),
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got = got.Add(time.Second)
		require.NotEqualValues(t, tt.args.src, got)
	})

	t.Run("slice nil", func(t *testing.T) {
		tt := TestCase[[]int]{
			args: args[[]int]{
				src: ([]int)(nil),
			},
			want: ([]int)(nil),
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)
	})

	t.Run("map nil", func(t *testing.T) {
		tt := TestCase[map[int]int]{
			args: args[map[int]int]{
				src: (map[int]int)(nil),
			},
			want: (map[int]int)(nil),
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)
	})

	t.Run("object nil", func(t *testing.T) {
		tt := TestCase[*test]{
			args: args[*test]{
				src: (*test)(nil),
			},
			want: (*test)(nil),
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)
	})

	t.Run("map", func(t *testing.T) {
		tt := TestCase[map[string]string]{
			args: args[map[string]string]{
				src: map[string]string{"foo": "bar"},
			},
			want: map[string]string{"foo": "bar"},
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got["foo"] = "zoo"
		require.NotEqualValues(t, tt.args.src, got)

		cp := tt.args.src
		cp["foo"] = "zoo"
		require.EqualValues(t, tt.args.src, cp)
	})

	t.Run("object ptr", func(t *testing.T) {
		tt := TestCase[*test]{
			args: args[*test]{
				src: &test{name: "foo"},
			},
			want: &test{name: "foo"},
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got.Age = 101
		require.NotEqualValues(t, tt.args.src, got)

		cp := tt.args.src
		cp.Age = 101
		require.EqualValues(t, tt.args.src, cp)
	})

	t.Run("object", func(t *testing.T) {
		tt := TestCase[test]{
			args: args[test]{
				src: test{name: ""},
			},
			want: test{name: ""},
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got.Age = 101
		require.NotEqualValues(t, tt.args.src, got)

		cp := tt.args.src
		cp.Age = 101
		require.NotEqualValues(t, tt.args.src, cp)
	})

	t.Run("slice", func(t *testing.T) {
		tt := TestCase[[]string]{
			args: args[[]string]{
				src: []string{"foo"},
			},
			want: []string{"foo"},
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got[0] = "zoo"
		require.NotEqualValues(t, tt.args.src, got)

		cp := tt.args.src
		cp[0] = "zoo"
		require.EqualValues(t, tt.args.src, cp)
	})

	t.Run("slice object", func(t *testing.T) {
		tt := TestCase[[]test]{
			args: args[[]test]{
				src: []test{{name: "foo"}},
			},
			want: []test{{name: "foo"}},
		}
		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got[0] = test{name: "zoo"}
		require.NotEqualValues(t, tt.args.src, got)

		cp := tt.args.src
		cp[0] = test{name: "zoo"}
		require.EqualValues(t, tt.args.src, cp)
	})

	t.Run("slice object ptr", func(t *testing.T) {
		tt := TestCase[[]*test]{
			args: args[[]*test]{
				src: []*test{{name: "foo"}},
			},
			want: []*test{{name: "foo"}},
		}

		got := Copy(tt.args.src)
		require.EqualValues(t, tt.want, got)

		got[0] = &test{name: "zoo"}
		require.NotEqualValues(t, tt.args.src, got)

		cp := tt.args.src
		cp[0] = &test{name: "zoo"}
		require.EqualValues(t, tt.args.src, cp)
	})

	// not supported
	t.Run("chan", func(t *testing.T) {
		tt := TestCase[chan int]{
			args: args[chan int]{
				src: make(chan int, 1),
			},
		}
		defer func() {
			assert.NotNil(t, recover())
		}()
		_ = Copy(tt.args.src)
	})

	t.Run("func", func(t *testing.T) {
		tt := TestCase[func()]{
			args: args[func()]{
				src: func() {},
			},
		}
		defer func() {
			assert.NotNil(t, recover())
		}()
		_ = Copy(tt.args.src)
	})
}

func Test_copyRecursive(t *testing.T) {
	t.Run("invalid ptr", func(t *testing.T) {
		var ptr *test
		original := reflect.ValueOf(ptr)
		cpy := reflect.New(original.Type()).Elem()
		copyRecursive[*test](original, cpy)

		assert.EqualValues(t, ptr, cpy.Interface().(*test))
	})

	t.Run("slice nil", func(t *testing.T) {
		var ptr []int
		original := reflect.ValueOf(ptr)
		cpy := reflect.New(original.Type()).Elem()
		copyRecursive[*test](original, cpy)

		assert.EqualValues(t, ptr, cpy.Interface().([]int))
	})

	t.Run("map nil", func(t *testing.T) {
		var ptr map[int]int
		original := reflect.ValueOf(ptr)
		cpy := reflect.New(original.Type()).Elem()
		copyRecursive[*test](original, cpy)

		assert.EqualValues(t, ptr, cpy.Interface().(map[int]int))
	})
}

type MyInt int

func (m MyInt) DeepCopy() MyInt {
	return MyInt(123456789)
}

var _ Interface[MyInt] = (*MyInt)(nil)

func TestDeepCopyInterface(t *testing.T) {
	assert.EqualValues(t, 1, Copy(1))

	data := MyInt(10)
	assert.EqualValues(t, MyInt(123456789), Copy(data))
}
