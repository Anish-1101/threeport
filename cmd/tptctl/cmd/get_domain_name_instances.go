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

// GetDomainNameInstancesCmd represents the domain-name-instances command
var GetDomainNameInstancesCmd = &cobra.Command{
	Use:          "domain-name-instances",
	Example:      "tptctl get domain-name-instances",
	Short:        "Get domain name instances from the system",
	Long:         `Get domain name instances from the system.`,
	SilenceUsage: true,
	PreRun:       commandPreRunFunc,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, _, apiEndpoint, requestedControlPlane := getClientContext(cmd)

		// get domain name instances
		domainNameInstances, err := client.GetDomainNameInstances(apiClient, apiEndpoint)
		if err != nil {
			cli.Error("failed to retrieve domain name instances", err)
			os.Exit(1)
		}

		// write the output
		if len(*domainNameInstances) == 0 {
			cli.Info(fmt.Sprintf(
				"No domain name instances currently managed by %s threeport control plane",
				requestedControlPlane,
			))
			os.Exit(0)
		}
		writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
		fmt.Fprintln(writer, "NAME\t GATEWAY DEFINITION\t KUBERNETES RUNTIME INSTANCE\t WORKLOAD INSTANCE\t AGE")
		metadataErr := false
		var domainNameDefErr error
		var kubernetesRuntimeInstErr error
		var workloadInstErr error
		for _, d := range *domainNameInstances {
			// get domain name definition name for instance
			var domainNameDef string
			domainNameDefinition, err := client.GetDomainNameDefinitionByID(apiClient, apiEndpoint, *d.DomainNameDefinitionID)
			if err != nil {
				metadataErr = true
				domainNameDefErr = err
				domainNameDef = "<error>"
			} else {
				domainNameDef = *domainNameDefinition.Name
			}
			// get kubernetes runtime instance name for instance
			var kubernetesRuntimeInst string
			kubernetesRuntimeInstance, err := client.GetKubernetesRuntimeInstanceByID(apiClient, apiEndpoint, *d.KubernetesRuntimeInstanceID)
			if err != nil {
				metadataErr = true
				kubernetesRuntimeInstErr = err
				kubernetesRuntimeInst = "<error>"
			} else {
				kubernetesRuntimeInst = *kubernetesRuntimeInstance.Name
			}
			// get workload instance instance name for instance
			var workloadInst string
			workloadInstance, err := client.GetWorkloadInstanceByID(apiClient, apiEndpoint, *d.WorkloadInstanceID)
			if err != nil {
				metadataErr = true
				workloadInstErr = err
				workloadInst = "<error>"
			} else {
				workloadInst = *workloadInstance.Name
			}
			fmt.Fprintln(
				writer, *d.Name, "\t", domainNameDef, "\t", kubernetesRuntimeInst, "\t",
				workloadInst, "\t", util.GetAge(d.CreatedAt),
			)
		}
		writer.Flush()

		if metadataErr {
			if domainNameDefErr != nil {
				cli.Error("encountered an error retrieving domain name definition info", domainNameDefErr)
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
	GetCmd.AddCommand(GetDomainNameInstancesCmd)
	GetDomainNameInstancesCmd.Flags().StringVarP(
		&cliArgs.ControlPlaneName,
		"control-plane-name", "i", "", "Optional. Name of control plane. Will default to current control plane if not provided.",
	)
}
