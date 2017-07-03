package example

import (
	"github.com/v2pro/plz/lang/app"
	"io/ioutil"
	"os"
	"github.com/v2pro/plz"
	"fmt"
)

func Example_app_finish_hook() {
	app.BeforeFinish = append(app.BeforeFinish, func(kv []interface{}) {
		ioutil.WriteFile("/tmp/hello", []byte("world"), os.ModeAppend|0666)
	})
	app.AfterPanic = append(app.AfterPanic, func(recovered interface{}, kv []interface{}) int {
		fmt.Println("panic", recovered)
		return 2
	})
	plz.RunApp(func() int {
		// os.Exit(0)
		return 0
	})
	// Output:
}
