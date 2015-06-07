package main

import (
	"encoding/json"
	"fmt"
	"github.com/nvellon/hal"
)

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

func main() {
	// Creating HAL Resources
	r := hal.NewResource(Response{Count: 10, Total: 20}, "/tasks")
	r.AddNewLink("next", "/tasks=page=2")

	t1 := hal.NewResource(Task{Id: 1, Name: "Some Task"}, "/tasks/1")
	t2 := hal.NewResource(Task{Id: 2, Name: "Some Task"}, "/tasks/2")
	t3 := hal.NewResource(Task{Id: 3, Name: "Some Task"}, "/tasks/3")

	// Embeding tasks
	r.Embed("tasks", t1)
	r.Embed("tasks", t2)
	r.Embed("tasks", t3)

	// JSON Encoding
	j, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", j)
}
