package counselor

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
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
	ioutil.WriteFile("/tmp/conf/testNS/testObj/toggle",
		append([]byte("data_format:toggle\n"), toggleItem...), 0666)
	ioutil.WriteFile("/tmp/conf/testNS/testObj/true",
		[]byte("data_format:json\ntrue"), 0666)
	ioutil.WriteFile("/tmp/conf/testNS/testObj/false",
		[]byte("data_format:json\nfalse"), 0666)
	SourceDir = "/tmp/conf"
	pid1 := "10248"
	should.True(ShouldUse("testNS", "testObj", "divide_buckets_by", pid1))
	pid2 := "9876"
	should.False(ShouldUse("testNS", "testObj", "divide_buckets_by", pid2))
}
