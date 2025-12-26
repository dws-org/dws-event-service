package cmd

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/dws-org/dws-event-service/configs"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	SilenceUsage: true,
}

// You will additionally define flags and handle configuration in your init() function.
func init() {
	configs.Init()
	envCfg := configs.GetEnvConfig()

	serviceName := strings.TrimSpace(envCfg.Service.Name)
	if serviceName == "" {
		serviceName = "Service Template"
	}

	serviceSlug := strings.TrimSpace(envCfg.Service.Slug)
	if serviceSlug == "" {
		serviceSlug = strings.ReplaceAll(strings.ToLower(serviceName), " ", "-")
	}

	serviceDescription := strings.TrimSpace(envCfg.Service.Description)
	if serviceDescription == "" {
		serviceDescription = fmt.Sprintf("%s microservice CLI", serviceName)
	}

	rootCmd.Use = serviceSlug
	rootCmd.Short = serviceDescription
	rootCmd.Long = fmt.Sprintf("%s\nService: %s\nVersion: %s", serviceDescription, serviceName, envCfg.Service.Version)
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		// Start server by default when no command is provided
		run()
	}

	rootCmd.AddCommand(ServerStartCmd) // add server start command
	rootCmd.AddCommand(VersionCmd)     // add version command
}

var embedFs embed.FS

func Execute(fs embed.FS) {
	embedFs = fs

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
