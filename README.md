Hal
===

[![GoDoc](https://godoc.org/github.com/nvellon/hal?status.svg)](https://godoc.org/github.com/nvellon/hal)

Go implementation of [HAL standard](http://stateless.co/hal_specification.html).

This is a work in progress... Everything might/will change.

Usage
-----

Gohal gives a way to translate structs/objects/entities/resources into HAL format, which can be easily translated into Json or Xml.

It provides the interface Mapper which, when implemented by a struct, returns Json/Xml Marshaler-friendly structure:

```go
type Mapper interface {
	GetMap() Entry
}
```

For a given Product struct:

```go
type Product struct {
	Code int
	Name string
	Price int
}
```

Implementint the HalEncoder interface:

```go
func (p Product) GetMap() hal.Entry {
	return hal.Entry{
		"name":  p.Name,
		"price": p.Price,
	}
}
```

This way you define which fields you want translated and which ones not (notice "Code" is not there).

Then you can just create a HAL Resource for a Product object by:

```go
p := Product{1, "Some Product", 10}

pr := hal.NewResource(p, "http://rest.api/products/some_product")
```

And when you need the Json encoded, you can do json.Marsal:

```go
j, err := json.Marshal(&pr)
if err != nil {
	fmt.Printf("%s", err)
}

fmt.Printf("%s", j)
```

Output:
```json
{
	"_links": {
		"self": {"href": "http://rest.api/products/some_product"}
	},
	"name": "Some product",
	"price": 10
}
```