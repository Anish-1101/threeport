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

var createKubernetesRuntimeConfigPath string

// CreateKubernetesRuntimeCmd represents the kubernetes-runtime command
var CreateKubernetesRuntimeCmd = &cobra.Command{
	Use:     "kubernetes-runtime",
	Example: "tptctl create kubernetes-runtime --config /path/to/config.yaml",
	Short:   "Create a new kubernetes runtime",
	Long: `Create a new kubernetes runtime. This command creates a new kubernetes runtime definition
and kubernetes runtime instance based on the kubernetes runtime config.`,
	SilenceUsage: true,
	PreRun:       commandPreRunFunc,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, _, apiEndpoint, _ := getClientContext(cmd)

		// load kubernetes runtime config
		configContent, err := os.ReadFile(createKubernetesRuntimeConfigPath)
		if err != nil {
			cli.Error("failed to read config file", err)
			os.Exit(1)
		}
		var kubernetesRuntimeConfig config.KubernetesRuntimeConfig
		if err := yaml.Unmarshal(configContent, &kubernetesRuntimeConfig); err != nil {
			cli.Error("failed to unmarshal config file yaml content", err)
			os.Exit(1)
		}

		// create kubernetes runtime
		kubernetesRuntime := kubernetesRuntimeConfig.KubernetesRuntime
		krd, kri, err := kubernetesRuntime.Create(apiClient, apiEndpoint)
		if err != nil {
			cli.Error("failed to create kubernetes runtime", err)
			os.Exit(1)
		}

		cli.Info(fmt.Sprintf("kubernetes runtime definition %s created", *krd.Name))
		cli.Info(fmt.Sprintf("kubernetes runtime instance %s created", *kri.Name))
		cli.Complete(fmt.Sprintf("kubernetes runtime %s created", kubernetesRuntimeConfig.KubernetesRuntime.Name))
	},
}

func init() {
	CreateCmd.AddCommand(CreateKubernetesRuntimeCmd)

	CreateKubernetesRuntimeCmd.Flags().StringVarP(
		&createKubernetesRuntimeConfigPath,
		"config", "c", "", "Path to file with kubernetes runtime config.",
	)
	CreateKubernetesRuntimeCmd.MarkFlagRequired("config")
	CreateKubernetesRuntimeCmd.Flags().StringVarP(
		&cliArgs.ControlPlaneName,
		"control-plane-name", "i", "", "Optional. Name of control plane. Will default to current control plane if not provided.",
	)
}