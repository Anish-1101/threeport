// generated by 'threeport-codegen api-model' - do not edit

package v0

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	util "github.com/threeport/threeport/internal/util"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	"net/http"
)

// GetAwsAccounts fetches all aws accounts.
// TODO: implement pagination
func GetAwsAccounts(apiClient *http.Client, apiAddr string) (*[]v0.AwsAccount, error) {
	var awsAccounts []v0.AwsAccount

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-accounts", apiAddr, ApiVersion),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsAccounts, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &awsAccounts, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsAccounts); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsAccounts, nil
}

// GetAwsAccountByID fetches a aws account by ID.
func GetAwsAccountByID(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsAccount, error) {
	var awsAccount v0.AwsAccount

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-accounts/%d", apiAddr, ApiVersion, id),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsAccount, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsAccount, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsAccount); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsAccount, nil
}

// GetAwsAccountByName fetches a aws account by name.
func GetAwsAccountByName(apiClient *http.Client, apiAddr, name string) (*v0.AwsAccount, error) {
	var awsAccounts []v0.AwsAccount

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-accounts?name=%s", apiAddr, ApiVersion, name),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &v0.AwsAccount{}, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &v0.AwsAccount{}, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsAccounts); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	switch {
	case len(awsAccounts) < 1:
		return &v0.AwsAccount{}, errors.New(fmt.Sprintf("no workload definitions with name %s", name))
	case len(awsAccounts) > 1:
		return &v0.AwsAccount{}, errors.New(fmt.Sprintf("more than one workload definition with name %s returned", name))
	}

	return &awsAccounts[0], nil
}

// CreateAwsAccount creates a new aws account.
func CreateAwsAccount(apiClient *http.Client, apiAddr string, awsAccount *v0.AwsAccount) (*v0.AwsAccount, error) {
	jsonAwsAccount, err := util.MarshalObject(awsAccount)
	if err != nil {
		return awsAccount, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-accounts", apiAddr, ApiVersion),
		http.MethodPost,
		bytes.NewBuffer(jsonAwsAccount),
		http.StatusCreated,
	)
	if err != nil {
		return awsAccount, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsAccount, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsAccount); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return awsAccount, nil
}

// UpdateAwsAccount updates a aws account.
func UpdateAwsAccount(apiClient *http.Client, apiAddr string, awsAccount *v0.AwsAccount) (*v0.AwsAccount, error) {
	// capture the object ID, make a copy of the object, then remove fields that
	// cannot be updated in the API
	awsAccountID := *awsAccount.ID
	payloadAwsAccount := *awsAccount
	payloadAwsAccount.ID = nil
	payloadAwsAccount.CreatedAt = nil
	payloadAwsAccount.UpdatedAt = nil

	jsonAwsAccount, err := util.MarshalObject(payloadAwsAccount)
	if err != nil {
		return awsAccount, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-accounts/%d", apiAddr, ApiVersion, awsAccountID),
		http.MethodPatch,
		bytes.NewBuffer(jsonAwsAccount),
		http.StatusOK,
	)
	if err != nil {
		return awsAccount, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsAccount, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&payloadAwsAccount); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &payloadAwsAccount, nil
}

// DeleteAwsAccount deletes a aws account by ID.
func DeleteAwsAccount(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsAccount, error) {
	var awsAccount v0.AwsAccount

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-accounts/%d", apiAddr, ApiVersion, id),
		http.MethodDelete,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsAccount, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsAccount, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsAccount); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsAccount, nil
}

// GetAwsEksKubernetesRuntimeDefinitions fetches all aws eks kubernetes runtime definitions.
// TODO: implement pagination
func GetAwsEksKubernetesRuntimeDefinitions(apiClient *http.Client, apiAddr string) (*[]v0.AwsEksKubernetesRuntimeDefinition, error) {
	var awsEksKubernetesRuntimeDefinitions []v0.AwsEksKubernetesRuntimeDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-definitions", apiAddr, ApiVersion),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsEksKubernetesRuntimeDefinitions, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &awsEksKubernetesRuntimeDefinitions, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeDefinitions); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsEksKubernetesRuntimeDefinitions, nil
}

