package rotation_test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/countlog/output/rotation"
	"github.com/v2pro/plz/test/must"
	"os"
	"io/ioutil"
	"time"
)

func resetTestLogDir() {
	os.RemoveAll("/tmp/testlog")
	os.Mkdir("/tmp/testlog", 0744)
}

func reset(cfg rotation.Config) *rotation.Writer {
	resetTestLogDir()
	writer := must.Call(rotation.NewWriter, cfg)[0].(*rotation.Writer)
	return writer
}

func Test_write(t *testing.T) {
	t.Run("write to existing file", test.Case(func(ctx *countlog.Context) {
		resetTestLogDir()
		ioutil.WriteFile("/tmp/testlog/test.log", []byte("hello\n"), 0644)
		writer := must.Call(rotation.NewWriter, rotation.Config{
			WritePath: "/tmp/testlog/test.log",
		})[0].(*rotation.Writer)
		defer must.Call(writer.Close)
		must.Call(writer.Write, []byte("world"))
		content := must.Call(ioutil.ReadFile, "/tmp/testlog/test.log")[0].([]byte)
		must.Equal("hello\nworld", string(content))
	}))
	t.Run("write to new file", test.Case(func(ctx *countlog.Context) {
		writer := reset(rotation.Config{
			WritePath: "/tmp/testlog/test.log",
		})
		defer must.Call(writer.Close)
		must.Call(writer.Write, []byte("hello"))
		content := must.Call(ioutil.ReadFile, "/tmp/testlog/test.log")[0].([]byte)
		must.Equal("hello", string(content))
	}))
	t.Run("write to new dir", test.Case(func(ctx *countlog.Context) {
		writer := reset(rotation.Config{
			WritePath: "/tmp/testlog/newdir/test.log",
		})
		defer must.Call(writer.Close)
		must.Call(writer.Write, []byte("hello"))
		content := must.Call(ioutil.ReadFile, "/tmp/testlog/newdir/test.log")[0].([]byte)
		must.Equal("hello", string(content))
	}))
}

func Test_rotation(t *testing.T) {
	t.Run("rotate by time interval", test.Case(func(ctx *countlog.Context) {
		writer := reset(rotation.Config{
			WritePath: "/tmp/testlog/test.log",
			TriggerStrategy: &rotation.TriggerByInterval{
				Interval: time.Second,
			},
			ArchiveStrategy: &rotation.ArchiveByMove{
				NamingStrategy: &rotation.NameByTime{
					Directory: "/tmp/testlog",
					Pattern:   "test-2006-01-02T15-04-05.log",
				},
			},
			RetainStrategy: &rotation.RetainByCount{3},
			PurgeStrategy: &rotation.PurgeByDelete{},
		})
		defer must.Call(writer.Close)
		for i := 0; i < 60; i++ {
			writer.Write([]byte("hello\n"))
			time.Sleep(time.Second)
		}
	}))
}
