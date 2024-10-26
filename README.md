# AnyConfig

> 🚧 Work in Progress
> 
> This library is a working in progress, and as a result even the public API will probably change

With a single line of code, bring flexible configuration for to command line utilities or microservices.

Requires Go 1.22 or newer.

## Features

* ✅ __Configuration struct__: one struct to rule them all: define a single struct to hold all your configuration and its validation.
* ✅ __Environment variables__: automatically bind your configuration from environment variables using [Envconfig](https://github.com/sethvargo/go-envconfig).
* ✅ __File configuration__: read your configuration from a JSON or YAML file. Useful for command-line utilities that need a persistent configuration.
* ✅ __Validation__: Define simple validation rules for your configuration using [Validator](https://github.com/go-playground/validator).
* ✅ __Defaults__: provide clever defaults for your configuration
* 🚧 __Built-in documentation__: document your configuration in your code

## Usage

```go
package main

import (
    "log"
    
    "sustainyfacts.dev/anyconfig"
)

func main() {
  var conf struct {
    
  }

  if err := anyconfig.Read(&conf); err != nil {
    log.Fatal(err)
  }
}
```

## License

AnyCache is released under the Apache 2.0 license (see [LICENSE](LICENSE))