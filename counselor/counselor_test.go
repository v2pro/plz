package counselor

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
)

func Test_thousand_bucket(t *testing.T) {
	should := require.New(t)
	rawToggle := rawToggle{
		Variants: []*variant{
			{
				ItemName: "true",
				RuleFn:   "thousand_buckets",
				RuleArgs: map[string]interface{}{
					"segment_start": 0,
					"segment_end":   500,
				},
			},
		},
		DefaultVariant: "false",
	}
	toggleItem, err := json.Marshal(rawToggle)
	should.Nil(err)
	os.RemoveAll("/tmp/conf/testNS/testObj")
	os.MkdirAll("/tmp/conf/testNS/testObj", 0777)
	ioutil.WriteFile("/tmp/conf/testNS/testObj/toggle", toggleItem, 0666)
}
