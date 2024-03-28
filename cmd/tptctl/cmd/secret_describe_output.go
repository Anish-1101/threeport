// originally generated by 'threeport-sdk codegen api-model' but will not be regenerated - intended for modification

package cmd

import (
	"fmt"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	config "github.com/threeport/threeport/pkg/config/v0"
	"net/http"
)

// outputDescribeSecretDefinitionCmd produces the plain description
// output for the 'tptctl describe secret-definition' command
func outputDescribeSecretDefinitionCmd(
	secretDefinition *v0.SecretDefinition,
	secretDefinitionConfig *config.SecretDefinitionConfig,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	// output describe details
	fmt.Printf(
		"* SecretDefinition Name: %s\n",
		secretDefinitionConfig.SecretDefinition.Name,
	)
	fmt.Printf(
		"* Created: %s\n",
		*secretDefinition.CreatedAt,
	)
	fmt.Printf(
		"* Last Modified: %s\n",
		*secretDefinition.UpdatedAt,
	)

	return nil
}

// outputDescribeSecretInstanceCmd produces the plain description
// output for the 'tptctl describe secret-instance' command
func outputDescribeSecretInstanceCmd(
	secretInstance *v0.SecretInstance,
	secretInstanceConfig *config.SecretInstanceConfig,
	apiClient *http.Client,
	apiEndpoint string,
) error {
	// output describe details
	fmt.Printf(
		"* SecretInstance Name: %s\n",
		secretInstanceConfig.SecretInstance.Name,
	)
	fmt.Printf(
		"* Created: %s\n",
		*secretInstance.CreatedAt,
	)
	fmt.Printf(
		"* Last Modified: %s\n",
		*secretInstance.UpdatedAt,
	)

	return nil
}