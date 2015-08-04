Hal
===

[![Build Status](https://travis-ci.org/nvellon/hal.svg)](https://travis-ci.org/nvellon/hal)
[![Coverage Status](https://coveralls.io/repos/nvellon/hal/badge.svg?branch=master&service=github)](https://coveralls.io/github/nvellon/hal?branch=master)
[![GoDoc](https://godoc.org/github.com/nvellon/hal?status.svg)](https://godoc.org/github.com/nvellon/hal)

Go implementation of the [HAL standard](http://stateless.co/hal_specification.html).

This is a work in progress... Everything might/will change.

Usage
-----

The library gives a way of mapping Go Structs into HAL Resources by implementing the `hal.Mapper` interface. You only need to define which fields you want and how you want them translated.

```go
type Mapper interface {
	GetMap() Entry
}
```

For a given Product struct, this would be the `hal.Mapper` implementation:

```go
type Product struct {
	Code int
	Name string
	Price int
}

func (p Product) GetMap() hal.Entry {
	return hal.Entry{
		"name":  p.Name,
		"price": p.Price,
	}
}
```

Then you can just create a HAL Resource for a Product by:

```go
p := Product{
	Code: 1,
	Name: "Some Product",
	Price: 10
}

pr := hal.NewResource(p, "http://rest.api/products/1")
```

And pass it through `json.Marsal` when needed getting this result:

```json
{
	"_links": {
		"self": {"href": "http://rest.api/products/1"}
	},
	"name": "Some product",
	"price": 10
}
```

Embedded Resources
------------------

Let's say your API has to serve a list of Task structs.

Since for HAL standard everything is a resource, even the entire API response could be seen as a resource containing other embedded resources. Check this out:

```go
type (
	Response struct {
		Count int
		Total int
	}

	Task struct {
		Id   int
		Name string
	}
)

func (p Response) GetMap() hal.Entry {
	return hal.Entry{
		"count": p.Count,
		"total": p.Total,
	}
}

func (c Task) GetMap() hal.Entry {
	return hal.Entry{
		"id":   c.Id,
		"name": c.Name,
	}
}
```

Then you could create the Resources by doing something like this:

```go
// Creating Response resource
r := hal.NewResource(Response{Count: 10, Total: 20}, "/tasks")
r.AddNewLink("next", "/tasks=page=2")

// Creating Task resources
t1 := hal.NewResource(Task{Id: 1, Name: "Some Task"}, "/tasks/1")
t2 := hal.NewResource(Task{Id: 2, Name: "Some Task"}, "/tasks/2")

// Embedding
r.Embed("tasks", t1)
r.Embed("tasks", t2)
```

Output:

```json
{
  "_embedded": {
    "tasks": [
      {
        "_links": {
          "self": {
            "href": "/tasks/1"
          }
        },
        "id": 1,
        "name": "Some Task"
      },
      {
        "_links": {
          "self": {
            "href": "/tasks/2"
          }
        },
        "id": 2,
        "name": "Some Task"
      }
    ]
  },
  "_links": {
    "next": {
      "href": "/tasks=page=2"
    },
    "self": {
      "href": "/tasks"
    }
  },
  "count": 10,
  "total": 20
}
```

Todo
----

 * CURIEs support.
 * XML Marshaler.