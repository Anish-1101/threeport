/*
Copyright © 2023 Threeport admin@threeport.io
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/threeport/threeport/internal/workload/status"
	cli "github.com/threeport/threeport/pkg/cli/v0"
	client "github.com/threeport/threeport/pkg/client/v0"
	config "github.com/threeport/threeport/pkg/config/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
)

// GetWorkloadInstancesCmd represents the workload-instances command
var GetWorkloadInstancesCmd = &cobra.Command{
	Use:          "workload-instances",
	Example:      "tptctl get workload-instances",
	Short:        "Get workload instances from the system",
	Long:         `Get workload instances from the system.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		// get threeport config and extract threeport API endpoint
		threeportConfig, requestedInstance, err := config.GetThreeportConfig(cliArgs.InstanceName)
		if err != nil {
			cli.Error("failed to get threeport config", err)
			os.Exit(1)
		}
		apiEndpoint, err := threeportConfig.GetThreeportAPIEndpoint(requestedInstance)
		if err != nil {
			cli.Error("failed to get threeport API endpoint from config", err)
			os.Exit(1)
		}

		// get threeport API client
		cliArgs.AuthEnabled, err = threeportConfig.GetThreeportAuthEnabled(requestedInstance)
		if err != nil {
			cli.Error("failed to determine if auth is enabled on threeport API", err)
			os.Exit(1)
		}
		ca, clientCertificate, clientPrivateKey, err := threeportConfig.GetThreeportCertificatesForInstance(requestedInstance)
		if err != nil {
			cli.Error("failed to get threeport certificates from config", err)
			os.Exit(1)
		}
		apiClient, err := client.GetHTTPClient(cliArgs.AuthEnabled, ca, clientCertificate, clientPrivateKey, "")
		if err != nil {
			cli.Error("failed to create threeport API client", err)
			os.Exit(1)
		}

		// get workload instances
		workloadInstances, err := client.GetWorkloadInstances(apiClient, apiEndpoint)
		if err != nil {
			cli.Error("failed to retrieve workload instances", err)
			os.Exit(1)
		}

		// write the output
		if len(*workloadInstances) == 0 {
			cli.Info(fmt.Sprintf(
				"No workload instances currently managed by %s threeport control plane",
				requestedInstance,
			))
			os.Exit(0)
		}
		writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
		fmt.Fprintln(writer, "NAME\t WORKLOAD DEFINITION\t KUBERNETES RUNTIME INSTANCE\t STATUS\t AGE")
		metadataErr := false
		var workloadDefErr error
		var kubernetesRuntimeInstErr error
		var statusErr error
		for _, wi := range *workloadInstances {
			// get workload definition name for instance
			var workloadDef string
			workloadDefinition, err := client.GetWorkloadDefinitionByID(apiClient, apiEndpoint, *wi.WorkloadDefinitionID)
			if err != nil {
				metadataErr = true
				workloadDefErr = err
				workloadDef = "<error>"
			} else {
				workloadDef = *workloadDefinition.Name
			}
			// get kubernetes runtime instance name for instance
			var kubernetesRuntimeInst string
			kubernetesRuntimeInstance, err := client.GetKubernetesRuntimeInstanceByID(apiClient, apiEndpoint, *wi.KubernetesRuntimeInstanceID)
			if err != nil {
				metadataErr = true
				kubernetesRuntimeInstErr = err
				kubernetesRuntimeInst = "<error>"
			} else {
				kubernetesRuntimeInst = *kubernetesRuntimeInstance.Name
			}
			// get workload status
			var workloadInstStatus string
			workloadInstStatusDetail := status.GetWorkloadInstanceStatus(apiClient, apiEndpoint, &wi)
			if workloadInstStatusDetail.Error != nil {
				metadataErr = true
				statusErr = workloadInstStatusDetail.Error
				workloadInstStatus = "<error>"
			}
			workloadInstStatus = string(workloadInstStatusDetail.Status)
			fmt.Fprintln(
				writer, *wi.Name, "\t", workloadDef, "\t", kubernetesRuntimeInst, "\t",
				workloadInstStatus, "\t", util.GetAge(wi.CreatedAt),
			)
		}
		writer.Flush()

		if metadataErr {
			if workloadDefErr != nil {
				cli.Error("encountered an error retrieving workload definition info", workloadDefErr)
			}
			if kubernetesRuntimeInstErr != nil {
				cli.Error("encountered an error retrieving kubernetes runtime instance info", kubernetesRuntimeInstErr)
			}
			if statusErr != nil {
				cli.Error("encountered an error retrieving workload instance status", statusErr)
			}
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(GetWorkloadInstancesCmd)
	GetWorkloadInstancesCmd.Flags().StringVarP(
		&cliArgs.InstanceName,
		"threeport-instance", "i", "", "Optional. Name of control plane instance. Will default to current instance if not provided.",
	)
}
