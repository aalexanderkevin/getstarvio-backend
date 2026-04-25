package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Service  ServiceConfig
	DB       DBConfig
	JWT      JWTConfig
	Google   GoogleConfig
	Meta     MetaConfig
	Xendit   XenditConfig
	Worker   WorkerConfig
	LogLevel string
}

type ServiceConfig struct {
	Env       string
	Host      string
	Port      string
	APIPrefix string
}

type DBConfig struct {
	Host          string
	Port          string
	User          string
	Password      string
	Name          string
	SSLMode       string
	MigrationPath string
}

type JWTConfig struct {
	Secret           string
	AccessTTLMinutes int
	RefreshTTLHours  int
}

type GoogleConfig struct {
	ClientID          string
	AllowInsecureMock bool
}

type MetaConfig struct {
	APIVersion         string
	PhoneNumberID      string
	WABAID             string
	AccessToken        string
	WebhookVerifyToken string
	AppSecret          string
	HTTPTimeoutSeconds int
}

type XenditConfig struct {
	APIKey          string
	CallbackToken   string
	SuccessRedirect string
	FailureRedirect string
}

type WorkerConfig struct {
	PollIntervalSeconds int
}

var (
	once     sync.Once
	instance Config
)

func Instance() Config {
	once.Do(func() {
		instance = MustLoad()
	})
	return instance
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load() (Config, error) {
	if err := loadDotEnv(".env"); err != nil {
		return Config{}, err
	}

	cfg := Config{
		Service: ServiceConfig{
			Env:       getEnv("SERVICE_ENV", "development"),
			Host:      getEnv("SERVICE_HOST", "0.0.0.0"),
			Port:      getEnv("SERVICE_PORT", "8080"),
			APIPrefix: getEnv("SERVICE_API_PREFIX", "/v1"),
		},
		DB: DBConfig{
			Host:          getEnv("POSTGRES_HOST", "127.0.0.1"),
			Port:          getEnv("POSTGRES_PORT", "5432"),
			User:          getEnv("POSTGRES_USER", "postgres"),
			Password:      getEnv("POSTGRES_PASSWORD", "postgres"),
			Name:          getEnv("POSTGRES_DB", "getstarvio"),
			SSLMode:       getEnv("POSTGRES_SSLMODE", "disable"),
			MigrationPath: getEnv("POSTGRES_MIGRATION_PATH", "database/migrations"),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "change-me"),
			AccessTTLMinutes: getEnvAsInt("JWT_ACCESS_TTL_MINUTES", 60),
			RefreshTTLHours:  getEnvAsInt("JWT_REFRESH_TTL_HOURS", 720),
		},
		Google: GoogleConfig{
			ClientID:          getEnv("GOOGLE_CLIENT_ID", ""),
			AllowInsecureMock: getEnvAsBool("ALLOW_INSECURE_GOOGLE_MOCK", true),
		},
		Meta: MetaConfig{
			APIVersion:         getEnv("META_API_VERSION", "v22.0"),
			PhoneNumberID:      getEnv("META_PHONE_NUMBER_ID", ""),
			WABAID:             getEnv("META_WABA_ID", getEnv("META_WABAID", "")),
			AccessToken:        getEnv("META_ACCESS_TOKEN", ""),
			WebhookVerifyToken: getEnv("META_WEBHOOK_VERIFY_TOKEN", ""),
			AppSecret:          getEnv("META_APP_SECRET", ""),
			HTTPTimeoutSeconds: getEnvAsInt("META_HTTP_TIMEOUT_SECONDS", 30),
		},
		Xendit: XenditConfig{
			APIKey:          getEnv("XENDIT_API_KEY", ""),
			CallbackToken:   getEnv("XENDIT_CALLBACK_TOKEN", ""),
			SuccessRedirect: getEnv("XENDIT_SUCCESS_REDIRECT", ""),
			FailureRedirect: getEnv("XENDIT_FAILURE_REDIRECT", ""),
		},
		Worker:   WorkerConfig{PollIntervalSeconds: getEnvAsInt("WORKER_POLL_INTERVAL_SECONDS", 30)},
		LogLevel: strings.ToLower(getEnv("LOG_LEVEL", "info")),
	}

	if cfg.JWT.Secret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func loadDotEnv(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}

	lines := strings.Split(string(content), "\n")
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		sep := strings.Index(line, "=")
		if sep <= 0 {
			continue
		}

		key := strings.TrimSpace(line[:sep])
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		value := strings.TrimSpace(line[sep+1:])
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set env %s: %w", key, err)
		}
	}

	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return fallback
}
