package pickle_test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/countlog/pickle"
	"github.com/v2pro/plz/test/should"
)

func Test_level0(t *testing.T) {
	t.Run("int8", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, int8(100))[0].([]byte)
		must.Equal([]byte{100}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(int8(100), decoded)
	}))
	t.Run("uint8", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, uint8(100))[0].([]byte)
		must.Equal([]byte{100}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(uint8(100), decoded)
	}))
	t.Run("int16", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, int16(100))[0].([]byte)
		must.Equal([]byte{100, 0}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(int16(100), decoded)
	}))
	t.Run("uint16", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, uint16(100))[0].([]byte)
		must.Equal([]byte{100, 0}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(uint16(100), decoded)
	}))
	t.Run("int32", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, int32(100))[0].([]byte)
		must.Equal([]byte{100, 0, 0, 0}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(int32(100), decoded)
	}))
	t.Run("uint32", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, uint32(100))[0].([]byte)
		must.Equal([]byte{100, 0, 0, 0}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(uint32(100), decoded)
	}))
	t.Run("int64", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, int64(100))[0].([]byte)
		must.Equal([]byte{100, 0, 0, 0, 0, 0, 0, 0}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(int64(100), decoded)
	}))
	t.Run("uint64", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, uint64(100))[0].([]byte)
		must.Equal([]byte{100, 0, 0, 0, 0, 0, 0, 0}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(uint64(100), decoded)
	}))
	t.Run("int", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, int(100))[0].([]byte)
		must.Equal([]byte{100, 0, 0, 0, 0, 0, 0, 0}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(int(100), decoded)
	}))
	t.Run("uint", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, uint(100))[0].([]byte)
		must.Equal([]byte{100, 0, 0, 0, 0, 0, 0, 0}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(uint(100), decoded)
	}))
	t.Run("uintptr", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, uintptr(100))[0].([]byte)
		must.Equal([]byte{100, 0, 0, 0, 0, 0, 0, 0}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(uintptr(100), decoded)
	}))
	t.Run("float32", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, float32(100))[0].([]byte)
		must.Equal([]byte{0x0, 0x0, 0xc8, 0x42}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(float32(100), decoded)
	}))
	t.Run("float64", test.Case(func(ctx *countlog.Context) {
		encoded := must.Call(pickle.Marshal, float64(100))[0].([]byte)
		must.Equal([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x59, 0x40}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		must.Equal(float64(100), decoded)
	}))
	t.Run("array of int", test.Case(func(ctx *countlog.Context) {
		type TestObject [2]int
		obj := TestObject{1, 2}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		should.Equal(obj, decoded)
	}))
	t.Run("array of struct", test.Case(func(ctx *countlog.Context) {
		type TestStruct struct {
			Field1 int
			Field2 int
		}
		type TestObject [1]TestStruct
		obj := TestObject{TestStruct{1, 2}}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		should.Equal(obj, decoded)
	}))
	t.Run("simple struct", test.Case(func(ctx *countlog.Context) {
		type TestObject struct {
			Field1 int
			Field2 int
		}
		obj := TestObject{1, 2}
		encoded := must.Call(pickle.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		}, encoded[12:])
		decoded := must.Call(pickle.Unmarshal, encoded)[0]
		should.Equal(obj, decoded)
	}))
	t.Run("multiple struct", test.Case(func(ctx *countlog.Context) {
		type TestObj struct {
			Field int
		}
		stream := pickle.NewStream(nil)
		stream.Marshal(TestObj{1})
		stream.Marshal(TestObj{2})
		must.Nil(stream.Error)
		iter := pickle.NewIterator(stream.Buffer())
		obj := iter.Unmarshal()
		must.Nil(iter.Error)
		must.Equal(1, obj.(TestObj).Field)
		obj = iter.Unmarshal()
		must.Nil(iter.Error)
		must.Equal(2, obj.(TestObj).Field)
	}))
	t.Run("signature based format", test.Case(func(ctx *countlog.Context) {
		type TestVersion1 struct {
			Field1 int
		}
		type TestVersion2 struct {
			Field1 uint
			Field2 uint
		}
		obj := TestVersion1{1}
		api := pickle.Config{UseSignature: true}.Froze()
		encoded := must.Call(api.Marshal, obj)[0].([]byte)
		must.Equal([]byte{
			0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		}, encoded[12:])
		decoded := must.Call(pickle.UnmarshalCandidates, encoded,
			(*TestVersion2)(nil), (*TestVersion1)(nil))[0]
		should.Equal(obj, *decoded.(*TestVersion1))
	}))
}
