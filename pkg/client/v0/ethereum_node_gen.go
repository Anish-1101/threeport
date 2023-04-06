// generated by 'threeport-codegen api-model' - do not edit

package v0

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	client "github.com/threeport/threeport/pkg/client"
	"net/http"
)

// GetEthereumNodeDefinitionByID feteches a ethereum node definition by ID
func GetEthereumNodeDefinitionByID(id uint, apiAddr, apiToken string) (*v0.EthereumNodeDefinition, error) {
	var ethereumNodeDefinition v0.EthereumNodeDefinition

	response, err := GetResponse(
		fmt.Sprintf("%s/%s/ethereum-node-definitions/%d", apiAddr, ApiVersion, id),
		apiToken,
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &ethereumNodeDefinition, err
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &ethereumNodeDefinition, err
	}

	if err = json.Unmarshal(jsonData, &ethereumNodeDefinition); err != nil {
		return &ethereumNodeDefinition, err
	}

	return &ethereumNodeDefinition, nil
}

// GetEthereumNodeDefinitionByName feteches a ethereum node definition by name
func GetEthereumNodeDefinitionByName(name, apiAddr, apiToken string) (*v0.EthereumNodeDefinition, error) {
	var ethereumNodeDefinitions []v0.EthereumNodeDefinition

	response, err := GetResponse(
		fmt.Sprintf("%s/%s/ethereum-node-definitions?name=%s", apiAddr, ApiVersion, name),
		apiToken,
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &v0.EthereumNodeDefinition{}, err
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &v0.EthereumNodeDefinition{}, err
	}

	if err = json.Unmarshal(jsonData, &ethereumNodeDefinitions); err != nil {
		return &v0.EthereumNodeDefinition{}, err
	}

	switch {
	case len(ethereumNodeDefinitions) < 1:
		return &v0.EthereumNodeDefinition{}, errors.New(fmt.Sprintf("no workload definitions with name %s", name))
	case len(ethereumNodeDefinitions) > 1:
		return &v0.EthereumNodeDefinition{}, errors.New(fmt.Sprintf("more than one workload definition with name %s returned", name))
	}

	return &ethereumNodeDefinitions[0], nil
}

// CreateEthereumNodeDefinition creates a new ethereum node definition
func CreateEthereumNodeDefinition(ethereumNodeDefinition *v0.EthereumNodeDefinition, apiAddr, apiToken string) (*v0.EthereumNodeDefinition, error) {
	jsonEthereumNodeDefinition, err := client.MarshalObject(ethereumNodeDefinition)
	if err != nil {
		return ethereumNodeDefinition, err
	}

	response, err := GetResponse(
		fmt.Sprintf("%s/%s/ethereum-node-definitions", apiAddr, ApiVersion),
		apiToken,
		http.MethodGet,
		bytes.NewBuffer(jsonEthereumNodeDefinition),
		http.StatusCreated,
	)
	if err != nil {
		return ethereumNodeDefinition, err
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return ethereumNodeDefinition, err
	}

	if err = json.Unmarshal(jsonData, &ethereumNodeDefinition); err != nil {
		return ethereumNodeDefinition, err
	}

	return ethereumNodeDefinition, nil
}

// UpdateEthereumNodeDefinition updates a ethereum node definition
func UpdateEthereumNodeDefinition(ethereumNodeDefinition *v0.EthereumNodeDefinition, apiAddr, apiToken string, id uint) (*v0.EthereumNodeDefinition, error) {
	jsonEthereumNodeDefinition, err := client.MarshalObject(ethereumNodeDefinition)
	if err != nil {
		return ethereumNodeDefinition, err
	}

	response, err := GetResponse(
		fmt.Sprintf("%s/%s/ethereum-node-definitions/%d", apiAddr, ApiVersion, id),
		apiToken,
		http.MethodPatch,
		bytes.NewBuffer(jsonEthereumNodeDefinition),
		http.StatusOK,
	)
	if err != nil {
		return ethereumNodeDefinition, err
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return ethereumNodeDefinition, err
	}

	if err = json.Unmarshal(jsonData, &ethereumNodeDefinition); err != nil {
		return ethereumNodeDefinition, err
	}

	return ethereumNodeDefinition, nil
}

