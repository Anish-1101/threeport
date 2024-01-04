// generated by 'threeport-codegen api-model' - do not edit

package v0

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
	"net/http"
)

// GetRadiusWorkloadDefinitions fetches all radius workload definitions.
// TODO: implement pagination
func GetRadiusWorkloadDefinitions(apiClient *http.Client, apiAddr string) (*[]v0.RadiusWorkloadDefinition, error) {
	var radiusWorkloadDefinitions []v0.RadiusWorkloadDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-definitions", apiAddr, ApiVersion),
		http.MethodGet,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &radiusWorkloadDefinitions, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &radiusWorkloadDefinitions, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadDefinitions); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &radiusWorkloadDefinitions, nil
}

// GetRadiusWorkloadDefinitionByID fetches a radius workload definition by ID.
func GetRadiusWorkloadDefinitionByID(apiClient *http.Client, apiAddr string, id uint) (*v0.RadiusWorkloadDefinition, error) {
	var radiusWorkloadDefinition v0.RadiusWorkloadDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-definitions/%d", apiAddr, ApiVersion, id),
		http.MethodGet,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &radiusWorkloadDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &radiusWorkloadDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &radiusWorkloadDefinition, nil
}

// GetRadiusWorkloadDefinitionsByQueryString fetches radius workload definitions by provided query string.
func GetRadiusWorkloadDefinitionsByQueryString(apiClient *http.Client, apiAddr string, queryString string) (*[]v0.RadiusWorkloadDefinition, error) {
	var radiusWorkloadDefinitions []v0.RadiusWorkloadDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-definitions?%s", apiAddr, ApiVersion, queryString),
		http.MethodGet,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &radiusWorkloadDefinitions, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &radiusWorkloadDefinitions, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadDefinitions); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &radiusWorkloadDefinitions, nil
}

// GetRadiusWorkloadDefinitionByName fetches a radius workload definition by name.
func GetRadiusWorkloadDefinitionByName(apiClient *http.Client, apiAddr, name string) (*v0.RadiusWorkloadDefinition, error) {
	var radiusWorkloadDefinitions []v0.RadiusWorkloadDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-definitions?name=%s", apiAddr, ApiVersion, name),
		http.MethodGet,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &v0.RadiusWorkloadDefinition{}, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &v0.RadiusWorkloadDefinition{}, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadDefinitions); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	switch {
	case len(radiusWorkloadDefinitions) < 1:
		return &v0.RadiusWorkloadDefinition{}, errors.New(fmt.Sprintf("no radius workload definition with name %s", name))
	case len(radiusWorkloadDefinitions) > 1:
		return &v0.RadiusWorkloadDefinition{}, errors.New(fmt.Sprintf("more than one radius workload definition with name %s returned", name))
	}

	return &radiusWorkloadDefinitions[0], nil
}

// CreateRadiusWorkloadDefinition creates a new radius workload definition.
func CreateRadiusWorkloadDefinition(apiClient *http.Client, apiAddr string, radiusWorkloadDefinition *v0.RadiusWorkloadDefinition) (*v0.RadiusWorkloadDefinition, error) {
	ReplaceAssociatedObjectsWithNil(radiusWorkloadDefinition)
	jsonRadiusWorkloadDefinition, err := util.MarshalObject(radiusWorkloadDefinition)
	if err != nil {
		return radiusWorkloadDefinition, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-definitions", apiAddr, ApiVersion),
		http.MethodPost,
		bytes.NewBuffer(jsonRadiusWorkloadDefinition),
		map[string]string{},
		http.StatusCreated,
	)
	if err != nil {
		return radiusWorkloadDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return radiusWorkloadDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return radiusWorkloadDefinition, nil
}

// UpdateRadiusWorkloadDefinition updates a radius workload definition.
func UpdateRadiusWorkloadDefinition(apiClient *http.Client, apiAddr string, radiusWorkloadDefinition *v0.RadiusWorkloadDefinition) (*v0.RadiusWorkloadDefinition, error) {
	ReplaceAssociatedObjectsWithNil(radiusWorkloadDefinition)
	// capture the object ID, make a copy of the object, then remove fields that
	// cannot be updated in the API
	radiusWorkloadDefinitionID := *radiusWorkloadDefinition.ID
	payloadRadiusWorkloadDefinition := *radiusWorkloadDefinition
	payloadRadiusWorkloadDefinition.ID = nil
	payloadRadiusWorkloadDefinition.CreatedAt = nil
	payloadRadiusWorkloadDefinition.UpdatedAt = nil

	jsonRadiusWorkloadDefinition, err := util.MarshalObject(payloadRadiusWorkloadDefinition)
	if err != nil {
		return radiusWorkloadDefinition, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-definitions/%d", apiAddr, ApiVersion, radiusWorkloadDefinitionID),
		http.MethodPatch,
		bytes.NewBuffer(jsonRadiusWorkloadDefinition),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return radiusWorkloadDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return radiusWorkloadDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&payloadRadiusWorkloadDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	payloadRadiusWorkloadDefinition.ID = &radiusWorkloadDefinitionID
	return &payloadRadiusWorkloadDefinition, nil
}

// DeleteRadiusWorkloadDefinition deletes a radius workload definition by ID.
func DeleteRadiusWorkloadDefinition(apiClient *http.Client, apiAddr string, id uint) (*v0.RadiusWorkloadDefinition, error) {
	var radiusWorkloadDefinition v0.RadiusWorkloadDefinition

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-definitions/%d", apiAddr, ApiVersion, id),
		http.MethodDelete,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &radiusWorkloadDefinition, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &radiusWorkloadDefinition, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadDefinition); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &radiusWorkloadDefinition, nil
}

// GetRadiusWorkloadInstances fetches all radius workload instances.
// TODO: implement pagination
func GetRadiusWorkloadInstances(apiClient *http.Client, apiAddr string) (*[]v0.RadiusWorkloadInstance, error) {
	var radiusWorkloadInstances []v0.RadiusWorkloadInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-instances", apiAddr, ApiVersion),
		http.MethodGet,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &radiusWorkloadInstances, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &radiusWorkloadInstances, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadInstances); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &radiusWorkloadInstances, nil
}

