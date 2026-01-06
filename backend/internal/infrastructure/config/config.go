// Package config provides configuration management using viper.
// It supports loading configuration from YAML files and environment variables.
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration.
type Config struct {
	// Database contains PostgreSQL database configuration.
	Database DatabaseConfig

	// Redis contains Redis cache configuration.
	Redis RedisConfig

	// GRPC contains gRPC server configuration.
	GRPC GRPCConfig

	// Log contains logging configuration.
	Log LogConfig
}

// DatabaseConfig contains PostgreSQL database connection settings.
type DatabaseConfig struct {
	// Host is the database host address.
	Host string

	// Port is the database port number.
	Port int

	// User is the database username.
	User string

	// Password is the database password.
	Password string

	// DBName is the database name.
	DBName string

	// SSLMode is the SSL mode for database connection (disable, require, verify-ca, verify-full).
	SSLMode string

	// MaxOpenConns is the maximum number of open connections to the database.
	MaxOpenConns int

	// MaxIdleConns is the maximum number of idle connections in the pool.
	MaxIdleConns int

	// ConnMaxLifetime is the maximum amount of time a connection may be reused (in seconds).
	ConnMaxLifetime int
}

// RedisConfig contains Redis cache connection settings.
type RedisConfig struct {
	// Host is the Redis host address.
	Host string

	// Port is the Redis port number.
	Port int

	// Password is the Redis password (empty if no password).
	Password string

	// DB is the Redis database number (0-15).
	DB int

	// MaxRetries is the maximum number of retries before giving up.
	MaxRetries int

	// PoolSize is the maximum number of socket connections.
	PoolSize int

	// MinIdleConns is the minimum number of idle connections.
	MinIdleConns int
}

// GRPCConfig contains gRPC server configuration.
type GRPCConfig struct {
	// Port is the gRPC server port number.
	Port int

	// MaxRecvMsgSize is the maximum message size the server can receive (in bytes).
	MaxRecvMsgSize int

	// MaxSendMsgSize is the maximum message size the server can send (in bytes).
	MaxSendMsgSize int
}

// LogConfig contains logging configuration.
type LogConfig struct {
	// Level is the log level (debug, info, warn, error).
	Level string

	// Format is the log format (json, text, console).
	// console is an alias for text (human-readable format).
	Format string

	// OutputPaths is a list of paths to write logging output to.
	OutputPaths []string

	// ErrorOutputPaths is a list of paths to write error level logs to.
	ErrorOutputPaths []string
}

// LoadConfig loads configuration from file and environment variables.
// It reads from the specified config file path and environment variables.
// Environment variables take precedence over file configuration.
// Environment variable names should be in the format: FUCK_BOSS_<SECTION>_<FIELD>
// Example: FUCK_BOSS_DATABASE_HOST, FUCK_BOSS_REDIS_PORT
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Configure viper
	v.SetConfigType("yaml")
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Default config file locations
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("../config")
	}

	// Enable environment variable support
	v.SetEnvPrefix("FUCK_BOSS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file (optional, will not error if file doesn't exist)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is OK, we can use defaults and env vars
	}

	// Unmarshal configuration
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Apply defaults for zero values (viper may not apply defaults during unmarshal)
	applyDefaults(&cfg)

	// Validate configuration
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// applyDefaults applies default values to zero-value fields in the config.
func applyDefaults(cfg *Config) {
	// Database defaults
	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}
	if cfg.Database.User == "" {
		cfg.Database.User = "postgres"
	}
	if cfg.Database.DBName == "" {
		cfg.Database.DBName = "fuck_boss"
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 100
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 10
	}
	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 3600
	}

	// Redis defaults
	if cfg.Redis.Host == "" {
		cfg.Redis.Host = "localhost"
	}
	if cfg.Redis.Port == 0 {
		cfg.Redis.Port = 6379
	}
	if cfg.Redis.MaxRetries == 0 {
		cfg.Redis.MaxRetries = 3
	}
	if cfg.Redis.PoolSize == 0 {
		cfg.Redis.PoolSize = 50
	}
	if cfg.Redis.MinIdleConns == 0 {
		cfg.Redis.MinIdleConns = 5
	}

	// gRPC defaults
	if cfg.GRPC.Port == 0 {
		cfg.GRPC.Port = 50051
	}
	if cfg.GRPC.MaxRecvMsgSize == 0 {
		cfg.GRPC.MaxRecvMsgSize = 4194304 // 4MB
	}
	if cfg.GRPC.MaxSendMsgSize == 0 {
		cfg.GRPC.MaxSendMsgSize = 4194304 // 4MB
	}

	// Log defaults
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}
	if len(cfg.Log.OutputPaths) == 0 {
		cfg.Log.OutputPaths = []string{"stdout"}
	}
	if len(cfg.Log.ErrorOutputPaths) == 0 {
		cfg.Log.ErrorOutputPaths = []string{"stderr"}
	}
}

