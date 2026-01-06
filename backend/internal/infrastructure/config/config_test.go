package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_WithFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
database:
  host: test-db
  port: 5433
  user: testuser
  password: testpass
  dbname: testdb
  sslmode: require
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600

redis:
  host: test-redis
  port: 6380
  password: redispass
  db: 1
  max_retries: 3
  pool_size: 50
  min_idle_conns: 5

grpc:
  port: 50052
  max_recv_msg_size: 4194304
  max_send_msg_size: 4194304

log:
  level: debug
  format: text
  output_paths:
    - stdout
  error_output_paths:
    - stderr
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Database.Host != "test-db" {
		t.Errorf("Database.Host = %v, want test-db", cfg.Database.Host)
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("Database.Port = %v, want 5433", cfg.Database.Port)
	}
	if cfg.Redis.Host != "test-redis" {
		t.Errorf("Redis.Host = %v, want test-redis", cfg.Redis.Host)
	}
	if cfg.GRPC.Port != 50052 {
		t.Errorf("GRPC.Port = %v, want 50052", cfg.GRPC.Port)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %v, want debug", cfg.Log.Level)
	}
}

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Load config without file (should use defaults)
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Check defaults
	if cfg.Database.Host != "localhost" {
		t.Errorf("Database.Host = %v, want localhost", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("Database.Port = %v, want 5432", cfg.Database.Port)
	}
	if cfg.Redis.Port != 6379 {
		t.Errorf("Redis.Port = %v, want 6379", cfg.Redis.Port)
	}
	if cfg.GRPC.Port != 50051 {
		t.Errorf("GRPC.Port = %v, want 50051", cfg.GRPC.Port)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level = %v, want info", cfg.Log.Level)
	}
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("FUCK_BOSS_DATABASE_HOST", "env-db")
	os.Setenv("FUCK_BOSS_DATABASE_PORT", "5434")
	os.Setenv("FUCK_BOSS_REDIS_HOST", "env-redis")
	defer func() {
		os.Unsetenv("FUCK_BOSS_DATABASE_HOST")
		os.Unsetenv("FUCK_BOSS_DATABASE_PORT")
		os.Unsetenv("FUCK_BOSS_REDIS_HOST")
	}()

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Environment variables should override defaults
	if cfg.Database.Host != "env-db" {
		t.Errorf("Database.Host = %v, want env-db", cfg.Database.Host)
	}
	if cfg.Database.Port != 5434 {
		t.Errorf("Database.Port = %v, want 5434", cfg.Database.Port)
	}
	if cfg.Redis.Host != "env-redis" {
		t.Errorf("Redis.Host = %v, want env-redis", cfg.Redis.Host)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				Database: DatabaseConfig{
					Host:         "localhost",
					Port:         5432,
					User:         "postgres",
					Password:     "password",
					DBName:       "testdb",
					SSLMode:      "disable",
					MaxOpenConns: 100,
					MaxIdleConns: 10,
				},
				Redis: RedisConfig{
					Host:         "localhost",
					Port:         6379,
					Password:     "",
					DB:           0,
					PoolSize:     50,
					MinIdleConns: 5,
				},
				GRPC: GRPCConfig{
					Port:           50051,
					MaxRecvMsgSize: 4194304,
					MaxSendMsgSize: 4194304,
				},
				Log: LogConfig{
					Level:  "info",
					Format: "json",
				},
			},
			wantErr: false,
		},
		{
			name: "missing database host",
			cfg: &Config{
				Database: DatabaseConfig{
					Host: "",
					Port: 5432,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid database port",
			cfg: &Config{
				Database: DatabaseConfig{
					Host: "localhost",
					Port: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			cfg: &Config{
				Database: DatabaseConfig{
					Host:         "localhost",
					Port:         5432,
					User:         "postgres",
					DBName:       "testdb",
					MaxOpenConns: 100,
				},
				Redis: RedisConfig{
					Host:     "localhost",
					Port:     6379,
					PoolSize: 50,
				},
				GRPC: GRPCConfig{
					Port:           50051,
					MaxRecvMsgSize: 4194304,
					MaxSendMsgSize: 4194304,
				},
				Log: LogConfig{
					Level:  "invalid",
					Format: "json",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log format",
			cfg: &Config{
				Database: DatabaseConfig{
					Host:         "localhost",
					Port:         5432,
					User:         "postgres",
					DBName:       "testdb",
					MaxOpenConns: 100,
				},
				Redis: RedisConfig{
					Host:     "localhost",
					Port:     6379,
					PoolSize: 50,
				},
				GRPC: GRPCConfig{
					Port:           50051,
					MaxRecvMsgSize: 4194304,
					MaxSendMsgSize: 4194304,
				},
				Log: LogConfig{
					Level:  "info",
					Format: "invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	dsn := cfg.GetDSN()
	expected := "host=localhost port=5432 user=postgres password=password dbname=testdb sslmode=disable"

	if dsn != expected {
		t.Errorf("GetDSN() = %v, want %v", dsn, expected)
	}
}

func TestRedisConfig_GetAddr(t *testing.T) {
	cfg := RedisConfig{
		Host: "localhost",
		Port: 6379,
	}

	addr := cfg.GetAddr()
	expected := "localhost:6379"

	if addr != expected {
		t.Errorf("GetAddr() = %v, want %v", addr, expected)
	}
}
