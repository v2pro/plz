package rotation

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"os"
	"io/ioutil"
	"time"
	"github.com/v2pro/plz/clock"
)

func Test_NameByTime(t *testing.T) {
	t.Run("next file", test.Case(func(ctx *countlog.Context) {
		namer := &NameByTime{
			Directory: "/tmp",
			Pattern: "test-2006.log",
		}
		defer clock.ResetNow()
		clock.Now = func() time.Time {
			return time.Unix(0, 0)
		}
		newPath := must.Call(namer.NextFile)[0]
		must.Equal("/tmp/test-1970.log", newPath)
	}))
	t.Run("list files", test.Case(func(ctx *countlog.Context) {
		os.RemoveAll("/tmp/testlog")
		os.MkdirAll("/tmp/testlog", 0755)
		ioutil.WriteFile("/tmp/testlog/test-2016.log", []byte{}, 0644)
		ioutil.WriteFile("/tmp/testlog/test-2015.log", []byte{}, 0644)
		ioutil.WriteFile("/tmp/testlog/test-abc.log", []byte{}, 0644)
		namer := &NameByTime{
			Directory: "/tmp/testlog",
			Pattern: "test-2006.log",
		}
		files := must.Call(namer.ListFiles)[0]
		must.Equal([]string{
			"/tmp/testlog/test-2015.log",
			"/tmp/testlog/test-2016.log",
		}, files)
	}))
}