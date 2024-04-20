// originally generated by 'threeport-sdk codegen api-model' but will not be regenerated - intended for modification

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	v0 "github.com/threeport/threeport/pkg/api/v0"
	client "github.com/threeport/threeport/pkg/client/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
)

// outputGetGatewaysCmd produces the tabular output for the
// 'tptctl get gateways' command.
func outputGetGatewaysCmd(
	gatewayInstances *[]v0.GatewayInstance,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "NAME\t GATEWAY DEFINITION\t GATEWAY INSTANCE\t AGE")
	var gatewayDefErr error
	for _, gatewayInstance := range *gatewayInstances {
		// get gateway definition name for instance
		var gatewayDef string
		gatewayDefinition, err := client.GetGatewayDefinitionByID(
			apiClient,
			apiEndpoint,
			*gatewayInstance.GatewayDefinitionID,
		)
		if err != nil {
			gatewayDefErr = err
			gatewayDef = "<error>"
		} else {
			gatewayDef = *gatewayDefinition.Name
		}

		fmt.Fprintln(
			writer,
			*gatewayInstance.Name, "\t",
			gatewayDef, "\t",
			*gatewayInstance.Name, "\t",
			util.GetAge(gatewayInstance.CreatedAt),
		)
	}
	writer.Flush()

	if gatewayDefErr != nil {
		return fmt.Errorf("encountered an error retrieving gateway definition info: %w", gatewayDefErr)
	}

	return nil
}

// outputGetGatewayDefinitionsCmd produces the tabular output for the
// 'tptctl get gateway-definitions' command.
func outputGetGatewayDefinitionsCmd(
	gatewayDefinitions *[]v0.GatewayDefinition,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "NAME\t PORTS\t SUBDOMAIN\t KUBERNETES SERVICE NAME\t AGE")
	var gatewayPortsErr error
	for _, g := range *gatewayDefinitions {
		// get gateway ports
		var gwPorts string
		gatewayPorts, err := client.GetGatewayPortsAsString(apiClient, apiEndpoint, *g.Common.ID)
		if err != nil {
			gatewayPortsErr = err
			gwPorts = "<error>"
		} else {
			gwPorts = gatewayPorts
		}
		fmt.Fprintln(
			writer,
			*g.Name, "\t",
			gwPorts, "\t",
			util.DerefString(g.SubDomain), "\t",
			util.DerefString(g.ServiceName), "\t",
			util.GetAge(g.CreatedAt),
		)
	}
	writer.Flush()

	if gatewayPortsErr != nil {
		return fmt.Errorf("encountered an error retrieving gateway ports info: %w", gatewayPortsErr)
	}

	return nil
}

// outputGetGatewayInstancesCmd produces the tabular output for the
// 'tptctl get gateway-instances' command.
func outputGetGatewayInstancesCmd(
	gatewayInstances *[]v0.GatewayInstance,
	apiClient *http.Client,
	apiEndpoint string,
) error {
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
			writer,
			*g.Name, "\t",
			gatewayDef, "\t",
			kubernetesRuntimeInst, "\t",
			workloadInst, "\t",
			util.GetAge(g.CreatedAt),
		)
	}
	writer.Flush()

	if metadataErr {
		multiError := util.MultiError{}
		if gatewayDefErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving gateway definition info: %w", gatewayDefErr),
			)
		}
		if kubernetesRuntimeInstErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving kubernetes runtime instance info: %w", kubernetesRuntimeInstErr),
			)
		}
		if workloadInstErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving workload instance info: %w", workloadInstErr),
			)
		}
		return multiError.Error()
	}

	return nil
}

// outputGetDomainNamesCmd produces the tabular output for the
// 'tptctl get domain-names' command.
func outputGetDomainNamesCmd(
	domainNameInstances *[]v0.DomainNameInstance,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "NAME\t DOMAIN NAME DEFINITION\t DOMAIN NAME INSTANCE\t AGE")
	var domainNameDefErr error
	// asdf
	for _, domainNameInstance := range *domainNameInstances {
		// get domain name definition name for instance
		var domainNameDef string
		domainNameDefinition, err := client.GetDomainNameDefinitionByID(
			apiClient,
			apiEndpoint,
			*domainNameInstance.DomainNameDefinitionID,
		)
		if err != nil {
			domainNameDefErr = err
			domainNameDef = "<error>"
		} else {
			domainNameDef = *domainNameDefinition.Name
		}

		fmt.Fprintln(
			writer,
			*domainNameInstance.Name, "\t",
			domainNameDef, "\t",
			*domainNameInstance.Name, "\t",
			util.GetAge(domainNameInstance.CreatedAt),
		)
	}
	writer.Flush()

	if domainNameDefErr != nil {
		return fmt.Errorf("encountered an error retrieving domain name definition info: %w", domainNameDefErr)
	}

	return nil
}

// outputGetDomainNameDefinitionsCmd produces the tabular output for the
// 'tptctl get domain-name-definitions' command.
func outputGetDomainNameDefinitionsCmd(
	domainNameDefinitions *[]v0.DomainNameDefinition,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "NAME\t ZONE\t ADMIN EMAIL\t AGE ")
	for _, dn := range *domainNameDefinitions {
		fmt.Fprintln(
			writer, *dn.Name,
			"\t", *dn.Zone, "\t",
			*dn.AdminEmail, "\t",
			util.GetAge(dn.CreatedAt),
		)
	}
	writer.Flush()

	return nil
}

// outputGetDomainNameInstancesCmd produces the tabular output for the
// 'tptctl get domain-name-instances' command.
func outputGetDomainNameInstancesCmd(
	domainNameInstances *[]v0.DomainNameInstance,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "NAME\t DOMAIN NAME DEFINITION\t KUBERNETES RUNTIME INSTANCE\t WORKLOAD INSTANCE\t AGE")
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
			writer,
			*d.Name, "\t",
			domainNameDef, "\t",
			kubernetesRuntimeInst, "\t",
			workloadInst, "\t",
			util.GetAge(d.CreatedAt),
		)
	}
	writer.Flush()

	if metadataErr {
		multiError := util.MultiError{}
		if domainNameDefErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving domain name definition info: %w", domainNameDefErr),
			)
		}
		if kubernetesRuntimeInstErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving kubernetes runtime instance info: %w", kubernetesRuntimeInstErr),
			)
		}
		if workloadInstErr != nil {
			multiError.AppendError(
				fmt.Errorf("encountered an error retrieving workload instance info: %w", workloadInstErr),
			)
		}
		return multiError.Error()
	}

	return nil
}
