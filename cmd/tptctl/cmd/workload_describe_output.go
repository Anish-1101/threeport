// originally generated by 'threeport-sdk codegen api-model' but will not be regenerated - intended for modification

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/threeport/threeport/internal/workload/status"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	cli "github.com/threeport/threeport/pkg/cli/v0"
	config "github.com/threeport/threeport/pkg/config/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
)

// outputDescribeWorkloadDefinitionCmd produces the plain description
// output for the 'tptctl describe workload-definition' command
func outputDescribeWorkloadDefinitionCmd(
	workloadDefinition *v0.WorkloadDefinition,
	workloadDefinitionConfig *config.WorkloadDefinitionConfig,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	// describe workload definition
	workloadStatus, err := workloadDefinitionConfig.WorkloadDefinition.Describe(apiClient, apiEndpoint)
	if err != nil {
		return fmt.Errorf("failed to describe workload definition: %w", err)
	}

	// output describe details
	fmt.Printf(
		"* WorkloadDefinition Name: %s\n",
		workloadDefinitionConfig.WorkloadDefinition.Name,
	)
	fmt.Printf(
		"* Created: %s\n",
		*workloadDefinition.CreatedAt,
	)
	fmt.Printf(
		"* Last Modified: %s\n",
		*workloadDefinition.UpdatedAt,
	)
	if len(*workloadStatus.WorkloadInstances) == 0 {
		fmt.Println("* No workload instances currently derived from this definition.")
	} else {
		fmt.Println("* Derived Workload Instances:")
		for _, workloadInst := range *workloadStatus.WorkloadInstances {
			fmt.Printf("  * %s\n", *workloadInst.Name)
		}
	}

	return nil
}

// outputDescribeWorkloadInstanceCmd produces the plain description
// output for the 'tptctl describe workload-instance' command
func outputDescribeWorkloadInstanceCmd(
	workloadInstance *v0.WorkloadInstance,
	workloadInstanceConfig *config.WorkloadInstanceConfig,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	// describe workload instance
	workloadStatus, err := workloadInstanceConfig.WorkloadInstance.Describe(
		apiClient,
		apiEndpoint,
	)
	if err != nil {
		cli.Error("failed to describe workload instance", err)
		os.Exit(1)
	}

	// output describe details
	fmt.Printf(
		"* WorkloadInstance Name: %s\n",
		workloadInstanceConfig.WorkloadInstance.Name,
	)
	fmt.Printf(
		"* Created: %s\n",
		*workloadInstance.CreatedAt,
	)
	fmt.Printf(
		"* Last Modified: %s\n",
		*workloadInstance.UpdatedAt,
	)
	fmt.Printf(
		"* Workload Status: %s\n",
		workloadStatus.Status,
	)
	if workloadStatus.Reason != "" {
		fmt.Printf(
			"* Workload Status Reason: %s\n",
			workloadStatus.Reason,
		)
	}
	if len(workloadStatus.Events) > 0 && workloadStatus.Status != status.WorkloadInstanceStatusHealthy {
		cli.Warning("Failed & Warning Events:")
		writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
		fmt.Fprintln(writer, "TYPE\t REASON\t MESSAGE\t AGE")
		for _, event := range workloadStatus.Events {
			fmt.Fprintln(
				writer, *event.Type, "\t", *event.Reason, "\t", *event.Message, "\t",
				util.GetAge(event.Timestamp),
			)
		}
		writer.Flush()
	}

	return nil
}
