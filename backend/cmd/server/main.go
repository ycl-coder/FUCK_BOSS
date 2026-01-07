// Package main is the entry point of the gRPC server application.
// It initializes all dependencies and starts the gRPC server.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	contentv1 "fuck_boss/backend/api/proto/content/v1"
	"fuck_boss/backend/internal/application/content"
	"fuck_boss/backend/internal/application/search"
	"fuck_boss/backend/internal/infrastructure/config"
	"fuck_boss/backend/internal/infrastructure/logger"
	"fuck_boss/backend/internal/infrastructure/persistence/postgres"
	redispersistence "fuck_boss/backend/internal/infrastructure/persistence/redis"
	grpchandler "fuck_boss/backend/internal/presentation/grpc"
	"fuck_boss/backend/internal/presentation/middleware"
)

func main() {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.NewLoggerFromConfig(&logger.LogConfig{
		Level:            cfg.Log.Level,
		Format:           cfg.Log.Format,
		OutputPaths:      cfg.Log.OutputPaths,
		ErrorOutputPaths: cfg.Log.ErrorOutputPaths,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting gRPC server...",
		zap.String("version", "1.0.0"),
		zap.Int("grpc_port", cfg.GRPC.Port),
	)

	// Connect to PostgreSQL
	db, err := connectDatabase(cfg.Database, log)
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		log.Info("Please check your database configuration:",
			zap.String("host", cfg.Database.Host),
			zap.Int("port", cfg.Database.Port),
			zap.String("user", cfg.Database.User),
			zap.String("database", cfg.Database.DBName),
		)
		log.Info("You can configure database connection via:",
			zap.String("config_file", "config/config.yaml"),
			zap.String("env_vars", "FUCK_BOSS_DATABASE_*"),
		)
		os.Exit(1)
	}
	defer db.Close()

	// Connect to Redis
	redisClient, err := connectRedis(cfg.Redis, log)
	if err != nil {
		log.Error("Failed to connect to Redis", zap.Error(err))
		log.Info("Please check your Redis configuration:",
			zap.String("host", cfg.Redis.Host),
			zap.Int("port", cfg.Redis.Port),
		)
		log.Info("You can configure Redis connection via:",
			zap.String("config_file", "config/config.yaml"),
			zap.String("env_vars", "FUCK_BOSS_REDIS_*"),
		)
		os.Exit(1)
	}
	defer redisClient.Close()

	// Run database migrations
	if err := runMigrations(db, log); err != nil {
		log.Error("Failed to run database migrations", zap.Error(err))
		os.Exit(1)
	}

	// Initialize repositories
	postRepo := postgres.NewPostRepository(db)
	cacheRepo := redispersistence.NewCacheRepository(redisClient)
	rateLimiter := redispersistence.NewRateLimiter(redisClient)

	// Initialize use cases
	createUseCase := content.NewCreatePostUseCase(postRepo, cacheRepo, rateLimiter)
	listUseCase := content.NewListPostsUseCase(postRepo, cacheRepo)
	getUseCase := content.NewGetPostUseCase(postRepo, cacheRepo)
	searchUseCase := search.NewSearchPostsUseCase(postRepo, cacheRepo)

	// Create gRPC service
	contentService := grpchandler.NewContentService(
		createUseCase,
		listUseCase,
		getUseCase,
		searchUseCase,
	)

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(log),
			middleware.LoggingInterceptor(log),
		),
		grpc.MaxRecvMsgSize(cfg.GRPC.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.GRPC.MaxSendMsgSize),
	)

	// Register services
	contentv1.RegisterContentServiceServer(grpcServer, contentService)

	// Enable reflection for gRPC tools (e.g., grpcurl, grpcui)
	reflection.Register(grpcServer)

	// Wrap gRPC server with gRPC Web support
	wrappedServer := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			// Allow all origins for development (in production, restrict this)
			return true
		}),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			// Allow all origins for WebSocket connections
			return true
		}),
	)

	// Start HTTP server for gRPC Web
	httpAddr := fmt.Sprintf(":%d", cfg.GRPC.Port)
	httpServer := &http.Server{
		Addr: httpAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wrappedServer.IsGrpcWebRequest(r) || wrappedServer.IsAcceptableGrpcCorsRequest(r) {
				wrappedServer.ServeHTTP(w, r)
			} else {
				// Fallback to standard gRPC for non-Web requests
				grpcServer.ServeHTTP(w, r)
			}
		}),
	}

	log.Info("gRPC server listening",
		zap.String("address", httpAddr),
		zap.Int("max_recv_msg_size", cfg.GRPC.MaxRecvMsgSize),
		zap.Int("max_send_msg_size", cfg.GRPC.MaxSendMsgSize),
		zap.Bool("grpc_web_enabled", true),
	)

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("gRPC server started (with gRPC Web support)")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Wait for interrupt signal or server error
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverErrors:
		log.Error("Server error", zap.Error(err))
		os.Exit(1)
	case sig := <-shutdown:
		log.Info("Shutdown signal received", zap.String("signal", sig.String()))
	}

	// Graceful shutdown
	log.Info("Starting graceful shutdown...")

	// Create shutdown context with timeout
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelShutdown()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Error during HTTP server shutdown", zap.Error(err))
	}

	// Stop gRPC server gracefully
	grpcServer.GracefulStop()

	// Close database connection
	if err := db.Close(); err != nil {
		log.Warn("Error closing database connection", zap.Error(err))
	}

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		log.Warn("Error closing Redis connection", zap.Error(err))
	}

	// Check if shutdown completed within timeout
	select {
	case <-shutdownCtx.Done():
		log.Warn("Shutdown timeout exceeded, forcing stop")
		grpcServer.Stop()
	default:
		log.Info("Graceful shutdown completed")
	}
}

