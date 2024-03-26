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

var createDomainNameInstanceConfigPath string

// CreateDomainNameInstanceCmd represents the domain-name-instance command
var CreateDomainNameInstanceCmd = &cobra.Command{
	Use:          "domain-name-instance",
	Example:      "tptctl create domain-name-instance --config /path/to/config.yaml",
	Short:        "Create a new domain name instance",
	Long:         `Create a new domain name instance.`,
	SilenceUsage: true,
	PreRun:       commandPreRunFunc,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, _, apiEndpoint, _ := getClientContext(cmd)

		// load domain name instance config
		configContent, err := os.ReadFile(createDomainNameInstanceConfigPath)
		if err != nil {
			cli.Error("failed to read config file", err)
			os.Exit(1)
		}
		var domainNameInstanceConfig config.DomainNameInstanceConfig
		if err := yaml.Unmarshal(configContent, &domainNameInstanceConfig); err != nil {
			cli.Error("failed to unmarshal config file yaml content", err)
			os.Exit(1)
		}

		// create domain name instance
		domainNameInstance := domainNameInstanceConfig.DomainNameInstance
		wi, err := domainNameInstance.Create(apiClient, apiEndpoint)
		if err != nil {
			cli.Error("failed to create domain name instance", err)
			os.Exit(1)
		}

		cli.Complete(fmt.Sprintf("domain name instance %s created\n", *wi.Name))
	},
}

func init() {
	CreateCmd.AddCommand(CreateDomainNameInstanceCmd)

	CreateDomainNameInstanceCmd.Flags().StringVarP(
		&createDomainNameInstanceConfigPath,
		"config", "c", "", "Path to file with domain name instance config.",
	)
	CreateDomainNameInstanceCmd.MarkFlagRequired("config")
	CreateDomainNameInstanceCmd.Flags().StringVarP(
		&cliArgs.ControlPlaneName,
		"control-plane-name", "i", "", "Optional. Name of control plane. Will default to current control plane if not provided.",
	)
}