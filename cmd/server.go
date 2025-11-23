package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oskargbc/dws-event-service.git/configs"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/logger"
	"github.com/oskargbc/dws-event-service.git/internal/router"
	"github.com/oskargbc/dws-event-service.git/internal/services"

	"github.com/spf13/cobra"
)

var (
	// host           string
	// port           int
	ServerStartCmd = &cobra.Command{
		Use:   "server",
		Short: `Start the server`,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func run() {
	envConfig := configs.GetEnvConfig()

	// Initialize database connection first - this will panic if connection fails
	dbService := services.GetDatabaseSeviceInstance()
	defer dbService.DbDisconnect()

	// Verify database connection is working before starting server
	logger := logger.NewLogrusLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbService.HealthCheck(ctx); err != nil {
		logger.Fatalf("Database health check failed: %v. Make sure DATABASE_URL is set correctly and database is accessible.", err)
	}
	logger.Infoln("Database connection verified successfully")

	router := router.NewGinRouter(envConfig.Server.GinMode)

	server := &http.Server{
		Addr:    envConfig.Server.Port,
		Handler: router,
	}
	logger.Infof("Starting %s on port %s", envConfig.Service.Name, envConfig.Server.Port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	i := <-quit
	logger.Println("Server receive a signal: ", i.String())

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("server shutdown error: %s\n", err)
	}
	logger.Println("Server exiting")
}