// loadConfig loads application configuration.
func loadConfig() (*config.Config, error) {
	// Try to load from config file first
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// If config file doesn't exist, try to load with defaults
		// LoadConfig will apply defaults automatically
		cfg, err = config.LoadConfig("")
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration: %w", err)
		}
	}

	return cfg, nil
}

// connectDatabase connects to PostgreSQL database.
func connectDatabase(cfg config.DatabaseConfig, log logger.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.DBName),
		zap.Int("max_open_conns", cfg.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
	)

	return db, nil
}

// connectRedis connects to Redis.
func connectRedis(cfg config.RedisConfig, log logger.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Info("Redis connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB),
		zap.Int("pool_size", cfg.PoolSize),
	)

	return client, nil
}

// runMigrations runs database migrations.
func runMigrations(db *sql.DB, log logger.Logger) error {
	// For now, we'll use a simple approach: check if tables exist
	// In production, you should use a migration tool like golang-migrate
	ctx := context.Background()

	// Check if posts table exists
	var exists bool
	err := db.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'posts')",
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if posts table exists: %w", err)
	}

	if !exists {
		log.Info("Running database migrations...")

		// Create cities table
		citiesSQL := `
		CREATE TABLE IF NOT EXISTS cities (
			code VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			pinyin VARCHAR(100),
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_cities_name ON cities(name);
		`

		if _, err := db.ExecContext(ctx, citiesSQL); err != nil {
			return fmt.Errorf("failed to create cities table: %w", err)
		}

		// Create posts table
		postsSQL := `
		CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			company_name VARCHAR(100) NOT NULL,
			city_code VARCHAR(50) NOT NULL,
			city_name VARCHAR(100) NOT NULL,
			content TEXT NOT NULL,
			occurred_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			CONSTRAINT posts_city_code_fkey FOREIGN KEY (city_code) REFERENCES cities(code)
		);

		CREATE INDEX IF NOT EXISTS idx_posts_city_code ON posts(city_code);
		CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_posts_company_name ON posts(company_name);
		CREATE INDEX IF NOT EXISTS idx_posts_search ON posts USING GIN (to_tsvector('simple', company_name || ' ' || content));
		`

		if _, err := db.ExecContext(ctx, postsSQL); err != nil {
			return fmt.Errorf("failed to create posts table: %w", err)
		}

		log.Info("Database migrations completed")
	} else {
		log.Info("Database tables already exist, skipping migrations")
	}

	return nil
}
