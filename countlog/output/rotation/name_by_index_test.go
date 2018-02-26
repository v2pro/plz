package rotation

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"os"
	"github.com/v2pro/plz/test/must"
	"io/ioutil"
)

func Test_NameByIndex(t *testing.T) {
	t.Run("directory not exists", test.Case(func(ctx *countlog.Context) {
		os.RemoveAll("/tmp/testlog")
		naming := &NameByIndex{
			Directory:  "/tmp/testlog",
			Pattern:    "test-{i}.log",
			StartIndex: 1,
		}
		must.Equal([]string{}, must.Call(naming.ListFiles)[0])
	}))
	t.Run("directory with indexed and not indexed files", test.Case(func(ctx *countlog.Context) {
		os.RemoveAll("/tmp/testlog")
		os.MkdirAll("/tmp/testlog", 0755)
		ioutil.WriteFile("/tmp/testlog/test-2.log", []byte{}, 0644)
		ioutil.WriteFile("/tmp/testlog/test-a.log", []byte{}, 0644)
		naming := &NameByIndex{
			Directory:  "/tmp/testlog",
			Pattern:    "test-%d.log",
			StartIndex: 1,
		}
		must.Equal([]string{
			"/tmp/testlog/test-2.log",
		}, must.Call(naming.ListFiles)[0])
	}))
}