// GetAwsEksKubernetesRuntimeDefinitionByID fetches a aws eks kubernetes runtime definition by ID.
func GetAwsEksKubernetesRuntimeDefinitionByID(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsEksKubernetesRuntimeDefinition, error) {
	var awsEksKubernetesRuntimeDefinition v0.AwsEksKubernetesRuntimeDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-definitions/%d", apiAddr, ApiVersion, id),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsEksKubernetesRuntimeDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsEksKubernetesRuntimeDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsEksKubernetesRuntimeDefinition, nil
}

// GetAwsEksKubernetesRuntimeDefinitionByName fetches a aws eks kubernetes runtime definition by name.
func GetAwsEksKubernetesRuntimeDefinitionByName(apiClient *http.Client, apiAddr, name string) (*v0.AwsEksKubernetesRuntimeDefinition, error) {
	var awsEksKubernetesRuntimeDefinitions []v0.AwsEksKubernetesRuntimeDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-definitions?name=%s", apiAddr, ApiVersion, name),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &v0.AwsEksKubernetesRuntimeDefinition{}, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &v0.AwsEksKubernetesRuntimeDefinition{}, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeDefinitions); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	switch {
	case len(awsEksKubernetesRuntimeDefinitions) < 1:
		return &v0.AwsEksKubernetesRuntimeDefinition{}, errors.New(fmt.Sprintf("no workload definitions with name %s", name))
	case len(awsEksKubernetesRuntimeDefinitions) > 1:
		return &v0.AwsEksKubernetesRuntimeDefinition{}, errors.New(fmt.Sprintf("more than one workload definition with name %s returned", name))
	}

	return &awsEksKubernetesRuntimeDefinitions[0], nil
}

