package spi

import (
	"os"
	"github.com/v2pro/plz/nfmt"
)

var OnError = func(err error) {
	nfmt.Fprintf(os.Stderr, "countlog encountered error: %(err)s\n", "err", err.Error())
	os.Stderr.Sync()
}