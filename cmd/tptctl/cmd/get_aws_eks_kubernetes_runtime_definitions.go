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

// GetAwsEksKubernetesRuntimeInstancesCmd represents the aws-eks-kubernetes-runtime-instances command
var GetAwsEksKubernetesRuntimeInstancesCmd = &cobra.Command{
	Use:          "aws-eks-kubernetes-runtime-instances",
	Example:      "tptctl get aws-eks-kubernetes-runtime-instances",
	Short:        "Get AWS EKS kubernetes runtime instances from the system",
	Long:         `Get AWS EKS kubernetes runtime instances from the system.`,
	SilenceUsage: true,
	PreRun:       commandPreRunFunc,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, _, apiEndpoint, requestedControlPlane := getClientContext(cmd)

		// get AWS EKS kubernetes runtime instances
		awsEksKubernetesRuntimeInstances, err := client.GetAwsEksKubernetesRuntimeInstances(apiClient, apiEndpoint)
		if err != nil {
			cli.Error("failed to retrieve AWS EKS kubernetes runtime instances", err)
			os.Exit(1)
		}

		// write the output
		if len(*awsEksKubernetesRuntimeInstances) == 0 {
			cli.Info(fmt.Sprintf(
				"No AWS EKS kubernetes runtime instances currently managed by %s threeport control plane",
				requestedControlPlane,
			))
			os.Exit(0)
		}
		writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
		fmt.Fprintln(writer, "NAME\t REGION\t KUBERNETES RUNTIME INSTANCE NAME\t AWS EKS KUBERNETES DEFINITION NAME\t RECONCILED\t AGE")
		metadataErr := false
		var kubernetesRuntimeInstanceErr error
		var awsEksKubernetesRuntimeDefinitionErr error
		for _, aekri := range *awsEksKubernetesRuntimeInstances {
			var kubernetesRuntimeInstanceName string
			// get kubernetes runtime instance
			kubernetesRuntimeInstance, err := client.GetKubernetesRuntimeInstanceByID(apiClient, apiEndpoint, *aekri.KubernetesRuntimeInstanceID)
			if err != nil {
				metadataErr = true
				kubernetesRuntimeInstanceErr = err
				kubernetesRuntimeInstanceName = "<error>"
			}
			kubernetesRuntimeInstanceName = *kubernetesRuntimeInstance.Name

			var awsEksKubernetesRuntimeDefinitionName string
			// get AWS EKS kubernetes runtime definition
			awsEksKubernetesRuntimeDefinition, err := client.GetAwsEksKubernetesRuntimeDefinitionByID(apiClient, apiEndpoint, *aekri.AwsEksKubernetesRuntimeDefinitionID)
			if err != nil {
				metadataErr = true
				awsEksKubernetesRuntimeDefinitionErr = err
				awsEksKubernetesRuntimeDefinitionName = "<error>"
			}
			awsEksKubernetesRuntimeDefinitionName = *awsEksKubernetesRuntimeDefinition.Name

			fmt.Fprintln(
				writer, *aekri.Name, "\t", *aekri.Region, "\t", kubernetesRuntimeInstanceName, "\t", awsEksKubernetesRuntimeDefinitionName, "\t",
				*aekri.Reconciled, "\t", util.GetAge(aekri.CreatedAt),
			)
		}
		writer.Flush()

		if metadataErr {
			if kubernetesRuntimeInstanceErr != nil {
				cli.Error("encountered an error retrieving kubernetes runtime instance info", kubernetesRuntimeInstanceErr)
			}
			if awsEksKubernetesRuntimeDefinitionErr != nil {
				cli.Error("encountered an error retrieving AWS EKS kubernetes runtime definition info", awsEksKubernetesRuntimeDefinitionErr)
			}
			os.Exit(1)
		}
	},
}

func init() {
	GetCmd.AddCommand(GetAwsEksKubernetesRuntimeInstancesCmd)
	GetAwsEksKubernetesRuntimeInstancesCmd.Flags().StringVarP(
		&cliArgs.ControlPlaneName,
		"control-plane-name", "i", "", "Optional. Name of control plane. Will default to current control plane if not provided.",
	)
}