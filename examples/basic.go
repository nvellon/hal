package main

import (
	"encoding/json"
	"fmt"
	"github.com/nvellon/gohal"
)

type Product struct {
	gohal.Resource
	Name string `json:"name"`
}

type Category struct {
	gohal.Resource
	Name string `json:"name"`
}

func (p *Product) ToResource(b []byte) error {
	return nil
}
func (p *Product) ToMap() map[string]interface{} {
	return nil
}

func (c *Category) ToResource(b []byte) error {
	return nil
}
func (c *Category) ToMap() map[string]interface{} {
	return nil
}

func main() {

	p := Product{}
	p.Name = "some product"

	lp := gohal.NewLink("self", "http://localhost/products/some_product")
	p.AddLink(&lp)

	c := Category{}
	c.Name = "some category"
	lc := gohal.NewLink("self", "http://localhost/categories/some_category")
	c.AddLink(&lc)

	p.Embed(&c)

	j, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", j)
}
