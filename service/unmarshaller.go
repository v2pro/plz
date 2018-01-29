package service

import (
	"github.com/v2pro/plz/countlog"
)

type Unmarshaller interface {
	Unmarshal(ctx *countlog.Context, obj interface{}, input interface{}) error
}
