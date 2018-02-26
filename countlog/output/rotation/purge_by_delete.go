package rotation

import (
	"os"
	"github.com/v2pro/plz"
)

type PurgeByDelete struct {
}

func (strategy *PurgeByDelete) Purge(purgeSet []Archive) error {
	var errs []error
	for _, archive := range purgeSet {
		err := os.Remove(archive.Path)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return plz.MergeErrors(errs...)
}
