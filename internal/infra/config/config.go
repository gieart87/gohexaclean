package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	App       AppConfig       `yaml:"app"`
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
	Logger    LoggerConfig    `yaml:"logger"`
	JWT       JWTConfig       `yaml:"jwt"`
	CORS      CORSConfig      `yaml:"cors"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Telemetry TelemetryConfig `yaml:"telemetry"`
	Metrics   MetricsConfig   `yaml:"metrics"`
}

type AppConfig struct {
	Name  string `yaml:"name"`
	Env   string `yaml:"env"`
	Debug bool   `yaml:"debug"`
}

type ServerConfig struct {
	HTTP HTTPConfig `yaml:"http"`
	GRPC GRPCConfig `yaml:"grpc"`
}

type HTTPConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type GRPCConfig struct {
	Port              int           `yaml:"port"`
	MaxConnectionIdle time.Duration `yaml:"max_connection_idle"`
	MaxConnectionAge  time.Duration `yaml:"max_connection_age"`
}

type DatabaseConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	User         string        `yaml:"user"`
	Password     string        `yaml:"password"`
	Name         string        `yaml:"name"`
	SSLMode      string        `yaml:"ssl_mode"`
	MaxOpenConns int           `yaml:"max_open_conns"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxLifetime  time.Duration `yaml:"max_lifetime"`
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

type JWTConfig struct {
	Secret  string        `yaml:"secret"`
	Expired time.Duration `yaml:"expired"`
}

type CORSConfig struct {
	AllowOrigins []string `yaml:"allow_origins"`
	AllowMethods []string `yaml:"allow_methods"`
	AllowHeaders []string `yaml:"allow_headers"`
}

type RateLimitConfig struct {
	Enabled bool          `yaml:"enabled"`
	Max     int           `yaml:"max"`
	Window  time.Duration `yaml:"window"`
}

type TelemetryConfig struct {
	Enabled           bool   `yaml:"enabled"`
	ServiceName       string `yaml:"service_name"`
	CollectorEndpoint string `yaml:"collector_endpoint"`
}

type MetricsConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

// Load loads configuration from YAML file and environment variables
func Load(configPath string) (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	// Read YAML config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables
	overrideFromEnv(&cfg)

	return &cfg, nil
}

// overrideFromEnv overrides config values with environment variables if they exist
func overrideFromEnv(cfg *Config) {
	if v := os.Getenv("APP_NAME"); v != "" {
		cfg.App.Name = v
	}
	if v := os.Getenv("APP_ENV"); v != "" {
		cfg.App.Env = v
	}
	if v := os.Getenv("APP_DEBUG"); v == "true" {
		cfg.App.Debug = true
	}

	if v := os.Getenv("HTTP_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Server.HTTP.Port)
	}
	if v := os.Getenv("GRPC_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Server.GRPC.Port)
	}

	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Database.Port)
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.Name = v
	}

	if v := os.Getenv("REDIS_HOST"); v != "" {
		cfg.Redis.Host = v
	}
	if v := os.Getenv("REDIS_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Redis.Port)
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}

	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}

	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Logger.Level = v
	}
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// GetRedisAddr returns the Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
