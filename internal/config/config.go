package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Service ServiceConfig
	DB      DBConfig
	JWT     JWTConfig
	Google  GoogleConfig
	Meta    MetaConfig
	Xendit  XenditConfig
	Worker  WorkerConfig
	Internal InternalConfig
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
	Secret             string
	AccessTTLMinutes   int
	RefreshTTLHours    int
}

type GoogleConfig struct {
	ClientID             string
	AllowInsecureMock    bool
}

type MetaConfig struct {
	APIVersion    string
	PhoneNumberID string
	AccessToken   string
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

type InternalConfig struct {
	Token string
}

var (
	once sync.Once
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
			APIVersion:    getEnv("META_API_VERSION", "v20.0"),
			PhoneNumberID: getEnv("META_PHONE_NUMBER_ID", ""),
			AccessToken:   getEnv("META_ACCESS_TOKEN", ""),
		},
		Xendit: XenditConfig{
			APIKey:          getEnv("XENDIT_API_KEY", ""),
			CallbackToken:   getEnv("XENDIT_CALLBACK_TOKEN", ""),
			SuccessRedirect: getEnv("XENDIT_SUCCESS_REDIRECT", ""),
			FailureRedirect: getEnv("XENDIT_FAILURE_REDIRECT", ""),
		},
		Worker: WorkerConfig{PollIntervalSeconds: getEnvAsInt("WORKER_POLL_INTERVAL_SECONDS", 30)},
		Internal: InternalConfig{Token: getEnv("INTERNAL_API_TOKEN", "change-me-internal-token")},
		LogLevel: strings.ToLower(getEnv("LOG_LEVEL", "info")),
	}

	if cfg.JWT.Secret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
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
