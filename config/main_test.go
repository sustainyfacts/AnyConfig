package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadFromEnv(t *testing.T) {
	type MyConfig struct {
		Port     int    `env:"PORT"`
		Username string `env:"USERNAME"`
	}

	os.Setenv("PORT", "8080")
	os.Setenv("USERNAME", "sustainyfacts")

	var conf MyConfig
	err := Read(&conf) // Read the configuration

	if assert.NoError(t, err) {
		assert.Equal(t, MyConfig{Port: 8080, Username: "sustainyfacts"}, conf)
	}
}

func Test_ValidateConfig(t *testing.T) {
	type MyConfig struct {
		Port     int    `env:"PORT" validate:"hostname_port"`
		Username string `env:"USERNAME" validate:"required,alphanum"`
	}

	os.Setenv("PORT", "-1")
	os.Setenv("USERNAME", "sustainyfacts")

	var conf MyConfig
	err := Read(&conf) // Read the configuration

	assert.Error(t, err, "Config should not be valid")
}

func Test_ReadFromJsonFile(t *testing.T) {
	type MyConfig struct {
		Port     int    `json:"port"`
		Username string `json:"username"`
	}
	const fileName = "myconfig.json"
	if file, err := os.Create(fileName); assert.NoError(t, err) {
		file.WriteString(`{"port":8080,"username":"sustainyfacts"}`)
		defer func() { os.Remove(fileName) }() // Cleanup after test
		file.Close()
	}

	var conf MyConfig
	err := ReadWithConfig(&conf, Config{File: fileName}) // Read the configuration

	if assert.NoError(t, err) {
		assert.Equal(t, MyConfig{Port: 8080, Username: "sustainyfacts"}, conf)
	}
}

func Test_ReadFromInvalidJsonFile(t *testing.T) {
	type MyConfig struct {
		Port     int    `json:"port"`
		Username string `json:"username"`
	}
	const fileName = "myconfig.json"
	if file, err := os.Create(fileName); assert.NoError(t, err) {
		file.WriteString(`{"port":8080,"username":sustainyfacts"}`)
		defer func() { os.Remove(fileName) }() // Cleanup after test
		file.Close()
	}

	var conf MyConfig
	err := ReadWithConfig(&conf, Config{File: fileName}) // Read the configuration

	assert.Error(t, err, "json file is invalid")
}

func Test_ReadFromYamlFile(t *testing.T) {
	type MyConfig struct {
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
	}
	const fileName = "myconfig.yaml"
	if file, err := os.Create(fileName); assert.NoError(t, err) {
		file.WriteString(`# Example YAML configuration
port: 8080 # Comment
username: sustainyfacts
`)
		defer func() { os.Remove(fileName) }() // Cleanup after test
		file.Close()
	}

	var conf MyConfig
	err := ReadWithConfig(&conf, Config{File: fileName}) // Read the configuration

	if assert.NoError(t, err) {
		assert.Equal(t, MyConfig{Port: 8080, Username: "sustainyfacts"}, conf)
	}
}
