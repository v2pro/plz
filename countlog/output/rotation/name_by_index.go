package rotation

import (
	"os"
	"path"
	"fmt"
)

type NameByIndex struct {
	Directory  string
	Pattern    string
	StartIndex int
}

func (naming *NameByIndex) ListFiles() ([]string, error) {
	dir, err := os.Open(naming.Directory)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	names, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(names))
	for _, name := range names {
		var idx int
		n, err := fmt.Sscanf(name, naming.Pattern, &idx)
		if err != nil {
			continue
		}
		if n != 1 {
			continue
		}
		paths = append(paths, path.Join(naming.Directory, name))
	}
	return paths, nil
}
