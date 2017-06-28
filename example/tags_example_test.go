package example

import (
	"encoding/json"
	"fmt"
	"reflect"
	"github.com/v2pro/plz/tagging"
)

type Order struct {
	OrderId   int `json:"order_id"`
	ProductId int `json:"product_id"`
}

func (order *Order) DefineTags() tagging.Tags {
	return tagging.D(
		tagging.S("comment", "some more info about the struct itself"),
		tagging.F(&order.OrderId, "validation", "required", "tag_is_not_only_string", 100),
		tagging.F(&order.ProductId, "validation", "required"),
	)
}

func Example_tags_defined_by_struct() {
	structTags := tagging.Get(reflect.TypeOf(Order{}))
	fieldTags, _ := json.Marshal(structTags.Fields["OrderId"])
	fmt.Println(string(fieldTags))
	// Output: {"json":"order_id","tag_is_not_only_string":100,"validation":"required"}
}

type Product struct {
	ProductId int `json:"product_id"`
}

func Example_tags_defined_externally() {
	tagging.Define(func(p *Product) tagging.Tags {
		return tagging.D(
			tagging.S("comment", "some more info about the struct itself"),
			tagging.F(&p.ProductId, "validation", "required", "tag_is_not_only_string", 100),
		)
	})
	structTags := tagging.Get(reflect.TypeOf(Product{}))
	fieldTags, _ := json.Marshal(structTags.Fields["ProductId"])
	fmt.Println(structTags.Struct["comment"])
	fmt.Println(string(fieldTags))
	// Output:
	// some more info about the struct itself
	// {"json":"product_id","tag_is_not_only_string":100,"validation":"required"}
}
