package example

import (
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/app"
	"io/ioutil"
	"os"
)

func Example_app_finish_hook() {
	app.Spi.Append(app.Config{
		AfterFinish: func(kv []interface{}) {
			ioutil.WriteFile("/tmp/hello", []byte("world"), os.ModeAppend|0666)
		},
	})
	plz.RunApp(func() int {
		// os.Exit(0)
		return 0
	})
	// Output:
}
