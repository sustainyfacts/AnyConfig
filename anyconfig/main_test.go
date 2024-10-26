package anyconfig

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
		Port     int    `env:"PORT" validate:"gt=0"`
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
	err := Read(&conf, WithFile(fileName)) // Read the configuration

	if assert.NoError(t, err) {
		assert.Equal(t, MyConfig{Port: 8080, Username: "sustainyfacts"}, conf)
	}
}

func Test_ReadFromFileInHomeDir(t *testing.T) {
	type MyConfig struct {
		Port     int    `json:"port"`
		Username string `json:"username"`
	}
	const fileName = ".mysecretconfig.json"
	home, _ := os.UserHomeDir()
	filePath := home + string(os.PathSeparator) + ".mysecretconfig.json"
	if file, err := os.Create(filePath); assert.NoError(t, err) {
		file.WriteString(`{"port":8080,"username":"sustainyfacts"}`)
		defer func() { os.Remove(filePath) }() // Cleanup after test
		file.Close()
	}

	var conf MyConfig
	err := Read(&conf, WithFile(fileName)) // Read the configuration

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
	err := Read(&conf, WithFile(fileName)) // Read the configuration

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
	err := Read(&conf, WithFile(fileName)) // Read the configuration

	if assert.NoError(t, err) {
		assert.Equal(t, MyConfig{Port: 8080, Username: "sustainyfacts"}, conf)
	}
}

func Test_ReadFromFileAndEnv(t *testing.T) {
	type MyConfig struct {
		Port     int    `json:"port" env:"PORT"`
		Username string `json:"username" env:"USERNAME, overwrite"`
	}
	const fileName = "defaults.json"
	if file, err := os.Create(fileName); assert.NoError(t, err) {
		file.WriteString(`{"port":8080,"username":"default_user"}`)
		defer func() { os.Remove(fileName) }() // Cleanup after test
		file.Close()
	}

	// Override USERNAME
	if err := os.Setenv("USERNAME", "sustainyfacts"); err != nil {
		assert.NoError(t, err)
	}

	var conf MyConfig
	err := Read(&conf, WithFile(fileName)) // Read the configuration

	if assert.NoError(t, err) {
		assert.Equal(t, MyConfig{Port: 8080, Username: "sustainyfacts"}, conf)
	}
}

func Test_ReadFile(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		testName    string
		fileName    string
		fullPath    string
		expectError bool
	}{
		{"from home dir", "hello.txt", home + "/hello.txt", false},
		{"from home dir with ~", "~/hello.txt", home + "/hello.txt", false},
		{"file not exist in home dir", "~/hello.txt", "hello.txt", true},
		{"relative dir", "./hello.txt", "hello.txt", false},
		{"relative dir to home", ".secret.txt", home + "/.secret.txt", false},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			if file, err := os.Create(tc.fullPath); assert.NoError(t, err) {
				file.WriteString(`Hello World!`)
				defer func() { os.Remove(tc.fullPath) }() // Cleanup after test
				file.Close()
			}

			bytes, err := readFile(tc.fileName)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.EqualValues(t, "Hello World!", bytes)
				assert.NoError(t, err)
			}
		})
	}

}

func Test_FullExample(t *testing.T) {
	type Server struct {
		Port     int    `env:"PORT" validate:"gte=0"`
		Hostname string `json:"host" yaml:"host" env:"USERNAME" validate:"hostname"`
	}
	type Logging struct {
		Environment string `json:"env" yaml:"env" env:"ENV, overwrite" validate:"oneof=prod staging dev"` // ENV variable with override file
		Level       string `env:"LEVEL" validate:"oneof=debug info warn error"`
	}
	type MyConfig struct {
		Server  Server  `env:", prefix=SERVER_"`
		Logging Logging `env:", prefix=LOGGING_"`
	}
	const fileName = "config.yaml"
	if file, err := os.Create(fileName); assert.NoError(t, err) {
		file.WriteString(`# This is am example yaml config file
server:
  host: example.com
logging:
  env: dev
  level: debug`)
		defer func() { os.Remove(fileName) }() // Cleanup after test
		file.Close()
	}

	os.Setenv("LOGGING_ENV", "prod")
	os.Setenv("SERVER_PORT", "8080")

	var conf MyConfig
	err := Read(&conf, WithFile(fileName)) // Read the configuration

	if assert.NoError(t, err) {
		assert.Equal(t, MyConfig{Server{Port: 8080, Hostname: "example.com"}, Logging{Environment: "prod", Level: "debug"}}, conf)
	}
}
