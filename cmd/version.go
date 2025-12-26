package cmd

import (
	"github.com/dws-org/dws-event-service/configs"
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Show the service version information",
	Example: "event-service version",
	Run: func(cmd *cobra.Command, args []string) {
		envCfg := configs.GetEnvConfig()

		serviceName := envCfg.Service.Name
		if serviceName == "" {
			serviceName = "Service Template"
		}

		version := envCfg.Service.Version
		if version == "" {
			version = "v0.0.0"
		}

		fmt.Printf("%s version: %s\n", serviceName, version)
	},
}
