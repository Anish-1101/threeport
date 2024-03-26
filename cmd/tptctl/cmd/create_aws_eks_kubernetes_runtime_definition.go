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

var createAwsEksKubernetesRuntimeDefinitionConfigPath string

// CreateAwsEksKubernetesRuntimeDefinitionCmd represents the aws-eks-kubernetes-runtime-definition command
var CreateAwsEksKubernetesRuntimeDefinitionCmd = &cobra.Command{
	Use:          "aws-eks-kubernetes-runtime-definition",
	Example:      "tptctl create aws-eks-kubernetes-runtime-definition --config /path/to/config.yaml",
	Short:        "Create a new AWS EKS kubernetes runtime definition",
	Long:         `Create a new AWS EKS kubernetes runtime definition.`,
	SilenceUsage: true,
	PreRun:       commandPreRunFunc,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, _, apiEndpoint, _ := getClientContext(cmd)

		// load AWS EKS kubernetes runtime definition config
		configContent, err := os.ReadFile(createAwsEksKubernetesRuntimeDefinitionConfigPath)
		if err != nil {
			cli.Error("failed to read config file", err)
			os.Exit(1)
		}
		var awsEksKubernetesRuntimeDefinitionConfig config.AwsEksKubernetesRuntimeDefinitionConfig
		if err := yaml.Unmarshal(configContent, &awsEksKubernetesRuntimeDefinitionConfig); err != nil {
			cli.Error("failed to unmarshal config file yaml content", err)
			os.Exit(1)
		}

		// create AWS EKS kubernetes runtime definition
		awsEksKubernetesRuntimeDefinition := awsEksKubernetesRuntimeDefinitionConfig.AwsEksKubernetesRuntimeDefinition
		wd, err := awsEksKubernetesRuntimeDefinition.Create(apiClient, apiEndpoint)
		if err != nil {
			cli.Error("failed to create AWS EKS kubernetes runtime definition", err)
			os.Exit(1)
		}

		cli.Complete(fmt.Sprintf("kubernetes runtime definition %s created", *wd.Name))
	},
}

func init() {
	CreateCmd.AddCommand(CreateAwsEksKubernetesRuntimeDefinitionCmd)

	CreateAwsEksKubernetesRuntimeDefinitionCmd.Flags().StringVarP(
		&createAwsEksKubernetesRuntimeDefinitionConfigPath,
		"config", "c", "", "Path to file with AWS EKS kubernetes runtime definition config.",
	)
	CreateAwsEksKubernetesRuntimeDefinitionCmd.MarkFlagRequired("config")
	CreateAwsEksKubernetesRuntimeDefinitionCmd.Flags().StringVarP(
		&cliArgs.ControlPlaneName,
		"control-plane-name", "i", "", "Optional. Name of control plane. Will default to current control plane if not provided.",
	)
}