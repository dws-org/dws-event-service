package cmd

import (
	"github.com/oskargbc/dws-event-service.git/configs"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/logger"
	"github.com/oskargbc/dws-event-service.git/internal/router"
	"github.com/oskargbc/dws-event-service.git/internal/services"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	dbService := services.GetDatabaseSeviceInstance()
	defer dbService.DbDisconnect()

	router := router.NewGinRouter(envConfig.Server.GinMode)

	logger := logger.NewLogrusLogger()

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("server shutdown error: %s\n", err)
	}
	logger.Println("Server exiting")
}
