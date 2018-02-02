package test

import (
	"testing"
	"time"
	"github.com/v2pro/plz/nfmt"
	"github.com/stretchr/testify/require"
	"fmt"
)

func Test_string_to_string(t *testing.T) {
	should := require.New(t)
	should.Equal("ahellob", fmt.Sprintf("a%sb", "hello"))
	should.Equal("ahellob", nfmt.Sprintf("a%(key)sb", "key", "hello"))
}

func Test_int_to_string(t *testing.T) {
	should := require.New(t)
	should.Equal("%!s(int=100)", fmt.Sprintf("%s", 100))
	should.Equal("100", nfmt.Sprintf("%(key)s", "key", 100))
}

func Test_bytes_to_string(t *testing.T) {
	should := require.New(t)
	should.Equal("hello", fmt.Sprintf("%s", []byte("hello")))
	should.Equal("hello", nfmt.Sprintf("%(key)s", "key", []byte("hello")))

}

func Test_printf(t *testing.T) {
	nfmt.Printf("%(key)s", "key", "hello")
}

func Benchmark_string_to_string(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		nfmt.Sprintf("%(key)s", "key", "hello")
		//fmt.Sprintf("%s", "hello")
	}
}

func Benchmark_time_now(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		time.Now().String()
	}
}
