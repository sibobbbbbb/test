package main

import (
    "context"
    "fmt"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "go.uber.org/zap"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"

    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/config"
    grpcHandler "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/delivery/grpc"
    httpHandler "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/delivery/http"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/repository/postgres"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/repository/redis"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/usecase"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/pkg/token"
    pb "github.com/sibobbbbbb/backend-engineer-challenge/proto/user"
)

func main() {
    // Inisialisasi logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    zap.ReplaceGlobals(logger)

    // Load konfigurasi
    cfg := config.LoadConfig()

    // Koneksi ke database
    db, err := initDB(cfg)
    if err != nil {
        logger.Fatal("Failed to connect to database", zap.Error(err))
    }
    defer db.Close()

    // Inisialisasi repository
    userRepo := postgres.NewUserRepository(db)
    redisRepo := redis.NewRedisRepository(
        cfg.GetRedisConnectionString(),
        cfg.RedisPassword,
        cfg.RedisDB,
    )

    // Inisialisasi token manager
    tokenManager := token.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)

    // Inisialisasi usecase
    userUsecase := usecase.NewUserUsecase(
        userRepo,    // Untuk operasi database
        redisRepo,   // Untuk operasi token
        tokenManager,
        cfg.ContextTimeout,
    )

    // Buat HTTP router
    router := mux.NewRouter()
    httpHandler.NewUserHandler(router, userUsecase)

    // Buat gRPC server
    grpcServer := grpc.NewServer()
    userHandler := grpcHandler.NewUserHandler(userUsecase)
    pb.RegisterUserServiceServer(grpcServer, userHandler)
    reflection.Register(grpcServer)

    // Start HTTP server
    httpServer := &http.Server{
        Addr:    fmt.Sprintf("%s:%s", cfg.ServerAddress, cfg.ServerPort),
        Handler: router,
    }

    // Jalankan server secara concurrent
    go func() {
        logger.Info("Starting HTTP server", zap.String("address", httpServer.Addr))
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("HTTP server error", zap.Error(err))
        }
    }()

    // Jalankan gRPC server
    go func() {
        lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.ServerAddress, cfg.GRPCPort))
        if err != nil {
            logger.Fatal("Failed to listen", zap.Error(err))
        }

        logger.Info("Starting gRPC server", zap.String("address", lis.Addr().String()))
        if err := grpcServer.Serve(lis); err != nil {
            logger.Fatal("gRPC server error", zap.Error(err))
        }
    }()

    // Tunggu sinyal untuk graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info("Shutting down servers...")

    // Shutdown HTTP server
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := httpServer.Shutdown(ctx); err != nil {
        logger.Error("HTTP server shutdown error", zap.Error(err))
    }

    // Shutdown gRPC server
    grpcServer.GracefulStop()

    logger.Info("Servers stopped")
}

// initDB menginisialisasi koneksi database
func initDB(cfg *config.Config) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", cfg.GetDBConnectionString())
    if err != nil {
        return nil, err
    }

    // Configure database
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    // Pastikan database terhubung
    if err := db.Ping(); err != nil {
        return nil, err
    }

    // Jalankan migrasi database
    if err := runMigration(db); err != nil {
        return nil, err
    }

    return db, nil
}

// runMigration menjalankan migrasi database
func runMigration(db *sqlx.DB) error {
    // Buat tabel users jika belum ada
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id VARCHAR(36) PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        email VARCHAR(100) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        fullname VARCHAR(100) NOT NULL,
        role VARCHAR(20) NOT NULL DEFAULT 'user',
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );
    
    CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
    CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
    `

    _, err := db.Exec(query)
    return err
}