// GetEthereumNodeInstanceByID feteches a ethereum node instance by ID
func GetEthereumNodeInstanceByID(id uint, apiAddr, apiToken string) (*v0.EthereumNodeInstance, error) {
	var ethereumNodeInstance v0.EthereumNodeInstance

	response, err := GetResponse(
		fmt.Sprintf("%s/%s/ethereum-node-instances/%d", apiAddr, ApiVersion, id),
		apiToken,
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &ethereumNodeInstance, err
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return &ethereumNodeInstance, err
	}

	if err = json.Unmarshal(jsonData, &ethereumNodeInstance); err != nil {
		return &ethereumNodeInstance, err
	}

	return &ethereumNodeInstance, nil
}

// GetEthereumNodeInstanceByName feteches a ethereum node instance by name
func GetEthereumNodeInstanceByName(name, apiAddr, apiToken string) (*v0.EthereumNodeInstance, error) {
	var ethereumNodeInstances []v0.EthereumNodeInstance

	response, err := GetResponse(
		fmt.Sprintf("%s/%s/ethereum-node-instances?name=%s", apiAddr, ApiVersion, name),
		apiToken,
		http.MethodGet,
		new(bytes.Buffer),
		http.StatusOK,
	)
	if err != nil {
		return &v0.EthereumNodeInstance{}, err
	}

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return &v0.EthereumNodeInstance{}, err
	}

	if err = json.Unmarshal(jsonData, &ethereumNodeInstances); err != nil {
		return &v0.EthereumNodeInstance{}, err
	}

	switch {
	case len(ethereumNodeInstances) < 1:
		return &v0.EthereumNodeInstance{}, errors.New(fmt.Sprintf("no workload definitions with name %s", name))
	case len(ethereumNodeInstances) > 1:
		return &v0.EthereumNodeInstance{}, errors.New(fmt.Sprintf("more than one workload definition with name %s returned", name))
	}

	return &ethereumNodeInstances[0], nil
}

// CreateEthereumNodeInstance creates a new ethereum node instance
func CreateEthereumNodeInstance(ethereumNodeInstance *v0.EthereumNodeInstance, apiAddr, apiToken string) (*v0.EthereumNodeInstance, error) {
	jsonEthereumNodeInstance, err := client.MarshalObject(ethereumNodeInstance)
	if err != nil {
		return ethereumNodeInstance, err
	}

	response, err := GetResponse(
		fmt.Sprintf("%s/%s/ethereum-node-instances", apiAddr, ApiVersion),
		apiToken,
		http.MethodGet,
		bytes.NewBuffer(jsonEthereumNodeInstance),
		http.StatusCreated,
	)
	if err != nil {
		return ethereumNodeInstance, err
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return ethereumNodeInstance, err
	}

	if err = json.Unmarshal(jsonData, &ethereumNodeInstance); err != nil {
		return ethereumNodeInstance, err
	}

	return ethereumNodeInstance, nil
}

// UpdateEthereumNodeInstance updates a ethereum node instance
func UpdateEthereumNodeInstance(ethereumNodeInstance *v0.EthereumNodeInstance, apiAddr, apiToken string) (*v0.EthereumNodeInstance, error) {
	jsonEthereumNodeInstance, err := client.MarshalObject(ethereumNodeInstance)
	if err != nil {
		return ethereumNodeInstance, err
	}

	response, err := GetResponse(
		fmt.Sprintf("%s/%s/ethereum-node-instances/%d", apiAddr, ApiVersion, *ethereumNodeInstance.ID),
		apiToken,
		http.MethodPatch,
		bytes.NewBuffer(jsonEthereumNodeInstance),
		http.StatusOK,
	)
	if err != nil {
		return ethereumNodeInstance, err
	}

	jsonData, err := json.Marshal(response.Data[0])
	if err != nil {
		return ethereumNodeInstance, err
	}

	if err = json.Unmarshal(jsonData, &ethereumNodeInstance); err != nil {
		return ethereumNodeInstance, err
	}

	return ethereumNodeInstance, nil
}
