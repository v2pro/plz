package tags

import (
	"github.com/json-iterator/go/require"
	"reflect"
	"testing"
)

func Test_define_externally(t *testing.T) {
	type TestObject struct {
		a int
	}
	Define(func(obj *TestObject) Tags {
		return D(
			S("k", "v"),
			F(&obj.a, "k2", "v2"),
		)
	})
	should := require.New(t)
	structTags := Get(reflect.TypeOf(TestObject{}))
	should.Equal(1, len(structTags.Fields))
	should.Equal(FieldTags{"k2": "v2"}, structTags.Fields["a"])
	should.Equal(map[string]interface{}{"k": "v"}, structTags.Struct)
}

type TestObject1 struct {
	a int
	b int
}

func (obj *TestObject1) DefineTags() Tags {
	return D(
		S("k", "v"),
		F(&obj.a, "k2", "v2"),
		F(&obj.b, "k3", "v3"),
	)
}

func Test_define_internally(t *testing.T) {
	should := require.New(t)
	structTags := Get(reflect.TypeOf(TestObject1{}))
	should.Equal(2, len(structTags.Fields))
	should.Equal(FieldTags{"k2": "v2"}, structTags.Fields["a"])
	should.Equal(FieldTags{"k3": "v3"}, structTags.Fields["b"])
	should.Equal(map[string]interface{}{"k": "v"}, structTags.Struct)
}

func Test_merge_with_string_defined_tags(t *testing.T) {
	type TestObject struct {
		Hello int `json:"hello" jsoniter:"abc"`
	}
	Define(func(obj *TestObject) Tags {
		return D(
			S("k", "v"),
			F(&obj.Hello, "k2", "v2"),
		)
	})
	should := require.New(t)
	structTags := Get(reflect.TypeOf(TestObject{}))
	should.Equal(1, len(structTags.Fields))
	should.Equal(FieldTags{"k2": "v2", "json": "hello", "jsoniter": "abc"}, structTags.Fields["Hello"])
	should.Equal(map[string]interface{}{"k": "v"}, structTags.Struct)
}
