package utils

import (
	"github.com/go-playground/validator/v10"
	"log"
	"os"
)

type EnvSetting struct {
	PostgresDsn               string `validate:"required"`
	Port                      string `validate:"required"`
	JwtSecret                 string `validate:"required"`
	RedisDsn                  string `validate:"required"`
	MailUsername              string `validate:"required"`
	MailPort                  string `validate:"required"`
	MailHost                  string `validate:"required"`
	MailPassword              string `validate:"required"`
	MailFrom                  string `validate:"required"`
	Environment               string `validate:"required"`
	CloudFlareBucket          string `validate:"required"`
	CloudFlareBucketUrl       string `validate:"required"`
	CloudFlareAccountID       string `validate:"required"`
	CloudFlareAccessKeyID     string `validate:"required"`
	CloudFlareAccessSecretKey string `validate:"required"`
}

var Settings EnvSetting

var validate *validator.Validate

func EnvInit() {
	Settings.PostgresDsn = os.Getenv("POSTGRES_DSN")
	Settings.Port = os.Getenv("PORT")
	Settings.JwtSecret = os.Getenv("JWT_SECRET")
	Settings.RedisDsn = os.Getenv("REDIS_DSN")
	Settings.MailUsername = os.Getenv("MAIL_USERNAME")
	Settings.MailPort = os.Getenv("MAIL_PORT")
	Settings.MailHost = os.Getenv("MAIL_HOST")
	Settings.MailPassword = os.Getenv("MAIL_PASSWORD")
	Settings.MailFrom = os.Getenv("MAIL_FROM")
	Settings.Environment = os.Getenv("ENVIRONMENT")
	//cloudflare r2 bucket

	Settings.CloudFlareBucket = os.Getenv("CLOUDFLARE_BUCKET")
	Settings.CloudFlareBucketUrl = os.Getenv("CLOUDFLARE_BUCKET_URL")
	Settings.CloudFlareAccessSecretKey = os.Getenv("CLOUDFLARE_ACCESS_SECRET_KEY")
	Settings.CloudFlareAccountID = os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	Settings.CloudFlareAccessKeyID = os.Getenv("CLOUDFLARE_ACCESS_KEY_ID")

	validate = validator.New()
	err := validate.Struct(Settings)

	if err != nil {
		validationError := err.(validator.ValidationErrors)
		log.Println("Error loading .env variables")
		log.Fatal(validationError.Error())
	}

}
