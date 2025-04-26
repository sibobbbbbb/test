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

	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/config"
	grpcHandler "github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/delivery/grpc"
	httpHandler "github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/delivery/http"
	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/repository/postgres"
	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/usecase"
	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/pkg/grpc_client"
	pb "github.com/sibobbbbbb/backend-engineer-challenge/proto/book"
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

	// Inisialisasi gRPC client untuk Category Service
	categoryClient, err := grpc_client.NewCategoryClient(cfg.GetCategoryServiceAddress())
	if err != nil {
		logger.Fatal("Failed to create category client", zap.Error(err))
	}

	// Inisialisasi repository
	bookRepo := postgres.NewBookRepository(db)

	// Inisialisasi usecase
	bookUsecase := usecase.NewBookUsecase(bookRepo, categoryClient)

	// Buat HTTP router
	router := mux.NewRouter()
	httpHandler.NewBookHandler(router, bookUsecase)

	// Buat gRPC server
	grpcServer := grpc.NewServer()
	bookHandler := grpcHandler.NewBookHandler(bookUsecase)
	pb.RegisterBookServiceServer(grpcServer, bookHandler)
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
	// Buat tabel books jika belum ada
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id VARCHAR(36) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		isbn VARCHAR(50) NOT NULL UNIQUE,
		published_year INTEGER NOT NULL,
		category_ids TEXT[] NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
	CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);
	CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);
	`

	_, err := db.Exec(query)
	return err
}