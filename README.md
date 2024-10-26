# AnyConfig

With a single line of code, bring flexible configuration for to command line utilities or microservices.

Requires Go 1.22 or newer.

## Features

* ✅ __Configuration struct__: one struct to rule them all: define a single struct to hold all your configuration, its validation and documentation.
* ✅ __Environment variables__: automatically bind your configuration from environment variables using [Envconfig](https://github.com/sethvargo/go-envconfig).
* ✅ __File configuration__: read your configuration from a JSON or YAML file. Useful for command-line utilities that need a persistent configuration.
* ✅ __Validation__: Define simple validation rules for your configuration using [Validator](https://github.com/go-playground/validator).
* ✅ __Defaults__: provide clever defaults for your configuration

Notes: this is a very simple wrapper around [Envconfig](https://github.com/sethvargo/go-envconfig), and most of the features provide this library and json/yaml unmarshallers.

## Usage

### Simple example

```go
package main

import (
    "log"
    
    "sustainyfacts.dev/anyconfig"
)

func main() {
  var conf struct {
    Port int    `env:"PORT"`
    Host string `env:"USERNAME"`
  }

  if err := anyconfig.Read(&conf); err != nil {
    log.Fatal(err)
  }
}
```

### Complex example

```go
package main

import (
    "log"
    
    "sustainyfacts.dev/anyconfig"
)

// Configuration
type Config struct {
  Server struct {
    Port     int    `env:"PORT" validate:"gte=0"`
    Hostname string `json:"host" yaml:"host" env:"USERNAME" validate:"hostname"`
  } `env:", prefix=SERVER_"`
  Logging struct {
    Environment string `json:"env" yaml:"env" env:"ENV, overwrite" validate:"oneof=prod staging dev"`
    Level       string `env:"LEVEL" validate:"oneof=debug info warn error"`
  } `env:", prefix=LOGGING_"`
}

func main() {
  var conf Config
  if err := anyconfig.Read(&conf, anyconfig.WithFile("conf.yaml")); err != nil {
    log.Fatal(err)
  }
}
```

Together with the following YAML file:

```yaml
# This is am example yaml config file
server:
  host: example.com
  port: 8080
logging:
  env: dev
  level: debug
```

Or with the following Enviroment variables:

```bash
SERVER_HOST=example.com
SERVER_PORT=8080
LOGGING_ENV=dev
LOGGING_LEVEL=debug
```

For the complete reference for "env:" annotations, please refer to the project [Envconfig](https://github.com/sethvargo/go-envconfig).

## License

AnyCache is released under the Apache 2.0 license (see [LICENSE](LICENSE))