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

// Options for reading the configuration
type Option func(*config)

// Read the configuration from a file. Usage examples:
//
//	WithFile("~/.config.yaml") // reference to home directory of current user
//	WithFile("/home/user/.config.yaml") // fully qualified path
//	WithFile(".config.yaml") // will first look for ~/.conf.yaml and then .conf.yaml in the current directory
func WithFile(file string) Option {
	return func(c *config) {
		c.file = file
	}
}

// Defines how the configuration should be read
type config struct {
	file string
}

// Validator instance
var validate = validator.New()

// Read the configuration using the specified options into a struct given as pointer
// S should be a struct with annotated fields
func Read[S any](c *S, options ...Option) error {
	ctx := context.Background()

	config := new(config)
	// Apply all the functional options to configure the client.
	for _, opt := range options {
		opt(config)
	}

	// Read from file
	if config.file != "" {
		if err := unmarshalFromFile(config.file, c); err != nil {
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

// Unmarshals the configuration from a file,
// as json or yaml depending on the file extension
func unmarshalFromFile[S any](file string, c *S) (err error) {
	var bytes []byte
	if bytes, err = readFile(file); err != nil {
		return err
	}

	if strings.HasSuffix(file, ".json") {
		if err = json.Unmarshal(bytes, c); err != nil {
			return err
		}
	} else if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		if err := yaml.Unmarshal(bytes, c); err != nil {
			return err
		}
	}
	return nil
}

// Reads a file using absolute path, or relative to home or current directory
func readFile(file string) ([]byte, error) {
	home, _ := os.UserHomeDir()
	switch file[0] {
	case '/':
		return os.ReadFile(file)
	case '~':
		return os.ReadFile(home + "/" + file[1:])
	}

	// First try in the home directory
	if bytes, err := os.ReadFile(home + "/" + file); err == nil {
		return bytes, nil
	} else {
		// Fallback to current directory
		return os.ReadFile(file)
	}
}
