package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// EnvConfig defines the expected environment variables for the application.
type EnvConfig struct {
	// App
	AppPort string `envconfig:"PORT"    default:"8080"`
	AppEnv  string `envconfig:"APP_ENV" default:"dev"`
	FrontendURL string `envconfig:"FRONTEND_URL" default:"http://localhost:3000"`

	// DB
	DBHost         string `envconfig:"DB_HOST"           required:"true"`
	DBPort         string `envconfig:"DB_PORT"           required:"true"`
	DBUser         string `envconfig:"DB_USER"           required:"true"`
	DBPassword     string `envconfig:"DB_PASSWORD"       required:"true"`
	DBName         string `envconfig:"DB_NAME"           required:"true"`
	DBSSLMode      string `envconfig:"DB_SSL_MODE"       required:"true"`
	DBMaxOpenConns int    `envconfig:"DB_MAX_OPEN_CONNS" default:"20"`
	DBMaxIdleConns int    `envconfig:"DB_MAX_IDLE_CONNS" default:"10"`
	
	// Mail
	MailSMTPHost     string `envconfig:"MAIL_SMTP_HOST"     default:"smtp.gmail.com"`
	MailSMTPPort     int    `envconfig:"MAIL_SMTP_PORT"     default:"587"`
	MailSMTPUser     string `envconfig:"MAIL_SMTP_USER"     required:"true"`
	MailSMTPPassword string `envconfig:"MAIL_SMTP_PASSWORD" required:"true"`
	MailEmailFrom    string `envconfig:"MAIL_EMAIL_FROM"    required:"true"`

	// OAuth
	GoogleRedirectURL  string `envconfig:"GOOGLE_REDIRECT_URL"`
	GoogleClientID     string `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET"`

	// S3
	AWSS3Region     string `envconfig:"AWS_S3_REGION"      default:"ap-southeast-1"`
	AWSS3AccessKey  string `envconfig:"AWS_S3_ACCESS_KEY"  required:"true"`
	AWSS3SecretKey  string `envconfig:"AWS_S3_SECRET_KEY"  required:"true"`
	AWSS3BucketName string `envconfig:"AWS_S3_BUCKET_NAME" required:"true"`



	// Redis
	RedisHost     string `envconfig:"REDIS_HOST"     default:"localhost"`
	RedisPort     string `envconfig:"REDIS_PORT"     default:"6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`

	// JWT
	JWTSecret         string        `envconfig:"JWT_SECRET"          default:"super-secret-key-12345"`
	JWTAccessDuration time.Duration `envconfig:"JWT_ACCESS_DURATION" default:"24h"`

	// RSA Keys
	RSAPrivateKeyPath string `envconfig:"RSA_PRIVATE_KEY_PATH" default:"config/key-pem/private_key.pem"`
	RSAPublicKeyPath  string `envconfig:"RSA_PUBLIC_KEY_PATH"  default:"config/key-pem/public-key.pem"`
}

var c *EnvConfig

// findProjectRoot walks up from CWD to find the directory containing go.mod
func findProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	cur := wd
	for {
		if _, err := os.Stat(filepath.Join(cur, "go.mod")); err == nil {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			// Reached filesystem root
			return wd
		}
		cur = parent
	}
}

// Init initializes config
func Init(env string) {
	if strings.TrimSpace(env) == "" {
		env = "dev"
	}

	root := findProjectRoot()

	// Load base .env (optional) from the project root
	_ = godotenv.Load(filepath.Join(root, ".env"))

	// Load env-specific file from the project root (optional)
	envFile := filepath.Join(root, fmt.Sprintf(".env.%s", env))
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Overload(envFile); err != nil {
			log.Printf("warning: could not load env file %s: %v", envFile, err)
		}
	} else if os.IsNotExist(err) {
		log.Printf("warning: env file %s not found; relying on environment variables", envFile)
	} else {
		log.Printf("warning: cannot stat env file %s: %v", envFile, err)
	}

	// Validate environment variables using envconfig
	var envCfg EnvConfig
	if err := envconfig.Process("", &envCfg); err != nil {
		log.Fatalf("[env] validation failed: %v", err)
	}

	c = &envCfg
}

// GetConfig returns config
func GetConfig() *EnvConfig {
	return c
}
