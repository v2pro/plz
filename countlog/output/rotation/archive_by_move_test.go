package rotation

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"os"
	"io/ioutil"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/clock"
	"time"
)

func Test_ArchiveByMove(t *testing.T) {
	t.Run("", test.Case(func(ctx *countlog.Context) {
		os.Remove("/tmp/test.log")
		ioutil.WriteFile("/tmp/test.log", []byte("hello"), 0644)
		archiver := &ArchiveByMove{
			NamingStrategy: &NameByTime{
				Directory: "/tmp",
				Pattern: "test-2006.log",
			},
		}
		defer clock.ResetNow()
		clock.Now = func() time.Time {
			return time.Unix(0, 0)
		}
		must.Call(archiver.Archive, "/tmp/test.log")
		content := must.Call(ioutil.ReadFile, "/tmp/test-1970.log")[0]
		must.Equal([]byte("hello"), content)
	}))
}
