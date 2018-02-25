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

func Test_rotation(t *testing.T) {
	t.Run("write to new file", test.Case(func(ctx *countlog.Context) {
		os.RemoveAll("/tmp/testlog")
		os.Mkdir("/tmp/testlog", 0744)
		writer := must.Call(rotation.NewWriter, rotation.Config{
			WritePath: "/tmp/testlog/test.log",
		})[0].(*rotation.Writer)
		must.Call(writer.Write, []byte("hello"))
		must.Call(writer.Close)
		content := must.Call(ioutil.ReadFile, "/tmp/testlog/test.log")[0]
		must.Equal([]byte("hello"), content)
	}))
}
