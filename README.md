# configuration

usage:
```go
package main

import (
    "fmt"

    "github.com/go-chocolate/configuration/configuration"
)

type Config struct {
    Name  string
    Value int
}

func main() {
    var c Config
    configuration.MustLoad(&c)
    fmt.Println(c)
}

```