// setDefaults sets default configuration values.
func setDefaults(v *viper.Viper) {
	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "")
	v.SetDefault("database.dbname", "fuck_boss")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", 3600)

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.max_retries", 3)
	v.SetDefault("redis.pool_size", 50)
	v.SetDefault("redis.min_idle_conns", 5)

	// gRPC defaults
	v.SetDefault("grpc.port", 50051)
	v.SetDefault("grpc.max_recv_msg_size", 4194304) // 4MB
	v.SetDefault("grpc.max_send_msg_size", 4194304) // 4MB

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output_paths", []string{"stdout"})
	v.SetDefault("log.error_output_paths", []string{"stderr"})
}

// validateConfig validates the configuration and returns an error if validation fails.
func validateConfig(cfg *Config) error {
	// Validate database configuration
	if cfg.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return fmt.Errorf("database.port must be between 1 and 65535")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("database.user is required")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("database.dbname is required")
	}
	if cfg.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("database.max_open_conns must be greater than 0")
	}
	if cfg.Database.MaxIdleConns < 0 {
		return fmt.Errorf("database.max_idle_conns must be non-negative")
	}
	if cfg.Database.ConnMaxLifetime < 0 {
		return fmt.Errorf("database.conn_max_lifetime must be non-negative")
	}

	// Validate Redis configuration
	if cfg.Redis.Host == "" {
		return fmt.Errorf("redis.host is required")
	}
	if cfg.Redis.Port <= 0 || cfg.Redis.Port > 65535 {
		return fmt.Errorf("redis.port must be between 1 and 65535")
	}
	if cfg.Redis.DB < 0 || cfg.Redis.DB > 15 {
		return fmt.Errorf("redis.db must be between 0 and 15")
	}
	if cfg.Redis.PoolSize <= 0 {
		return fmt.Errorf("redis.pool_size must be greater than 0")
	}
	if cfg.Redis.MinIdleConns < 0 {
		return fmt.Errorf("redis.min_idle_conns must be non-negative")
	}

	// Validate gRPC configuration
	if cfg.GRPC.Port <= 0 || cfg.GRPC.Port > 65535 {
		return fmt.Errorf("grpc.port must be between 1 and 65535")
	}
	if cfg.GRPC.MaxRecvMsgSize <= 0 {
		return fmt.Errorf("grpc.max_recv_msg_size must be greater than 0")
	}
	if cfg.GRPC.MaxSendMsgSize <= 0 {
		return fmt.Errorf("grpc.max_send_msg_size must be greater than 0")
	}

	// Validate log configuration
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[strings.ToLower(cfg.Log.Level)] {
		return fmt.Errorf("log.level must be one of: debug, info, warn, error")
	}

	validFormats := map[string]bool{
		"json":    true,
		"text":    true,
		"console": true, // console is an alias for text (human-readable format)
	}
	if !validFormats[strings.ToLower(cfg.Log.Format)] {
		return fmt.Errorf("log.format must be one of: json, text, console")
	}

	return nil
}

// GetDSN returns the PostgreSQL data source name (DSN) string.
func (c *DatabaseConfig) GetDSN() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
	return dsn
}

// GetAddr returns the Redis address string in the format "host:port".
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
