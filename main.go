package main

import (
	"net/http"
	"os"

	"github.com/DaveBlooman/api-common/logger"
	"github.com/DaveBlooman/bootstrap-api/appconfig"
	"github.com/DaveBlooman/bootstrap-api/routes"
	"github.com/joho/godotenv"
)

func init() {
	var envFile string

	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "development")
	}

	if isProduction() {
		envFile = "/app/config/production.env"

		err := setupProduction()
		if err != nil {
			message := map[string]interface{}{"event": "ErrorLoadingConfig"}
			logger.Fatal(message)
		}
	}

	if envFile == "" {
		envFile = "config/" + os.Getenv("APP_ENV") + ".env"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		message := map[string]interface{}{"event": "ErrorLoadingEnvFile"}
		logger.Fatal(message)
	}
}

func main() {
	message := map[string]interface{}{"event": "ApiStarted"}
	logger.Info(message)
	router := routes.APIRouter()

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		message = map[string]interface{}{"message": err.Error()}
		logger.Fatal(message)
	}
}

func isProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

func setupProduction() error {
	config, err := appconfig.LoadConfig()
	if err != nil {
		return err
	}

	os.Setenv("REDIS_HOST", config.RedisHost)
	os.Setenv("S3_BUCKET", config.S3Bucket)

	return nil
}
