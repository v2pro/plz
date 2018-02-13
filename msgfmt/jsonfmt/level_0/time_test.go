package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
	"time"
	"github.com/v2pro/plz/reflect2"
)

func Test_time(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect2.TypeOf(time.Time{}))
	should.Equal(`"0001-01-01T00:00:00Z"`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf(time.Time{}))))

}