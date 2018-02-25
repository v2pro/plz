package rotation_test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/countlog/output/rotation"
	"github.com/v2pro/plz/test/must"
	"os"
	"io/ioutil"
)

func resetTestLogDir() {
	os.RemoveAll("/tmp/testlog")
	os.Mkdir("/tmp/testlog", 0744)
}

func Test_rotation(t *testing.T) {
	t.Run("write to existing file", test.Case(func(ctx *countlog.Context) {
		resetTestLogDir()
		ioutil.WriteFile("/tmp/testlog/test.log", []byte("hello\n"), 0644)
		writer := must.Call(rotation.NewWriter, rotation.Config{
			WritePath: "/tmp/testlog/test.log",
		})[0].(*rotation.Writer)
		must.Call(writer.Write, []byte("world"))
		must.Call(writer.Close)
		content := must.Call(ioutil.ReadFile, "/tmp/testlog/test.log")[0].([]byte)
		must.Equal("hello\nworld", string(content))
	}))
	t.Run("write to new file", test.Case(func(ctx *countlog.Context) {
		resetTestLogDir()
		writer := must.Call(rotation.NewWriter, rotation.Config{
			WritePath: "/tmp/testlog/test.log",
		})[0].(*rotation.Writer)
		must.Call(writer.Write, []byte("hello"))
		must.Call(writer.Close)
		content := must.Call(ioutil.ReadFile, "/tmp/testlog/test.log")[0].([]byte)
		must.Equal("hello", string(content))
	}))
}
