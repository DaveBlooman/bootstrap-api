package appconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Config Cosmos configuration data
type Config struct {
	RedisHost string `json:"redis_host"`
	S3Bucket  string `json:"s3_bucket"`
}

var readFile = ioutil.ReadFile

// LoadConfig loads AWS config
func LoadConfig() (Config, error) {
	var path string

	if os.Getenv("APP_ENV") == "production" {
		path = "/app/config.json"
	} else {
		_, filedir, _, _ := runtime.Caller(0)
		path = fmt.Sprintf("%s/../config/config.json", filepath.Dir(filedir))
	}

	file, error := readFile(path)
	if error != nil {
		message := map[string]interface{}{"event": "ConfigLoadError"}
		log.Fatal(message)
		os.Exit(1)
	}

	var config Config
	err := json.Unmarshal(file, &config)

	return config, err
}
