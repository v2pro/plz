package counselor

import (
	"path"
	"plugin"
)

func loadFn(fn string) (plugin.Symbol, error) {
	filename := path.Join(SourceDir, "counselor_plugins", "go", fn)
	so, err := plugin.Open(filename)
	if err != nil {
		return nil, err
	}
	return so.Lookup("fn")
}
