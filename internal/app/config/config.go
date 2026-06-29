// pkg/app/config/config.go
package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type HTTPConfig struct {
	ListenAddr        string        // ":8888" or "0.0.0.0:8888"
	ShutdownGrace     time.Duration // 10s
	ReadHeaderTimeout time.Duration // 5s
}

type GRPCConfig struct {
	ListenAddr string
}

type DBConfig struct {
	DSN                  string
	AutoMigrate          bool
	Seed                 bool
	LogQueries           bool
	ConnectRetries       int
	ConnectRetryInterval time.Duration
	// Connection pool settings
	MaxOpenConns    int           // max open connections (0 = unlimited)
	MaxIdleConns    int           // max idle connections
	ConnMaxLifetime time.Duration // maximum connection lifetime
	ConnMaxIdleTime time.Duration // maximum idle time for a connection
}

type SecurityConfig struct {
	TokenHashSecret        string // HMAC secret for opaque token hashing
	RefreshTokenHashSecret string // Deprecated: use TokenHashSecret
}

type Config struct {
	AppEnv   string
	HTTP     HTTPConfig
	GRPC     GRPCConfig
	DB       DBConfig
	Push     PushConfig
	Security SecurityConfig
}

type PushConfig struct {
	FirebaseCredentialsFile string
	IdentityBaseURL         string
}

func Load() (Config, error) {
	_ = godotenv.Load(".env") // do not fail if the file does not exist

	env := NormalizeAppEnv(get("APP_ENV", "dev"))
	port := get("GOLANG_PORT", "8888")
	grpcPort := get("GRPC_PORT", "9090")

	addr := ":" + port
	grpcAddr := ":" + grpcPort
	if IsLocalAppEnv(env) {
		addr = "0.0.0.0:" + port
		grpcAddr = "0.0.0.0:" + grpcPort
	}

	legacyTokenHashSecret := get("REFRESH_TOKEN_HASH_SECRET", "")
	tokenHashSecret := get("TOKEN_HASH_SECRET", legacyTokenHashSecret)

	cfg := Config{
		AppEnv: env,
		HTTP: HTTPConfig{
			ListenAddr:        addr,
			ShutdownGrace:     10 * time.Second,
			ReadHeaderTimeout: getDuration("HTTP_READ_HEADER_TIMEOUT", 5*time.Second),
		},
		GRPC: GRPCConfig{
			ListenAddr: grpcAddr,
		},
		DB: DBConfig{
			DSN:                  pgDSN(env),
			AutoMigrate:          get("DB_AUTO_MIGRATE", "true") == "true",
			Seed:                 get("DB_SEED", "false") == "true",
			LogQueries:           get("DB_LOG", "false") == "true",
			ConnectRetries:       getInt("DB_CONNECT_RETRIES", 10),
			ConnectRetryInterval: getDuration("DB_CONNECT_RETRY_INTERVAL", 2*time.Second),
			// Connection pool: sensible production defaults
			MaxOpenConns:    getInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getDuration("DB_CONN_MAX_IDLE_TIME", 1*time.Minute),
		},
		Push: PushConfig{
			FirebaseCredentialsFile: get("PUSH_FIREBASE_CREDENTIALS_FILE", ""),
			IdentityBaseURL:         get("PUSH_IDENTITY_BASE_URL", "http://localhost:8888"),
		},
		Security: SecurityConfig{
			TokenHashSecret:        tokenHashSecret,
			RefreshTokenHashSecret: legacyTokenHashSecret,
		},
	}
	return cfg, nil
}

func pgDSN(appEnv string) string {
	if url := strings.TrimSpace(get("DB_URL", "")); url != "" {
		return url
	}
	host := get("DB_HOST", "postgres")
	user := get("DB_USER", "postgres")
	pass := get("DB_PASS", "")
	name := get("DB_NAME", "postgres")
	port := get("DB_PORT", "5432")

	if IsLocalAppEnv(appEnv) {
		host = get("DB_HOST_LOCAL", "127.0.0.1")
		port = get("DB_PORT", port)
	}

	var builder strings.Builder
	builder.WriteString("host=")
	builder.WriteString(host)
	builder.WriteString(" user=")
	builder.WriteString(user)
	builder.WriteString(" password=")
	builder.WriteString(pass)
	builder.WriteString(" dbname=")
	builder.WriteString(name)
	builder.WriteString(" port=")
	builder.WriteString(port)
	return builder.String()
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return def
}

func getDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func NormalizeAppEnv(value string) string {
	env := strings.ToLower(strings.TrimSpace(value))
	switch env {
	case "", "dev", "development", "local", "localhost":
		return "localhost"
	case "prod", "production":
		return "production"
	default:
		return env
	}
}

func IsLocalAppEnv(value string) bool {
	return NormalizeAppEnv(value) == "localhost"
}

func IsProductionAppEnv(value string) bool {
	return NormalizeAppEnv(value) == "production"
}
