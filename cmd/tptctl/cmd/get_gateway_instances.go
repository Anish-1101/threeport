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

// GetGatewayInstancesCmd represents the gateway-instances command
var GetGatewayInstancesCmd = &cobra.Command{
	Use:          "gateway-instances",
	Example:      "tptctl get gateway-instances",
	Short:        "Get gateway instances from the system",
	Long:         `Get gateway instances from the system.`,
	SilenceUsage: true,
	PreRun:       commandPreRunFunc,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, _, apiEndpoint, requestedControlPlane := getClientContext(cmd)

		// get gateway instances
		gatewayInstances, err := client.GetGatewayInstances(apiClient, apiEndpoint)
		if err != nil {
			cli.Error("failed to retrieve gateway instances", err)
			os.Exit(1)
		}

		// write the output
		if len(*gatewayInstances) == 0 {
			cli.Info(fmt.Sprintf(
				"No gateway instances currently managed by %s threeport control plane",
				requestedControlPlane,
			))
			os.Exit(0)
		}
		writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
		fmt.Fprintln(writer, "NAME\t GATEWAY DEFINITION\t KUBERNETES RUNTIME INSTANCE\t WORKLOAD INSTANCE\t AGE")
		metadataErr := false
		var gatewayDefErr error
		var kubernetesRuntimeInstErr error
		var workloadInstErr error
		for _, g := range *gatewayInstances {
			// get gateway definition name for instance
			var gatewayDef string
			gatewayDefinition, err := client.GetGatewayDefinitionByID(apiClient, apiEndpoint, *g.GatewayDefinitionID)
			if err != nil {
				metadataErr = true
				gatewayDefErr = err
				gatewayDef = "<error>"
			} else {
				gatewayDef = *gatewayDefinition.Name
			}
			// get kubernetes runtime instance name for instance
			var kubernetesRuntimeInst string
			kubernetesRuntimeInstance, err := client.GetKubernetesRuntimeInstanceByID(apiClient, apiEndpoint, *g.KubernetesRuntimeInstanceID)
			if err != nil {
				metadataErr = true
				kubernetesRuntimeInstErr = err
				kubernetesRuntimeInst = "<error>"
			} else {
				kubernetesRuntimeInst = *kubernetesRuntimeInstance.Name
			}
			// get workload instance instance name for instance
			var workloadInst string
			workloadInstance, err := client.GetWorkloadInstanceByID(apiClient, apiEndpoint, *g.WorkloadInstanceID)
			if err != nil {
				metadataErr = true
				workloadInstErr = err
				workloadInst = "<error>"
			} else {
				workloadInst = *workloadInstance.Name
			}
			fmt.Fprintln(
				writer, *g.Name, "\t", gatewayDef, "\t", kubernetesRuntimeInst, "\t",
				workloadInst, "\t", util.GetAge(g.CreatedAt),
			)
		}
		writer.Flush()

		if metadataErr {
			if gatewayDefErr != nil {
				cli.Error("encountered an error retrieving gateway definition info", gatewayDefErr)
			}
			if kubernetesRuntimeInstErr != nil {
				cli.Error("encountered an error retrieving kubernetes runtime instance info", kubernetesRuntimeInstErr)
			}
			if workloadInstErr != nil {
				cli.Error("encountered an error retrieving workload instance info", workloadInstErr)
			}
			os.Exit(1)
		}
	},
}

func init() {
	GetCmd.AddCommand(GetGatewayInstancesCmd)
	GetGatewayInstancesCmd.Flags().StringVarP(
		&cliArgs.ControlPlaneName,
		"control-plane-name", "i", "", "Optional. Name of control plane. Will default to current control plane if not provided.",
	)
}