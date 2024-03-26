/*
Copyright © 2023 Threeport admin@threeport.io
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	cli "github.com/threeport/threeport/pkg/cli/v0"
	client "github.com/threeport/threeport/pkg/client/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
)

// GetKubernetesRuntimeDefinitionsCmd represents the kubernetes-runtime-definitions command
var GetKubernetesRuntimeDefinitionsCmd = &cobra.Command{
	Use:          "kubernetes-runtime-definitions",
	Example:      "tptctl get kubernetes-runtime-definitions",
	Short:        "Get kubernetes runtime definitions from the system",
	Long:         `Get kubernetes runtime definitions from the system.`,
	SilenceUsage: true,
	PreRun:       commandPreRunFunc,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, _, apiEndpoint, requestedControlPlane := getClientContext(cmd)

		// get kubernetes runtime definitions
		kubernetesRuntimeDefinitions, err := client.GetKubernetesRuntimeDefinitions(apiClient, apiEndpoint)
		if err != nil {
			cli.Error("failed to retrieve kubernetes runtime definitions", err)
			os.Exit(1)
		}

		// write the output
		if len(*kubernetesRuntimeDefinitions) == 0 {
			cli.Info(fmt.Sprintf(
				"No kubernetes runtime definitions currently managed by %s threeport control plane",
				requestedControlPlane,
			))
			os.Exit(0)
		}
		writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
		fmt.Fprintln(writer, "NAME\t INFRA PROVIDER\t HIGH AVAILABILITY\t INFRA PROVIDER ACCOUNT\t AGE")
		metadataErr := false
		var kubernetesRuntimeDefErr error
		var kubernetesRuntimeInstErr error
		var statusErr error
		for _, krd := range *kubernetesRuntimeDefinitions {
			var ha bool
			if krd.HighAvailability == nil {
				ha = false
			} else {
				ha = *krd.HighAvailability
			}
			var providerAccountID string
			if krd.InfraProviderAccountName == nil {
				providerAccountID = "N/A"
			} else {
				providerAccountID = *krd.InfraProviderAccountName
			}
			fmt.Fprintln(
				writer, *krd.Name, "\t", *krd.InfraProvider, "\t", ha, "\t",
				providerAccountID, "\t", util.GetAge(krd.CreatedAt),
			)
		}
		writer.Flush()

		if metadataErr {
			if kubernetesRuntimeDefErr != nil {
				cli.Error("encountered an error retrieving kubernetes runtime definition info", kubernetesRuntimeDefErr)
			}
			if kubernetesRuntimeInstErr != nil {
				cli.Error("encountered an error retrieving kubernetes runtime definition info", kubernetesRuntimeInstErr)
			}
			if statusErr != nil {
				cli.Error("encountered an error retrieving kubernetes runtime definition status", statusErr)
			}
			os.Exit(1)
		}
	},
}

func init() {
	GetCmd.AddCommand(GetKubernetesRuntimeDefinitionsCmd)
	GetKubernetesRuntimeDefinitionsCmd.Flags().StringVarP(
		&cliArgs.ControlPlaneName,
		"control-plane-name", "i", "", "Optional. Name of control plane. Will default to current control plane if not provided.",
	)
}