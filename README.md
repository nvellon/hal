Gohal
=====

[![GoDoc](https://godoc.org/github.com/nvellon/gohal?status.svg)](https://godoc.org/github.com/nvellon/gohal)

Go implementation of [HAL standard](http://stateless.co/hal_specification.html).

This is a work in progress... Everything might/will change.

Usage
-----

Gohal gives a way to translate structs/objects/entities/resources into HAL-Json format.

It provides the interface HalEncoder which, when implemented by a struct, returns json.Marshal-friendly structure:

```go
type HalEncoder interface {
	Encode() map[string]interface{}
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
func (p Product) Encode() map[string]interface{} {
	return map[string]interface{}{
		"name":  p.Name,
		"price": p.Price,
	}
}
```

This way you define which fields you want translated and which ones not (notice "Code" is not there).

Then you can just create a HAL Resource for a Product object by:

```go
p := Product{1, "Some Product", 10}

pr := gohal.NewResource(p, "http://rest.api/products/some_product")
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