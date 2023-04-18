/*
Copyright © 2023 Threeport admin@threeport.io
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/threeport/threeport/internal/cli"
	config "github.com/threeport/threeport/pkg/config/v0"
)

var createWorkloadDefinitionConfigPath string

// CreateWorkloadDefinitionCmd represents the workload-definition command
var CreateWorkloadDefinitionCmd = &cobra.Command{
	Use:          "workload-definition",
	Example:      "tptctl create workload-definition --config /path/to/config.yaml",
	Short:        "Create a new workload definition",
	Long:         `Create a new workload definition.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		// get threeport config and extract threeport API endpoint
		threeportConfig := &config.ThreeportConfig{}
		if err := viper.Unmarshal(threeportConfig); err != nil {
			cli.Error("Failed to get threeport config", err)
			os.Exit(1)
		}
		apiEndpoint, err := threeportConfig.GetThreeportAPIEndpoint()
		if err != nil {
			cli.Error("failed to get threeport API endpoint from config", err)
			os.Exit(1)
		}

		// load workload definition config
		configContent, err := ioutil.ReadFile(createWorkloadDefinitionConfigPath)
		if err != nil {
			cli.Error("failed to read config file", err)
			os.Exit(1)
		}
		var workloadDefinitionConfig config.WorkloadDefinitionConfig
		if err := yaml.Unmarshal(configContent, &workloadDefinitionConfig); err != nil {
			cli.Error("failed to unmarshal config file yaml content", err)
			os.Exit(1)
		}

		// create workload definition
		workloadDefinition := workloadDefinitionConfig.WorkloadDefinition
		wd, err := workloadDefinition.Create(apiEndpoint)
		if err != nil {
			cli.Error("failed to create workload definition", err)
			os.Exit(1)
		}

		cli.Complete(fmt.Sprintf("workload definition %s created", *wd.Name))
	},
}

func init() {
	createCmd.AddCommand(CreateWorkloadDefinitionCmd)

	CreateWorkloadDefinitionCmd.Flags().StringVarP(
		&createWorkloadDefinitionConfigPath,
		"config", "c", "", "Path to file with workload definition config.",
	)
	CreateWorkloadDefinitionCmd.MarkFlagRequired("config")
}