package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog/pickle"
	"github.com/json-iterator/go"
	"encoding/json"
)

func Test_string_slice(t *testing.T) {
	should := require.New(t)
	encoded, err := pickle.Marshal([]string{"h", "i"})
	should.Nil(err)
	should.Equal([]byte{
		0x18, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, // sliceHeader
		0x20, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,                         // string header
		0x11, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,                         // string header
		'h', 'i'}, encoded[8:])
	decoded, err := pickle.ReadonlyConfig.Unmarshal(encoded, (*[]string)(nil))
	should.Nil(err)
	should.Equal([]string{"h", "i"}, *decoded.(*[]string))
	decoded, err = pickle.Unmarshal(encoded, (*[]string)(nil))
	should.Nil(err)
	should.Equal([]string{"h", "i"}, *decoded.(*[]string))
}

func Test_ptr_slice(t *testing.T) {
	type TestObject struct {
		Field1 int
		Field2 int
	}
	should := require.New(t)
	encoded, err := pickle.Marshal([]*TestObject{{1, 2}, {3, 4}})
	should.Nil(err)
	should.Equal([]byte{
		24, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0,
		16, 0, 0, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, // [0]
		3, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, // [1]
	}, encoded[8:])
	decoded, err := pickle.ReadonlyConfig.Unmarshal(encoded, (*[]*TestObject)(nil))
	should.Nil(err)
	should.Equal([]*TestObject{{1, 2}, {3, 4}}, *decoded.(*[]*TestObject))
	decoded, err = pickle.Unmarshal(encoded, (*[]*TestObject)(nil))
	should.Nil(err)
	should.Equal([]*TestObject{{1, 2}, {3, 4}}, *decoded.(*[]*TestObject))
}

func Benchmark_string_slice(b *testing.B) {
	data := []string{"hello", "world"}
	gocEncoded, err := pickle.Marshal(data)
	if err != nil {
		b.Error(err)
	}
	jsonEncoded, _ := jsoniter.Marshal(data)
	b.Run("goc encode", func(b *testing.B) {
		b.ReportAllocs()
		stream := pickle.DefaultConfig.NewStream(nil)
		for i := 0; i < b.N; i++ {
			stream.Reset(stream.Buffer()[:0])
			stream.Marshal(data)
		}
	})
	b.Run("goc decode", func(b *testing.B) {
		b.ReportAllocs()
		iter := pickle.DefaultConfig.NewIterator(gocEncoded)
		for i := 0; i < b.N; i++ {
			iter.Reset(append(([]byte)(nil), gocEncoded...))
			iter.Unmarshal(&data)
		}
	})
	b.Run("json encode", func(b *testing.B) {
		b.ReportAllocs()
		jsonEncoder := jsoniter.ConfigFastest.BorrowStream(nil)
		for i := 0; i < b.N; i++ {
			jsonEncoder.Reset(nil)
			jsonEncoder.WriteVal(data)
		}
	})
	b.Run("json decode", func(b *testing.B) {
		b.ReportAllocs()
		jsonDecoder := jsoniter.ConfigFastest.BorrowIterator(nil)
		for i := 0; i < b.N; i++ {
			jsonDecoder.ResetBytes(jsonEncoded)
			jsonDecoder.ReadVal(&data)
		}
	})
	b.Run("encoding/json decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			json.Unmarshal(jsonEncoded, &data)
		}
	})
}
