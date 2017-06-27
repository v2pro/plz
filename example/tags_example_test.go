package example

import (
	"encoding/json"
	"fmt"
	"github.com/v2pro/plz/tags"
	"reflect"
)

type Order struct {
	OrderId   int `json:"order_id"`
	ProductId int `json:"product_id"`
}

func (order *Order) DefineTags() tags.Tags {
	return tags.D(
		tags.D("comment", "some more info about the struct itself"),
		tags.D(&order.OrderId, "validation", "required", "tag_is_not_only_string", 100),
		tags.D(&order.ProductId, "validation", "required"),
	)
}

func Example_tags_defined_by_struct() {
	structTags := tags.Get(reflect.TypeOf(Order{}))
	fieldTags, _ := json.Marshal(structTags.Fields["OrderId"])
	fmt.Println(string(fieldTags))
	// Output: {"json":"order_id","tag_is_not_only_string":100,"validation":"required"}
}

type Product struct {
	ProductId int `json:"product_id"`
}

func Example_tags_defined_externally() {
	tags.Define(func(p *Product) tags.Tags {
		return tags.D(
			tags.D("comment", "some more info about the struct itself"),
			tags.D(&p.ProductId, "validation", "required", "tag_is_not_only_string", 100),
		)
	})
	structTags := tags.Get(reflect.TypeOf(Product{}))
	fieldTags, _ := json.Marshal(structTags.Fields["ProductId"])
	fmt.Println(structTags.Struct["comment"])
	fmt.Println(string(fieldTags))
	// Output:
	// some more info about the struct itself
	// {"json":"product_id","tag_is_not_only_string":100,"validation":"required"}
}
