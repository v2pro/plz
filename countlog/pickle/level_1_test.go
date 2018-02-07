package pickle_test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/countlog/pickle"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/test/should"
)

func Test_level1(t *testing.T) {
	t.Run("single ptr in array", test.Case(func(ctx *countlog.Context) {
		type TestObject [1]*uint8
		one := uint8(1)
		obj := TestObject{&one}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x1,
		}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
	}))
	t.Run("two ptrs in array", test.Case(func(ctx *countlog.Context) {
		type TestObject [2]*uint8
		one := uint8(1)
		obj := TestObject{&one, &one}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			1,
			1,
		}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(obj, decoded)
	}))
	t.Run("byte slice", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, []byte("hello"))[0].([]byte)
		must.Equal([]byte{
			0x18, 0, 0, 0, 0, 0, 0, 0,
			5, 0, 0, 0, 0, 0, 0, 0,
			5, 0, 0, 0, 0, 0, 0, 0,
			'h', 'e', 'l', 'l', 'o'}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal([]byte("hello"), decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal([]byte("hello"), decoded)
	}))
	t.Run("int slice", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, []int{1, 2, 3})[0].([]byte)
		must.Equal([]byte{
			0x18, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0,
			1, 0, 0, 0, 0, 0, 0, 0,
			2, 0, 0, 0, 0, 0, 0, 0,
			3, 0, 0, 0, 0, 0, 0, 0}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal([]int{1, 2, 3}, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal([]int{1, 2, 3}, decoded)
	}))
	t.Run("ptr int", test.Case(func(ctx *countlog.Context) {
		val := 100
		encoded := must.Call(pickle.Marshal, &val)[0].([]byte)
		must.Equal([]byte{0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 100, 0, 0, 0, 0, 0, 0, 0},
			encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal(&val, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(&val, decoded)
	}))
	t.Run("string", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, "hello")[0].([]byte)
		must.Equal([]byte{
			0x10, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0,
			'h', 'e', 'l', 'l', 'o'}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		must.Equal("hello", decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal("hello", decoded)
	}))
	t.Run("single ptr in struct", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field1 *uint8
		}
		one := uint8(1)
		obj := TestObject{&one}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x1,
		}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		should.Equal(obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		should.Equal(obj, decoded)
	}))
	t.Run("two ptrs in struct", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field1 *uint8
			Field2 *uint8
		}
		one := uint8(1)
		obj := TestObject{&one, &one}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			1,
			1,
		}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		should.Equal(obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		should.Equal(obj, decoded)
	}))
	t.Run("slice in struct", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field1 uint
			Field2 []uint
		}
		obj := TestObject{Field1: 1, Field2: []uint{2,3}}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x18, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		}, encoded[16:])
		decoded := must.Call(pickle.ReadonlyConfig.Unmarshal, encoded)[0]
		should.Equal(obj, decoded)
		decoded = must.Call(pickle.Unmarshal, encoded)[0]
		should.Equal(obj, decoded)
	}))
	t.Run("multiple struct", test.Case(func(ctx *countlog.Context) {
		type TestObj struct {
			Field []int
		}
		stream := pickle.NewStream(nil)
		stream.Marshal(TestObj{[]int{1}})
		stream.Marshal(TestObj{[]int{2}})
		must.Nil(stream.Error)
		iter := pickle.NewIterator(stream.Buffer())
		obj := iter.Unmarshal()
		must.Nil(iter.Error)
		must.Equal([]int{1}, obj.(TestObj).Field)
		obj = iter.Unmarshal()
		must.Nil(iter.Error)
		must.Equal([]int{2}, obj.(TestObj).Field)
	}))
}
