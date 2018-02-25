package rotation

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"os"
	"github.com/v2pro/plz/test/must"
	"io/ioutil"
)

func Test_TriggerBySize(t *testing.T) {
	t.Run("init size below threshold", test.Case(func(ctx *countlog.Context) {
		os.Remove("/tmp/test.log")
		file := must.Call(os.OpenFile, "/tmp/test.log",
			os.O_CREATE|os.O_RDWR, os.FileMode(0644))[0].(*os.File)
		trigger := &TriggerBySize{SizeInKB: 1}
		result := must.Call(trigger.UpdateStat, nil, file, make([]byte, 24))
		must.Equal(int64(0), result[0])
		must.Equal(false, result[1])
		result = must.Call(trigger.UpdateStat, result[0], file, make([]byte, 24))
		must.Equal(int64(24), result[0])
		must.Equal(false, result[1])
		result = must.Call(trigger.UpdateStat, result[0], file, make([]byte, 1000))
		must.Equal(int64(1024), result[0])
		must.Equal(true, result[1])
	}))
	t.Run("init size equal threshold", test.Case(func(ctx *countlog.Context) {
		os.Remove("/tmp/test.log")
		ioutil.WriteFile("/tmp/test.log", make([]byte, 1024), 0644)
		file := must.Call(os.OpenFile, "/tmp/test.log",
			os.O_CREATE|os.O_RDWR, os.FileMode(0644))[0].(*os.File)
		trigger := &TriggerBySize{SizeInKB: 1}
		result := must.Call(trigger.UpdateStat, nil, file, make([]byte, 24))
		must.Equal(int64(1024), result[0])
		must.Equal(true, result[1])
	}))
}