// CreateAwsEksKubernetesRuntimeDefinition creates a new aws eks kubernetes runtime definition.
func CreateAwsEksKubernetesRuntimeDefinition(apiClient *http.Client, apiAddr string, awsEksKubernetesRuntimeDefinition *v0.AwsEksKubernetesRuntimeDefinition) (*v0.AwsEksKubernetesRuntimeDefinition, error) {
	jsonAwsEksKubernetesRuntimeDefinition, err := util.MarshalObject(awsEksKubernetesRuntimeDefinition)
	if err != nil {
		return awsEksKubernetesRuntimeDefinition, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-definitions", apiAddr, ApiVersion),
		http.MethodPost,
		bytes.NewBuffer(jsonAwsEksKubernetesRuntimeDefinition),
		http.StatusCreated,
	)
	if err != nil {
		return awsEksKubernetesRuntimeDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsEksKubernetesRuntimeDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return awsEksKubernetesRuntimeDefinition, nil
}

// UpdateAwsEksKubernetesRuntimeDefinition updates a aws eks kubernetes runtime definition.
func UpdateAwsEksKubernetesRuntimeDefinition(apiClient *http.Client, apiAddr string, awsEksKubernetesRuntimeDefinition *v0.AwsEksKubernetesRuntimeDefinition) (*v0.AwsEksKubernetesRuntimeDefinition, error) {
	// capture the object ID, make a copy of the object, then remove fields that
	// cannot be updated in the API
	awsEksKubernetesRuntimeDefinitionID := *awsEksKubernetesRuntimeDefinition.ID
	payloadAwsEksKubernetesRuntimeDefinition := *awsEksKubernetesRuntimeDefinition
	payloadAwsEksKubernetesRuntimeDefinition.ID = nil
	payloadAwsEksKubernetesRuntimeDefinition.CreatedAt = nil
	payloadAwsEksKubernetesRuntimeDefinition.UpdatedAt = nil

	jsonAwsEksKubernetesRuntimeDefinition, err := util.MarshalObject(payloadAwsEksKubernetesRuntimeDefinition)
	if err != nil {
		return awsEksKubernetesRuntimeDefinition, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-definitions/%d", apiAddr, ApiVersion, awsEksKubernetesRuntimeDefinitionID),
		http.MethodPatch,
		bytes.NewBuffer(jsonAwsEksKubernetesRuntimeDefinition),
		http.StatusOK,
	)
	if err != nil {
		return awsEksKubernetesRuntimeDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsEksKubernetesRuntimeDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&payloadAwsEksKubernetesRuntimeDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &payloadAwsEksKubernetesRuntimeDefinition, nil
}

// DeleteAwsEksKubernetesRuntimeDefinition deletes a aws eks kubernetes runtime definition by ID.
func DeleteAwsEksKubernetesRuntimeDefinition(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsEksKubernetesRuntimeDefinition, error) {
	var awsEksKubernetesRuntimeDefinition v0.AwsEksKubernetesRuntimeDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-definitions/%d", apiAddr, ApiVersion, id),
		http.MethodDelete,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsEksKubernetesRuntimeDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsEksKubernetesRuntimeDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsEksKubernetesRuntimeDefinition, nil
}

// GetAwsEksKubernetesRuntimeInstances fetches all aws eks kubernetes runtime instances.
// TODO: implement pagination
func GetAwsEksKubernetesRuntimeInstances(apiClient *http.Client, apiAddr string) (*[]v0.AwsEksKubernetesRuntimeInstance, error) {
	var awsEksKubernetesRuntimeInstances []v0.AwsEksKubernetesRuntimeInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-instances", apiAddr, ApiVersion),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsEksKubernetesRuntimeInstances, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &awsEksKubernetesRuntimeInstances, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeInstances); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsEksKubernetesRuntimeInstances, nil
}

// GetAwsEksKubernetesRuntimeInstanceByID fetches a aws eks kubernetes runtime instance by ID.
func GetAwsEksKubernetesRuntimeInstanceByID(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsEksKubernetesRuntimeInstance, error) {
	var awsEksKubernetesRuntimeInstance v0.AwsEksKubernetesRuntimeInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-instances/%d", apiAddr, ApiVersion, id),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsEksKubernetesRuntimeInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsEksKubernetesRuntimeInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsEksKubernetesRuntimeInstance, nil
}

// GetAwsEksKubernetesRuntimeInstanceByName fetches a aws eks kubernetes runtime instance by name.
func GetAwsEksKubernetesRuntimeInstanceByName(apiClient *http.Client, apiAddr, name string) (*v0.AwsEksKubernetesRuntimeInstance, error) {
	var awsEksKubernetesRuntimeInstances []v0.AwsEksKubernetesRuntimeInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-instances?name=%s", apiAddr, ApiVersion, name),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &v0.AwsEksKubernetesRuntimeInstance{}, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &v0.AwsEksKubernetesRuntimeInstance{}, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeInstances); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	switch {
	case len(awsEksKubernetesRuntimeInstances) < 1:
		return &v0.AwsEksKubernetesRuntimeInstance{}, errors.New(fmt.Sprintf("no workload definitions with name %s", name))
	case len(awsEksKubernetesRuntimeInstances) > 1:
		return &v0.AwsEksKubernetesRuntimeInstance{}, errors.New(fmt.Sprintf("more than one workload definition with name %s returned", name))
	}

	return &awsEksKubernetesRuntimeInstances[0], nil
}

// CreateAwsEksKubernetesRuntimeInstance creates a new aws eks kubernetes runtime instance.
func CreateAwsEksKubernetesRuntimeInstance(apiClient *http.Client, apiAddr string, awsEksKubernetesRuntimeInstance *v0.AwsEksKubernetesRuntimeInstance) (*v0.AwsEksKubernetesRuntimeInstance, error) {
	jsonAwsEksKubernetesRuntimeInstance, err := util.MarshalObject(awsEksKubernetesRuntimeInstance)
	if err != nil {
		return awsEksKubernetesRuntimeInstance, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-instances", apiAddr, ApiVersion),
		http.MethodPost,
		bytes.NewBuffer(jsonAwsEksKubernetesRuntimeInstance),
		http.StatusCreated,
	)
	if err != nil {
		return awsEksKubernetesRuntimeInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsEksKubernetesRuntimeInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return awsEksKubernetesRuntimeInstance, nil
}

// UpdateAwsEksKubernetesRuntimeInstance updates a aws eks kubernetes runtime instance.
func UpdateAwsEksKubernetesRuntimeInstance(apiClient *http.Client, apiAddr string, awsEksKubernetesRuntimeInstance *v0.AwsEksKubernetesRuntimeInstance) (*v0.AwsEksKubernetesRuntimeInstance, error) {
	// capture the object ID, make a copy of the object, then remove fields that
	// cannot be updated in the API
	awsEksKubernetesRuntimeInstanceID := *awsEksKubernetesRuntimeInstance.ID
	payloadAwsEksKubernetesRuntimeInstance := *awsEksKubernetesRuntimeInstance
	payloadAwsEksKubernetesRuntimeInstance.ID = nil
	payloadAwsEksKubernetesRuntimeInstance.CreatedAt = nil
	payloadAwsEksKubernetesRuntimeInstance.UpdatedAt = nil

	jsonAwsEksKubernetesRuntimeInstance, err := util.MarshalObject(payloadAwsEksKubernetesRuntimeInstance)
	if err != nil {
		return awsEksKubernetesRuntimeInstance, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-instances/%d", apiAddr, ApiVersion, awsEksKubernetesRuntimeInstanceID),
		http.MethodPatch,
		bytes.NewBuffer(jsonAwsEksKubernetesRuntimeInstance),
		http.StatusOK,
	)
	if err != nil {
		return awsEksKubernetesRuntimeInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsEksKubernetesRuntimeInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&payloadAwsEksKubernetesRuntimeInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &payloadAwsEksKubernetesRuntimeInstance, nil
}

// DeleteAwsEksKubernetesRuntimeInstance deletes a aws eks kubernetes runtime instance by ID.
func DeleteAwsEksKubernetesRuntimeInstance(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsEksKubernetesRuntimeInstance, error) {
	var awsEksKubernetesRuntimeInstance v0.AwsEksKubernetesRuntimeInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-eks-kubernetes-runtime-instances/%d", apiAddr, ApiVersion, id),
		http.MethodDelete,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsEksKubernetesRuntimeInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsEksKubernetesRuntimeInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsEksKubernetesRuntimeInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsEksKubernetesRuntimeInstance, nil
}

// GetAwsRelationalDatabaseDefinitions fetches all aws relational database definitions.
// TODO: implement pagination
func GetAwsRelationalDatabaseDefinitions(apiClient *http.Client, apiAddr string) (*[]v0.AwsRelationalDatabaseDefinition, error) {
	var awsRelationalDatabaseDefinitions []v0.AwsRelationalDatabaseDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-definitions", apiAddr, ApiVersion),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsRelationalDatabaseDefinitions, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &awsRelationalDatabaseDefinitions, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseDefinitions); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsRelationalDatabaseDefinitions, nil
}

// GetAwsRelationalDatabaseDefinitionByID fetches a aws relational database definition by ID.
func GetAwsRelationalDatabaseDefinitionByID(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsRelationalDatabaseDefinition, error) {
	var awsRelationalDatabaseDefinition v0.AwsRelationalDatabaseDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-definitions/%d", apiAddr, ApiVersion, id),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsRelationalDatabaseDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsRelationalDatabaseDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsRelationalDatabaseDefinition, nil
}

// GetAwsRelationalDatabaseDefinitionByName fetches a aws relational database definition by name.
func GetAwsRelationalDatabaseDefinitionByName(apiClient *http.Client, apiAddr, name string) (*v0.AwsRelationalDatabaseDefinition, error) {
	var awsRelationalDatabaseDefinitions []v0.AwsRelationalDatabaseDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-definitions?name=%s", apiAddr, ApiVersion, name),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &v0.AwsRelationalDatabaseDefinition{}, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &v0.AwsRelationalDatabaseDefinition{}, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseDefinitions); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	switch {
	case len(awsRelationalDatabaseDefinitions) < 1:
		return &v0.AwsRelationalDatabaseDefinition{}, errors.New(fmt.Sprintf("no workload definitions with name %s", name))
	case len(awsRelationalDatabaseDefinitions) > 1:
		return &v0.AwsRelationalDatabaseDefinition{}, errors.New(fmt.Sprintf("more than one workload definition with name %s returned", name))
	}

	return &awsRelationalDatabaseDefinitions[0], nil
}

// CreateAwsRelationalDatabaseDefinition creates a new aws relational database definition.
func CreateAwsRelationalDatabaseDefinition(apiClient *http.Client, apiAddr string, awsRelationalDatabaseDefinition *v0.AwsRelationalDatabaseDefinition) (*v0.AwsRelationalDatabaseDefinition, error) {
	jsonAwsRelationalDatabaseDefinition, err := util.MarshalObject(awsRelationalDatabaseDefinition)
	if err != nil {
		return awsRelationalDatabaseDefinition, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-definitions", apiAddr, ApiVersion),
		http.MethodPost,
		bytes.NewBuffer(jsonAwsRelationalDatabaseDefinition),
		http.StatusCreated,
	)
	if err != nil {
		return awsRelationalDatabaseDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsRelationalDatabaseDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return awsRelationalDatabaseDefinition, nil
}

// UpdateAwsRelationalDatabaseDefinition updates a aws relational database definition.
func UpdateAwsRelationalDatabaseDefinition(apiClient *http.Client, apiAddr string, awsRelationalDatabaseDefinition *v0.AwsRelationalDatabaseDefinition) (*v0.AwsRelationalDatabaseDefinition, error) {
	// capture the object ID, make a copy of the object, then remove fields that
	// cannot be updated in the API
	awsRelationalDatabaseDefinitionID := *awsRelationalDatabaseDefinition.ID
	payloadAwsRelationalDatabaseDefinition := *awsRelationalDatabaseDefinition
	payloadAwsRelationalDatabaseDefinition.ID = nil
	payloadAwsRelationalDatabaseDefinition.CreatedAt = nil
	payloadAwsRelationalDatabaseDefinition.UpdatedAt = nil

	jsonAwsRelationalDatabaseDefinition, err := util.MarshalObject(payloadAwsRelationalDatabaseDefinition)
	if err != nil {
		return awsRelationalDatabaseDefinition, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-definitions/%d", apiAddr, ApiVersion, awsRelationalDatabaseDefinitionID),
		http.MethodPatch,
		bytes.NewBuffer(jsonAwsRelationalDatabaseDefinition),
		http.StatusOK,
	)
	if err != nil {
		return awsRelationalDatabaseDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsRelationalDatabaseDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&payloadAwsRelationalDatabaseDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &payloadAwsRelationalDatabaseDefinition, nil
}

// DeleteAwsRelationalDatabaseDefinition deletes a aws relational database definition by ID.
func DeleteAwsRelationalDatabaseDefinition(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsRelationalDatabaseDefinition, error) {
	var awsRelationalDatabaseDefinition v0.AwsRelationalDatabaseDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-definitions/%d", apiAddr, ApiVersion, id),
		http.MethodDelete,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsRelationalDatabaseDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsRelationalDatabaseDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsRelationalDatabaseDefinition, nil
}

// GetAwsRelationalDatabaseInstances fetches all aws relational database instances.
// TODO: implement pagination
func GetAwsRelationalDatabaseInstances(apiClient *http.Client, apiAddr string) (*[]v0.AwsRelationalDatabaseInstance, error) {
	var awsRelationalDatabaseInstances []v0.AwsRelationalDatabaseInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-instances", apiAddr, ApiVersion),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsRelationalDatabaseInstances, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &awsRelationalDatabaseInstances, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseInstances); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsRelationalDatabaseInstances, nil
}

// GetAwsRelationalDatabaseInstanceByID fetches a aws relational database instance by ID.
func GetAwsRelationalDatabaseInstanceByID(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsRelationalDatabaseInstance, error) {
	var awsRelationalDatabaseInstance v0.AwsRelationalDatabaseInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-instances/%d", apiAddr, ApiVersion, id),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsRelationalDatabaseInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsRelationalDatabaseInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsRelationalDatabaseInstance, nil
}

// GetAwsRelationalDatabaseInstanceByName fetches a aws relational database instance by name.
func GetAwsRelationalDatabaseInstanceByName(apiClient *http.Client, apiAddr, name string) (*v0.AwsRelationalDatabaseInstance, error) {
	var awsRelationalDatabaseInstances []v0.AwsRelationalDatabaseInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-instances?name=%s", apiAddr, ApiVersion, name),
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &v0.AwsRelationalDatabaseInstance{}, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &v0.AwsRelationalDatabaseInstance{}, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseInstances); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	switch {
	case len(awsRelationalDatabaseInstances) < 1:
		return &v0.AwsRelationalDatabaseInstance{}, errors.New(fmt.Sprintf("no workload definitions with name %s", name))
	case len(awsRelationalDatabaseInstances) > 1:
		return &v0.AwsRelationalDatabaseInstance{}, errors.New(fmt.Sprintf("more than one workload definition with name %s returned", name))
	}

	return &awsRelationalDatabaseInstances[0], nil
}

// CreateAwsRelationalDatabaseInstance creates a new aws relational database instance.
func CreateAwsRelationalDatabaseInstance(apiClient *http.Client, apiAddr string, awsRelationalDatabaseInstance *v0.AwsRelationalDatabaseInstance) (*v0.AwsRelationalDatabaseInstance, error) {
	jsonAwsRelationalDatabaseInstance, err := util.MarshalObject(awsRelationalDatabaseInstance)
	if err != nil {
		return awsRelationalDatabaseInstance, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-instances", apiAddr, ApiVersion),
		http.MethodPost,
		bytes.NewBuffer(jsonAwsRelationalDatabaseInstance),
		http.StatusCreated,
	)
	if err != nil {
		return awsRelationalDatabaseInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsRelationalDatabaseInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return awsRelationalDatabaseInstance, nil
}

// UpdateAwsRelationalDatabaseInstance updates a aws relational database instance.
func UpdateAwsRelationalDatabaseInstance(apiClient *http.Client, apiAddr string, awsRelationalDatabaseInstance *v0.AwsRelationalDatabaseInstance) (*v0.AwsRelationalDatabaseInstance, error) {
	// capture the object ID, make a copy of the object, then remove fields that
	// cannot be updated in the API
	awsRelationalDatabaseInstanceID := *awsRelationalDatabaseInstance.ID
	payloadAwsRelationalDatabaseInstance := *awsRelationalDatabaseInstance
	payloadAwsRelationalDatabaseInstance.ID = nil
	payloadAwsRelationalDatabaseInstance.CreatedAt = nil
	payloadAwsRelationalDatabaseInstance.UpdatedAt = nil

	jsonAwsRelationalDatabaseInstance, err := util.MarshalObject(payloadAwsRelationalDatabaseInstance)
	if err != nil {
		return awsRelationalDatabaseInstance, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-instances/%d", apiAddr, ApiVersion, awsRelationalDatabaseInstanceID),
		http.MethodPatch,
		bytes.NewBuffer(jsonAwsRelationalDatabaseInstance),
		http.StatusOK,
	)
	if err != nil {
		return awsRelationalDatabaseInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return awsRelationalDatabaseInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&payloadAwsRelationalDatabaseInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &payloadAwsRelationalDatabaseInstance, nil
}

// DeleteAwsRelationalDatabaseInstance deletes a aws relational database instance by ID.
func DeleteAwsRelationalDatabaseInstance(apiClient *http.Client, apiAddr string, id uint) (*v0.AwsRelationalDatabaseInstance, error) {
	var awsRelationalDatabaseInstance v0.AwsRelationalDatabaseInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/aws-relational-database-instances/%d", apiAddr, ApiVersion, id),
		http.MethodDelete,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &awsRelationalDatabaseInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &awsRelationalDatabaseInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&awsRelationalDatabaseInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &awsRelationalDatabaseInstance, nil
}
