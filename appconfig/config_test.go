package appconfig

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	var json = `{"s3_bucket": "foo"}`

	readFile = func(string) ([]byte, error) {
		return []byte(json), nil
	}

	config, err := LoadConfig()

	assert.Equal(t, config.S3Bucket, "foo")

	assert.Nil(t, err)
}

func TestLoadConfigWithInvalidJson(t *testing.T) {
	var invalidJSON = `{ "configuration": { "s3_bucket": "foo" }`

	readFile = func(string) ([]byte, error) {
		return []byte(invalidJSON), nil
	}

	config, err := LoadConfig()

	expected := false

	switch err.(type) {
	case *json.SyntaxError:
		expected = true
	}

	assert.True(t, expected, "Expected error to be of type json.SyntaxError")

	assert.Equal(t, config.S3Bucket, "")
}
