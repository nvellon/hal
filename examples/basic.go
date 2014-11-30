package main

import (
	"encoding/json"
	"fmt"
	"github.com/nvellon/gohal"
)

type (
	Category struct {
		Code int
		Name string
	}

	Product struct {
		Code     int
		Name     string
		Price    int
		Category Category
	}
)

func (p Product) Encode() map[string]interface{} {
	return map[string]interface{}{
		"name":  p.Name,
		"price": p.Price,
	}
}

func (c Category) Encode() map[string]interface{} {
	return map[string]interface{}{
		"name": c.Name,
	}
}

func main() {
	c := Category{1, "Some Category"}
	p := Product{1, "Some Product", 10, c}

	// Creating HAL Resources
	pr := gohal.NewResource(p, "http://some_host/products/some_product")
	cr := gohal.NewResource(p.Category, "http://some_host/categories/some_category")

	// Embeding category into product
	pr.Embed(cr)

	// JSON Encoding
	j, err := json.Marshal(&pr)
	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", j)
}
