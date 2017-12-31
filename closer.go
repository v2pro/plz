package plz

import (
	"io"
	"runtime"
	"fmt"
	"github.com/v2pro/plz/countlog"
)

var Recover = countlog.Recover
type MultiError []error

func (errs MultiError) Error() string {
	return "multiple errors"
}

func NewMultiError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return MultiError(errs)
}

func Close(resource io.Closer, properties ...interface{}) error {
	err := resource.Close()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		closedAt := fmt.Sprintf("%s:%d", file, line)
		properties = append(properties, "err", err)
		countlog.Error("event!close "+closedAt, properties...)
		return err
	}
	return nil
}

func CloseAll(resources []io.Closer, properties ...interface{}) error {
	var errs []error
	for _, resource := range resources {
		err := resource.Close()
		if err != nil {
			_, file, line, _ := runtime.Caller(1)
			closedAt := fmt.Sprintf("%s:%d", file, line)
			properties = append(properties, "err", err)
			countlog.Error("event!close "+closedAt, properties...)
			errs = append(errs, err)
		}
	}
	return NewMultiError(errs)
}

type funcResource struct {
	f func() error
}

func (res funcResource) Close() error {
	return res.f()
}

func WrapCloser(f func() error) io.Closer {
	return &funcResource{f}
}
