/*
Copyright © 2023 Threeport admin@threeport.io
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	cli "github.com/threeport/threeport/pkg/cli/v0"
	config "github.com/threeport/threeport/pkg/config/v0"
)

var createGatewayInstanceConfigPath string

// CreateGatewayInstanceCmd represents the gateway-instance command
var CreateGatewayInstanceCmd = &cobra.Command{
	Use:          "gateway-instance",
	Example:      "tptctl create gateway-instance --config /path/to/config.yaml",
	Short:        "Create a new gateway instance",
	Long:         `Create a new gateway instance.`,
	SilenceUsage: true,
	PreRun:       commandPreRunFunc,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, _, apiEndpoint, _ := getClientContext(cmd)

		// load gateway instance config
		configContent, err := os.ReadFile(createGatewayInstanceConfigPath)
		if err != nil {
			cli.Error("failed to read config file", err)
			os.Exit(1)
		}
		var gatewayInstanceConfig config.GatewayInstanceConfig
		if err := yaml.Unmarshal(configContent, &gatewayInstanceConfig); err != nil {
			cli.Error("failed to unmarshal config file yaml content", err)
			os.Exit(1)
		}

		// create gateway instance
		gatewayInstance := gatewayInstanceConfig.GatewayInstance
		wi, err := gatewayInstance.Create(apiClient, apiEndpoint)
		if err != nil {
			cli.Error("failed to create gateway instance", err)
			os.Exit(1)
		}

		cli.Complete(fmt.Sprintf("gateway instance %s created\n", *wi.Name))
	},
}

func init() {
	CreateCmd.AddCommand(CreateGatewayInstanceCmd)

	CreateGatewayInstanceCmd.Flags().StringVarP(
		&createGatewayInstanceConfigPath,
		"config", "c", "", "Path to file with gateway instance config.",
	)
	CreateGatewayInstanceCmd.MarkFlagRequired("config")
	CreateGatewayInstanceCmd.Flags().StringVarP(
		&cliArgs.ControlPlaneName,
		"control-plane-name", "i", "", "Optional. Name of control plane. Will default to current control plane if not provided.",
	)
}