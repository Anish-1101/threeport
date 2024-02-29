// originally generated by 'threeport-sdk codegen api-model' but will not be regenerated - intended for modification

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	v0 "github.com/threeport/threeport/pkg/api/v0"
	client "github.com/threeport/threeport/pkg/client/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
)

// outputGetControlPlanesCmd produces the tabular output for the
// 'tptctl get control-planes' command.
func outputGetControlPlanesCmd(
	controlPlaneInstances *[]v0.ControlPlaneInstance,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "NAME\t CONTROL PLANE DEFINITION\t CONTROL PLANE INSTANCE\t AUTH ENABLED\t GENESIS CONTROL PLANE\t KUBERNETES RUNTIME INSTANCE\t AGE")
	metadataErr := false
	var controlPlaneDefinitionErr error
	var kubernetesRuntimeInstErr error
	for _, ci := range *controlPlaneInstances {
		// get control plane definition name for instance
		var controlPlaneDefName string
		var controlPlaneDefAuth string
		controlPlaneDefinition, err := client.GetControlPlaneDefinitionByID(apiClient, apiEndpoint, *ci.ControlPlaneDefinitionID)
		if err != nil {
			metadataErr = true
			controlPlaneDefinitionErr = err
			controlPlaneDefName = "<error>"
			controlPlaneDefAuth = "<error>"
		} else {
			controlPlaneDefName = *controlPlaneDefinition.Name
			controlPlaneDefAuth = strconv.FormatBool(*controlPlaneDefinition.AuthEnabled)
		}
		// get kubernetes runtime instance name for instance
		var kubernetesRuntimeInst string
		kubernetesRuntimeInstance, err := client.GetKubernetesRuntimeInstanceByID(apiClient, apiEndpoint, *ci.KubernetesRuntimeInstanceID)
		if err != nil {
			metadataErr = true
			kubernetesRuntimeInstErr = err
			kubernetesRuntimeInst = "<error>"
		} else {
			kubernetesRuntimeInst = *kubernetesRuntimeInstance.Name
		}

		fmt.Fprintln(
			writer,
			controlPlaneDefName, "\t",
			controlPlaneDefName, "\t",
			*ci.Name, "\t",
			controlPlaneDefAuth, "\t",
			*ci.Genesis, "\t",
			kubernetesRuntimeInst, "\t",
			util.GetAge(ci.CreatedAt),
		)
	}
	writer.Flush()

	if metadataErr {
		multiError := util.MultiError{}
		if controlPlaneDefinitionErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving control plane definition info: %w", controlPlaneDefinitionErr),
			)
		}
		if kubernetesRuntimeInstErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving kubernetes runtime instance info: %w", kubernetesRuntimeInstErr),
			)
		}
		return multiError.Error()
	}

	return nil
}

// outputGetControlPlaneDefinitionsCmd produces the tabular output for the
// 'tptctl get control-plane-definitions' command.
func outputGetControlPlaneDefinitionsCmd(
	controlPlaneDefinitions *[]v0.ControlPlaneDefinition,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "NAME\t AUTH ENABLED\t AGE")
	for _, controlPlaneDefinition := range *controlPlaneDefinitions {
		fmt.Fprintln(
			writer,
			*controlPlaneDefinition.Name, "\t",
			*controlPlaneDefinition.AuthEnabled, "\t",
			util.GetAge(controlPlaneDefinition.CreatedAt),
		)
	}
	writer.Flush()

	return nil
}

// outputGetControlPlaneInstancesCmd produces the tabular output for the
// 'tptctl get control-plane-instances' command.
func outputGetControlPlaneInstancesCmd(
	controlPlaneInstances *[]v0.ControlPlaneInstance,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "NAME\t GENESIS CONTROL PLANE\t CONTROL PLANE DEFINITION\t KUBERNETES RUNTIME INSTANCE\t AGE")
	metadataErr := false
	var controlPlaneDefErr error
	var kubernetesRuntimeInstErr error
	for _, ci := range *controlPlaneInstances {
		// get control plane definition name for instance
		var controlPlaneDef string
		controlPlaneDefinition, err := client.GetControlPlaneDefinitionByID(apiClient, apiEndpoint, *ci.ControlPlaneDefinitionID)
		if err != nil {
			metadataErr = true
			controlPlaneDefErr = err
			controlPlaneDef = "<error>"
		} else {
			controlPlaneDef = *controlPlaneDefinition.Name
		}
		// get kubernetes runtime instance name for instance
		var kubernetesRuntimeInst string
		kubernetesRuntimeInstance, err := client.GetKubernetesRuntimeInstanceByID(apiClient, apiEndpoint, *ci.KubernetesRuntimeInstanceID)
		if err != nil {
			metadataErr = true
			kubernetesRuntimeInstErr = err
			kubernetesRuntimeInst = "<error>"
		} else {
			kubernetesRuntimeInst = *kubernetesRuntimeInstance.Name
		}
		fmt.Fprintln(
			writer,
			*ci.Name, "\t",
			*ci.Genesis, "\t",
			controlPlaneDef, "\t",
			kubernetesRuntimeInst, "\t",
			util.GetAge(ci.CreatedAt),
		)
	}
	writer.Flush()

	if metadataErr {
		multiError := util.MultiError{}
		if controlPlaneDefErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving control plane definition info: %w", controlPlaneDefErr),
			)
		}
		if kubernetesRuntimeInstErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving kubernetes runtime instance info: %w", kubernetesRuntimeInstErr),
			)
		}
		return multiError.Error()
	}

	return nil
}
