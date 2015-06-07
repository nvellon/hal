package main

import (
	"encoding/json"
	"fmt"
	"github.com/nvellon/hal"
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

func (p Product) GetMap() hal.Entry {
	return hal.Entry{
		"name":  p.Name,
		"price": p.Price,
	}
}

func (c Category) GetMap() hal.Entry {
	return hal.Entry{
		"name": c.Name,
	}
}

func main() {
	c := Category{
		Code: 1,
		Name: "Some Category",
	}

	p := Product{
		Code:     1,
		Name:     "Some Product",
		Price:    10,
		Category: c,
	}

	// Creating HAL Resources
	pr := hal.NewResource(p, "/products/1")
	cr1 := hal.NewResource(p.Category, "/categories/1")
	cr2 := hal.NewResource(p.Category, "/categories/2")
	cr3 := hal.NewResource(p.Category, "/categories/3")

	// Embeding categories into product
	pr.Embed("categories", cr1)
	pr.Embedded.Del("categories")
	pr.Embedded.Set("categories", cr2)
	pr.Embedded.Add("categories", cr3)

	// JSON Encoding
	j, err := json.MarshalIndent(pr, "", "  ")
	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", j)
}