// GetRadiusWorkloadInstanceByID fetches a radius workload instance by ID.
func GetRadiusWorkloadInstanceByID(apiClient *http.Client, apiAddr string, id uint) (*v0.RadiusWorkloadInstance, error) {
	var radiusWorkloadInstance v0.RadiusWorkloadInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-instances/%d", apiAddr, ApiVersion, id),
		http.MethodGet,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &radiusWorkloadInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &radiusWorkloadInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &radiusWorkloadInstance, nil
}

// GetRadiusWorkloadInstancesByQueryString fetches radius workload instances by provided query string.
func GetRadiusWorkloadInstancesByQueryString(apiClient *http.Client, apiAddr string, queryString string) (*[]v0.RadiusWorkloadInstance, error) {
	var radiusWorkloadInstances []v0.RadiusWorkloadInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-instances?%s", apiAddr, ApiVersion, queryString),
		http.MethodGet,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &radiusWorkloadInstances, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &radiusWorkloadInstances, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadInstances); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &radiusWorkloadInstances, nil
}

// GetRadiusWorkloadInstanceByName fetches a radius workload instance by name.
func GetRadiusWorkloadInstanceByName(apiClient *http.Client, apiAddr, name string) (*v0.RadiusWorkloadInstance, error) {
	var radiusWorkloadInstances []v0.RadiusWorkloadInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-instances?name=%s", apiAddr, ApiVersion, name),
		http.MethodGet,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &v0.RadiusWorkloadInstance{}, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &v0.RadiusWorkloadInstance{}, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadInstances); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	switch {
	case len(radiusWorkloadInstances) < 1:
		return &v0.RadiusWorkloadInstance{}, errors.New(fmt.Sprintf("no radius workload instance with name %s", name))
	case len(radiusWorkloadInstances) > 1:
		return &v0.RadiusWorkloadInstance{}, errors.New(fmt.Sprintf("more than one radius workload instance with name %s returned", name))
	}

	return &radiusWorkloadInstances[0], nil
}

// CreateRadiusWorkloadInstance creates a new radius workload instance.
func CreateRadiusWorkloadInstance(apiClient *http.Client, apiAddr string, radiusWorkloadInstance *v0.RadiusWorkloadInstance) (*v0.RadiusWorkloadInstance, error) {
	ReplaceAssociatedObjectsWithNil(radiusWorkloadInstance)
	jsonRadiusWorkloadInstance, err := util.MarshalObject(radiusWorkloadInstance)
	if err != nil {
		return radiusWorkloadInstance, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-instances", apiAddr, ApiVersion),
		http.MethodPost,
		bytes.NewBuffer(jsonRadiusWorkloadInstance),
		map[string]string{},
		http.StatusCreated,
	)
	if err != nil {
		return radiusWorkloadInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return radiusWorkloadInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return radiusWorkloadInstance, nil
}

// UpdateRadiusWorkloadInstance updates a radius workload instance.
func UpdateRadiusWorkloadInstance(apiClient *http.Client, apiAddr string, radiusWorkloadInstance *v0.RadiusWorkloadInstance) (*v0.RadiusWorkloadInstance, error) {
	ReplaceAssociatedObjectsWithNil(radiusWorkloadInstance)
	// capture the object ID, make a copy of the object, then remove fields that
	// cannot be updated in the API
	radiusWorkloadInstanceID := *radiusWorkloadInstance.ID
	payloadRadiusWorkloadInstance := *radiusWorkloadInstance
	payloadRadiusWorkloadInstance.ID = nil
	payloadRadiusWorkloadInstance.CreatedAt = nil
	payloadRadiusWorkloadInstance.UpdatedAt = nil

	jsonRadiusWorkloadInstance, err := util.MarshalObject(payloadRadiusWorkloadInstance)
	if err != nil {
		return radiusWorkloadInstance, fmt.Errorf("failed to marshal provided object to JSON: %w", err)
	}

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-instances/%d", apiAddr, ApiVersion, radiusWorkloadInstanceID),
		http.MethodPatch,
		bytes.NewBuffer(jsonRadiusWorkloadInstance),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return radiusWorkloadInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return radiusWorkloadInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&payloadRadiusWorkloadInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	payloadRadiusWorkloadInstance.ID = &radiusWorkloadInstanceID
	return &payloadRadiusWorkloadInstance, nil
}

// DeleteRadiusWorkloadInstance deletes a radius workload instance by ID.
func DeleteRadiusWorkloadInstance(apiClient *http.Client, apiAddr string, id uint) (*v0.RadiusWorkloadInstance, error) {
	var radiusWorkloadInstance v0.RadiusWorkloadInstance

	response, err := GetResponse(
		apiClient,
		fmt.Sprintf("%s/%s/radius-workload-instances/%d", apiAddr, ApiVersion, id),
		http.MethodDelete,
		new(bytes.Buffer),
		map[string]string{},
		http.StatusOK,
	)
	if err != nil {
		return &radiusWorkloadInstance, fmt.Errorf("call to threeport API returned unexpected response: %w", err)
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &radiusWorkloadInstance, fmt.Errorf("failed to marshal response data from threeport API: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&radiusWorkloadInstance); err != nil {
		return nil, fmt.Errorf("failed to decode object in response data from threeport API: %w", err)
	}

	return &radiusWorkloadInstance, nil
}