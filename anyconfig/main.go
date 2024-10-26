package anyconfig

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
)

/*
Copyright © 2023 The Authors (See AUTHORS file)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Defines how the configuration should be read
type Config struct {
	// File to read configuration from.
	//
	// If a relative path or name is specified, it will look for a file in
	// the user home directory first, and if not found, in the current directory
	File string
	// Specify a prefix to
	EnvironmentVariablePrefix string
}

var DefaultConfig = Config{}

// Validator instance
var validate = validator.New()

// Read the configuration using the default configuration into a struct given as pointer
// S should be a struct with annotated fields
func Read[S any](c *S) error {
	return ReadWithConfig(c, DefaultConfig)
}

// Read the configuration using the default configuration
func ReadWithConfig[S any](c *S, config Config) error {
	ctx := context.Background()

	// Read from file
	if config.File != "" {
		if err := readFromFile(config.File, c); err != nil {
			return err
		}
	}

	// Read environment variables
	if err := envconfig.Process(ctx, c); err != nil {
		return err
	}

	// Validate
	if err := validate.Struct(*c); err != nil {
		return err
	}

	return nil // All good
}

func readFromFile[S any](file string, c *S) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if strings.HasSuffix(file, ".json") {
		if err := json.Unmarshal(bytes, c); err != nil {
			return err
		}
	} else if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		if err := yaml.Unmarshal(bytes, c); err != nil {
			return err
		}
	}
	return nil
}
