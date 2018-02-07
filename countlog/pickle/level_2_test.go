package pickle_test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/countlog/pickle"
	"github.com/v2pro/plz/test/must"
)

func Test_level2(t *testing.T) {
	t.Run("ptr of string", test.Case(func(ctx *countlog.Context) {
		obj := "hello"
		encoded := must.Call(pickle.Marshal, &obj)[0].([]byte)
		must.Equal([]byte{
			0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x68, 0x65, 0x6c, 0x6c, 0x6f,
		}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal(&obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(&obj, decoded)
	}))
	t.Run("ptr of slice", test.Case(func(ctx *countlog.Context) {
		obj := []byte("hello")
		encoded := must.Call(pickle.Marshal, &obj)[0].([]byte)
		must.Equal([]byte{
			0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x18, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x68, 0x65, 0x6c, 0x6c, 0x6f,
		}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal(&obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(&obj, decoded)
	}))
	t.Run("string slice", test.Case(func(ctx *countlog.Context) {
		obj := []string{"h", "i"}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x18, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, // sliceHeader
			0x20, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,                         // string header
			0x11, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,                         // string header
			'h', 'i'}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
	}))
	t.Run("ptr slice", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field1 int
			Field2 int
		}
		obj := []*TestObject{{1, 2}, {3, 4}}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			24, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,
			16, 0, 0, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0, 0,
			1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, // [0]
			3, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, // [1]
		}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
	}))
	t.Run("nil struct within struct", test.Case(func(ctx *countlog.Context) {
		type SubObject struct {
			length uint
			set    []uint64
		}
		type TestObject struct {
			f1 uint
			f2 uint
			f3 *SubObject
		}
		obj := TestObject{}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
	}))
	t.Run("struct within struct", test.Case(func(ctx *countlog.Context) {
		type SubObject struct {
			length uint
			set    []uint64
		}
		type TestObject struct {
			f1 uint
			f2 uint
			f3 *SubObject
		}
		obj := TestObject{f1: 1, f2: 2, f3: &SubObject{length: 3, set: []uint64{100}}}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			24, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
	}))
	t.Run("multiple structs", test.Case(func(ctx *countlog.Context) {
		type TestObj struct {
			Field1 []int
			Field2 [][]byte
		}
		stream := pickle.NewStream(nil)
		stream.Marshal(TestObj{Field2: [][]byte{[]byte("hello")}})
		stream.Marshal(TestObj{Field2: [][]byte{[]byte("world")}})
		must.Nil(stream.Error)
		iter := pickle.NewIterator(stream.Buffer())
		obj := iter.Unmarshal()
		must.Nil(iter.Error)
		must.Equal([][]byte{[]byte("hello")}, obj.(TestObj).Field2)
		obj = iter.Unmarshal()
		must.Nil(iter.Error)
		must.Equal([][]byte{[]byte("world")}, obj.(TestObj).Field2)
	}))
}
