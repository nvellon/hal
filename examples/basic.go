package main

import (
	"encoding/json"
	"fmt"
	"github.com/nvellon/hal"
)

type (
	Product struct {
		Code  int
		Name  string
		Price int
	}
)

func (p Product) GetMap() hal.Entry {
	return hal.Entry{
		"name":  p.Name,
		"price": p.Price,
	}
}

func main() {
	p := Product{
		Code:  1,
		Name:  "Some Product",
		Price: 10,
	}

	// Creating HAL Resources
	pr := hal.NewResource(p, "/products/1")

	// Adding an extra link
	pr.AddNewLink("help", "/docs")

	// JSON Encoding
	j, err := json.MarshalIndent(pr, "", "  ")
	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", j)
}
