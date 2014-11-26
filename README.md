Gohal
=====

Go implementation of (HAL)[http://stateless.co/hal_specification.html] standard.

This is a work in progress... Everything might/will change.

Examples
--------

Marshal:

```
self := Link{"self", "http://localhost/"}

embedded := Resource{Link{"self", "http://localhost/embedded"}}

resource := Resource{[]Link{self}, []Resource{embedded}}

j := json.Marshal(resource)

fmt.Sprintf(j)

```