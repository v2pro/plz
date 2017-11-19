package witch

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/countlog"
	"bytes"
)

type fixedStateExporter struct {
	state map[string]interface{}
}

func (se *fixedStateExporter) ExportState() map[string]interface{} {
	return se.state
}

func Test_recursive_state(t *testing.T) {
	se1 := &fixedStateExporter{
		state: map[string]interface{}{
			"hello": "world",
		},
	}
	se1.state["myself"] = se1
	should := require.New(t)
	buf := bytes.NewBuffer(nil)
	marshalState(map[string]countlog.StateExporter{
		"se1": se1,
	}, buf)
	should.Contains(buf.String(), "myself")
}
