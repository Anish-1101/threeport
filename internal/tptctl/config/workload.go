package config

import (
	"encoding/json"
	"io/ioutil"

	v0 "github.com/threeport/threeport/pkg/api/v0"
	client "github.com/threeport/threeport/pkg/client/v0"

	"github.com/threeport/threeport/internal/tptctl/install"
)

// WorkloadConfig contains the attributes needed to manage a workload.
type WorkloadConfig struct {
	Name               string                   `yaml:"Name"`
	WorkloadDefinition WorkloadDefinitionConfig `yaml:"WorkloadDefinition"`
	WorkloadInstance   WorkloadInstanceConfig   `yaml:"WorkloadInstance"`
	//WorkloadServiceDependency WorkloadServiceDependencyConfig `yaml:"WorkloadServiceDependency"`
}

// WorkloadDefinitionConfig contains the attributes needed to manage a workload
// definition.
type WorkloadDefinitionConfig struct {
	Name         string `yaml:"Name"`
	YAMLDocument string `yaml:"YAMLDocument"`
	UserID       uint   `yaml:"UserID"`
}

// WorkloadInstanceConfig contains the attributes needed to manage a workload
// instance.
type WorkloadInstanceConfig struct {
	Name                   string `yaml:"Name"`
	WorkloadClusterName    string `yaml:"WorkloadClusterName"`
	WorkloadDefinitionName string `yaml:"WorkloadDefinitionName"`
}

// Create creates a workload in the Threeport API.
func (wc *WorkloadConfig) Create() error {
	// create the definition
	_, aerr := wc.WorkloadDefinition.Create()
	if aerr != nil {
		return aerr
	}

	// create the instance
	_, berr := wc.WorkloadInstance.Create()
	if berr != nil {
		return berr
	}

	//// create the service dependency
	//_, cerr := wc.WorkloadServiceDependency.Create()
	//if cerr != nil {
	//	return cerr
	//}

	return nil
}

// Create creates a workload definition in the Threeport API.
func (wdc *WorkloadDefinitionConfig) Create() (*v0.WorkloadDefinition, error) {
	// get the content of the yaml document
	definitionContent, err := ioutil.ReadFile(wdc.YAMLDocument)
	if err != nil {
		return nil, err
	}
	stringContent := string(definitionContent)

	// construct workload definition object
	workloadDefinition := &v0.WorkloadDefinition{
		Definition: v0.Definition{
			Name:   &wdc.Name,
			UserID: &wdc.UserID,
		},
		YAMLDocument: &stringContent,
	}

	// create workload definition in API
	wdJSON, err := json.Marshal(&workloadDefinition)
	if err != nil {
		return nil, err
	}
	wd, err := client.CreateWorkloadDefinition(wdJSON, install.GetThreeportAPIEndpoint(), "")
	if err != nil {
		return nil, err
	}

	return wd, nil
}

// Create creates a workload instance in the Threeport API.
func (wic *WorkloadInstanceConfig) Create() (*v0.WorkloadInstance, error) {
	// get workload cluster by name
	workloadCluster, err := client.GetClusterInstanceByName(
		wic.WorkloadClusterName,
		install.GetThreeportAPIEndpoint(), "",
	)
	if err != nil {
		return nil, err
	}

	// get workload definition by name
	workloadDefinition, err := client.GetWorkloadDefinitionByName(
		wic.WorkloadDefinitionName,
		install.GetThreeportAPIEndpoint(), "",
	)
	if err != nil {
		return nil, err
	}

	// construct workload instance object
	workloadInstance := &v0.WorkloadInstance{
		Instance: v0.Instance{
			Name: &wic.Name,
		},
		ClusterInstanceID:    workloadCluster.ID,
		WorkloadDefinitionID: workloadDefinition.ID,
	}

	// create workload instance in API
	wiJSON, err := json.Marshal(&workloadInstance)
	if err != nil {
		return nil, err
	}
	wi, err := client.CreateWorkloadInstance(wiJSON, install.GetThreeportAPIEndpoint(), "")
	if err != nil {
		return nil, err
	}

	return wi, nil
}