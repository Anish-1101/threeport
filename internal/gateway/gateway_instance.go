package gateway

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"gorm.io/datatypes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	workloadutil "github.com/threeport/threeport/internal/workload/util"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	client "github.com/threeport/threeport/pkg/client/v0"
	controller "github.com/threeport/threeport/pkg/controller/v0"
	util "github.com/threeport/threeport/pkg/util/v0"
)

// gatewayInstanceCreated performs reconciliation when a gateway instance
// has been created.
func gatewayInstanceCreated(
	r *controller.Reconciler,
	gatewayInstance *v0.GatewayInstance,
	log *logr.Logger,
) (int64, error) {

	// ensure attached object reference exists
	err := client.EnsureAttachedObjectReferenceExists(
		r.APIClient,
		r.APIServer,
		reflect.TypeOf(*gatewayInstance).String(),
		gatewayInstance.ID,
		gatewayInstance.WorkloadInstanceID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to ensure attached object reference exists: %w", err)
	}

	// initialize threeport object references
	kubernetesRuntimeInstance, gatewayDefinition, workloadInstance, err := getThreeportObjects(r, gatewayInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to get threeport objects: %w", err)
	}

	// validate threeport state before deploying virtual service
	err = validateThreeportState(r, gatewayDefinition, gatewayInstance, workloadInstance, kubernetesRuntimeInstance, log)
	if err != nil {
		return 0, fmt.Errorf("failed to validate threeport state: %w", err)
	}

	// configure workload resource instances
	workloadResourceInstances, err := configureWorkloadResourceInstances(r, gatewayDefinition, workloadInstance, kubernetesRuntimeInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to configure workload resource instances: %w", err)
	}

	// create workload resource instances
	for _, workloadResourceInstance := range *workloadResourceInstances {
		_, err := client.CreateWorkloadResourceInstance(r.APIClient, r.APIServer, &workloadResourceInstance)
		if err != nil {
			return 0, fmt.Errorf("failed to create workload resource instance: %w", err)
		}
	}

	// trigger a reconciliation of the workload instance
	workloadInstanceReconciled := false
	workloadInstance.Reconciled = &workloadInstanceReconciled
	_, err = client.UpdateWorkloadInstance(r.APIClient, r.APIServer, workloadInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to update workload instance: %w", err)
	}

	// update gateway instance
	gatewayInstanceReconciled := true
	gatewayInstance.Reconciled = &gatewayInstanceReconciled
	_, err = client.UpdateGatewayInstance(r.APIClient, r.APIServer, gatewayInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to update gateway instance: %w", err)
	}

	log.V(1).Info(
		"gateway instance created",
		"gatewayResourceInstanceID", gatewayInstance.ID,
	)

	return 0, nil
}

// gatewayInstanceUpdated performs reconciliation when a gateway instance
// has been updated
func gatewayInstanceUpdated(
	r *controller.Reconciler,
	gatewayInstance *v0.GatewayInstance,
	log *logr.Logger,
) (int64, error) {

	// initialize threeport object references
	kubernetesRuntimeInstance, gatewayDefinition, workloadInstance, err := getThreeportObjects(r, gatewayInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to get threeport objects: %w", err)
	}

	// validate threeport state before deploying virtual service
	err = validateThreeportState(r, gatewayDefinition, gatewayInstance, workloadInstance, kubernetesRuntimeInstance, log)
	if err != nil {
		return 0, fmt.Errorf("failed to validate threeport state: %w", err)
	}

	// configure workload resource instances
	updatedWorkloadResourceInstances, err := configureWorkloadResourceInstances(r, gatewayDefinition, workloadInstance, kubernetesRuntimeInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to configure workload resource instances: %w", err)
	}

	// get workload resource instances
	existingWorkloadResourceInstances, err := client.GetWorkloadResourceInstancesByWorkloadInstanceID(r.APIClient, r.APIServer, *gatewayInstance.WorkloadInstanceID)
	if err != nil {
		return 0, fmt.Errorf("failed to get workload resource instances: %w", err)
	}

	// get gateway instance objects
	gatewayInstanceObjects, err := getGatewayInstanceObjects(r, gatewayInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to get gateway instance objects: %w", err)
	}

	for _, resource := range gatewayInstanceObjects {

		// get workload resource instance for virtual service
		existingWorkloadResourceInstance, err := workloadutil.GetUniqueWorkloadResourceInstance(existingWorkloadResourceInstances, resource)
		if err != nil {
			return 0, fmt.Errorf("failed to get workload resource instance: %w", err)
		}

		// get workload resource instance for virtual service
		updatedWorkloadResourceInstance, err := workloadutil.GetUniqueWorkloadResourceInstance(updatedWorkloadResourceInstances, resource)
		if err != nil {
			return 0, fmt.Errorf("failed to get workload resource instance: %w", err)
		}

		// update the workload resource instance
		workloadResourceInstanceReconciled := false
		existingWorkloadResourceInstance.Reconciled = &workloadResourceInstanceReconciled
		existingWorkloadResourceInstance.JSONDefinition = updatedWorkloadResourceInstance.JSONDefinition
		_, err = client.UpdateWorkloadResourceInstance(r.APIClient, r.APIServer, existingWorkloadResourceInstance)
		if err != nil {
			return 0, fmt.Errorf("failed to update workload resource instance: %w", err)
		}
	}

	// trigger a reconciliation of the workload instance
	workloadInstanceReconciled := false
	workloadInstance.Reconciled = &workloadInstanceReconciled
	_, err = client.UpdateWorkloadInstance(r.APIClient, r.APIServer, workloadInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to update workload instance: %w", err)
	}

	// update gateway instance
	gatewayInstanceReconciled := true
	gatewayInstance.Reconciled = &gatewayInstanceReconciled
	_, err = client.UpdateGatewayInstance(r.APIClient, r.APIServer, gatewayInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to update gateway instance: %w", err)
	}

	log.V(1).Info(
		"gateway instance updated",
		"gatewayResourceInstanceID", gatewayInstance.ID,
	)

	return 0, nil
}

// gatewayInstanceDeleted performs reconciliation when a gateway instance
// has been deleted
func gatewayInstanceDeleted(
	r *controller.Reconciler,
	gatewayInstance *v0.GatewayInstance,
	log *logr.Logger,
) (int64, error) {
	// check that deletion is scheduled - if not there's a problem
	if gatewayInstance.DeletionScheduled == nil {
		return 0, errors.New("deletion notification receieved but not scheduled")
	}

	// check to see if reconciled - it should not be, but if so we should do no
	// more
	if gatewayInstance.DeletionConfirmed != nil {
		return 0, nil
	}

	// get workload resource instances
	workloadResourceInstances, err := client.GetWorkloadResourceInstancesByWorkloadInstanceID(r.APIClient, r.APIServer, *gatewayInstance.WorkloadInstanceID)
	if err != nil {
		if errors.Is(err, client.ErrorObjectNotFound) {
			// workload instance has already been deleted
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get workload resource instances: %w", err)
	}

	// get gateway instance objects
	gatewayInstanceObjects, err := getGatewayInstanceObjects(r, gatewayInstance)
	if err != nil {
		if errors.Is(err, client.ErrorObjectNotFound) {
			// workload instance has already been deleted
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get gateway instance objects: %w", err)
	}

	for _, resource := range gatewayInstanceObjects {

		// get workload resource instance for virtual service
		workloadResourceInstance, err := workloadutil.GetUniqueWorkloadResourceInstance(workloadResourceInstances, resource)
		if err != nil {
			// workload instance has already been deleted
			return 0, nil
		}

		// schedule workload resource instance for deletion
		scheduledForDeletion := time.Now().UTC()
		reconciledWorkloadResourceInstance := false
		workloadResourceInstance = &v0.WorkloadResourceInstance{
			Common:               v0.Common{ID: workloadResourceInstance.ID},
			ScheduledForDeletion: &scheduledForDeletion,
			Reconciled:           &reconciledWorkloadResourceInstance,
		}
		_, err = client.UpdateWorkloadResourceInstance(r.APIClient, r.APIServer, workloadResourceInstance)
		if err != nil {
			if errors.Is(err, client.ErrorObjectNotFound) {
				// workload resource instance has already been deleted
				return 0, nil
			}
			return 0, fmt.Errorf("failed to update workload resource instance: %w", err)
		}
	}

	// trigger a reconciliation of the workload instance
	if gatewayInstance.WorkloadInstanceID == nil {
		return 0, fmt.Errorf("failed to delete workload instance, workloadInstanceID is nil")
	}
	reconciledWorkloadInstance := false
	workloadInstance := &v0.WorkloadInstance{
		Common:         v0.Common{ID: gatewayInstance.WorkloadInstanceID},
		Reconciliation: v0.Reconciliation{Reconciled: &reconciledWorkloadInstance},
	}
	_, err = client.UpdateWorkloadInstance(r.APIClient, r.APIServer, workloadInstance)
	if err != nil {
		if errors.Is(err, client.ErrorObjectNotFound) {
			// workload resource instance has already been deleted
			return 0, nil
		}
		return 0, fmt.Errorf("failed to update workload instance: %w", err)
	}

	// delete the gateway instance that was scheduled for deletion
	deletionReconciled := true
	deletionTimestamp := time.Now().UTC()
	deletedGatewayInstance := v0.GatewayInstance{
		Common: v0.Common{
			ID: gatewayInstance.ID,
		},
		Reconciliation: v0.Reconciliation{
			Reconciled:           &deletionReconciled,
			DeletionAcknowledged: &deletionTimestamp,
			DeletionConfirmed:    &deletionTimestamp,
		},
	}
	_, err = client.UpdateGatewayInstance(
		r.APIClient,
		r.APIServer,
		&deletedGatewayInstance,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to confirm deletion of gateway instance in threeport API: %w", err)
	}
	_, err = client.DeleteGatewayInstance(
		r.APIClient,
		r.APIServer,
		*gatewayInstance.ID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to delete gateway instance in threeport API: %w", err)
	}

	return 0, nil
}

// getThreeportobjects returns all threeport objects required for
// gateway instance reconciliation
func getThreeportObjects(
	r *controller.Reconciler,
	gatewayInstance *v0.GatewayInstance,
) (*v0.KubernetesRuntimeInstance, *v0.GatewayDefinition, *v0.WorkloadInstance, error) {

	// get kubernetes runtime instance
	if gatewayInstance.KubernetesRuntimeInstanceID == nil {
		return nil, nil, nil, errors.New("kubernetes runtime instance ID is nil")
	}
	kubernetesRuntimeInstance, err := client.GetKubernetesRuntimeInstanceByID(r.APIClient, r.APIServer, *gatewayInstance.KubernetesRuntimeInstanceID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get gateway kubernetes runtime instance by ID: %w", err)
	}

	// get gateway definition
	if gatewayInstance.GatewayDefinitionID == nil {
		return nil, nil, nil, fmt.Errorf("gateway definition ID is nil")
	}
	gatewayDefinition, err := client.GetGatewayDefinitionByID(r.APIClient, r.APIServer, *gatewayInstance.GatewayDefinitionID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get gateway controller workload definition: %w", err)
	}

	// get workload instance
	if gatewayInstance.WorkloadInstanceID == nil {
		return nil, nil, nil, fmt.Errorf("workload instance ID is nil")
	}
	workloadInstance, err := client.GetWorkloadInstanceByID(r.APIClient, r.APIServer, *gatewayInstance.WorkloadInstanceID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get workload instance: %w", err)
	}

	return kubernetesRuntimeInstance, gatewayDefinition, workloadInstance, nil
}

// validateThreeportState validates the state of the threeport API
// prior to reconciling a gateway instance
func validateThreeportState(
	r *controller.Reconciler,
	gatewayDefinition *v0.GatewayDefinition,
	gatewayInstance *v0.GatewayInstance,
	workloadInstance *v0.WorkloadInstance,
	kubernetesRuntimeInstance *v0.KubernetesRuntimeInstance,
	log *logr.Logger,
) error {

	// ensure gateway and workload definition are reconciled
	definitionsReconciled, err := confirmDefinitionsReconciled(r, gatewayInstance)
	if err != nil {
		return fmt.Errorf("failed to determine if workload definition is reconciled: %w", err)
	}
	if !definitionsReconciled {
		return errors.New("workload definition not reconciled")
	}

	// ensure workload instance is reconciled
	if gatewayInstance.WorkloadInstanceID == nil {
		return fmt.Errorf("failed to determine if workload instance is reconciled, workload instance ID is nil")
	}
	workloadInstanceReconciled, err := client.ConfirmWorkloadInstanceReconciled(r, *gatewayInstance.WorkloadInstanceID)
	if err != nil {
		return fmt.Errorf("failed to determine if workload instance is reconciled: %w", err)
	}
	if !workloadInstanceReconciled {
		return errors.New("workload instance not reconciled")
	}

	// ensure gateway controller is deployed
	err = confirmGatewayControllerDeployed(r, kubernetesRuntimeInstance, log)
	if err != nil {
		return fmt.Errorf("failed to reconcile gateway controller: %w", err)
	}

	// ensure gateway controller instance is reconciled
	if kubernetesRuntimeInstance.GatewayControllerInstanceID == nil {
		return fmt.Errorf("failed to determine if gateway controller instance is reconciled, gateway controller instance ID is nil")
	}
	gatewayControllerInstanceReconciled, err := client.ConfirmWorkloadInstanceReconciled(r, *kubernetesRuntimeInstance.GatewayControllerInstanceID)
	if err != nil {
		return fmt.Errorf("failed to determine if gateway controller instance is reconciled: %w", err)
	}
	if !gatewayControllerInstanceReconciled {
		return errors.New("gateway controller instance not reconciled")
	}

	// confirm requested port exposed
	err = confirmGatewayPortExposed(r, gatewayInstance, kubernetesRuntimeInstance, gatewayDefinition, log)
	if err != nil {
		return fmt.Errorf("failed to confirm requested port is exposed: %w", err)
	}

	return nil
}

// confirmDefinitionsReconciled confirms the gateway definition related to a
// gateway instance is reconciled.
func confirmDefinitionsReconciled(
	r *controller.Reconciler,
	gatewayInstance *v0.GatewayInstance,
) (bool, error) {

	// get gateway definition
	if gatewayInstance.GatewayDefinitionID == nil {
		return false, fmt.Errorf("failed to get gateway definition from gateway instance, gateway definition ID is nil")
	}
	gatewayDefinition, err := client.GetGatewayDefinitionByID(r.APIClient, r.APIServer, *gatewayInstance.GatewayDefinitionID)
	if err != nil {
		return false, fmt.Errorf("failed to get gateway definition by workload definition ID: %w", err)
	}

	// if the gateway definition is not reconciled, return false
	if gatewayDefinition.Reconciled != nil && !*gatewayDefinition.Reconciled {
		return false, nil
	}

	// get workload definition
	if gatewayDefinition.WorkloadDefinitionID == nil {
		return false, fmt.Errorf("failed to get workload definition from gateway definition, workload definition ID is nil")
	}
	workloadDefinition, err := client.GetWorkloadDefinitionByID(r.APIClient, r.APIServer, *gatewayDefinition.WorkloadDefinitionID)
	if err != nil {
		return false, fmt.Errorf("failed to get workload definition by workload definition ID: %w", err)
	}

	// if the workload definition is not reconciled, return false
	if workloadDefinition.Reconciled != nil && !*workloadDefinition.Reconciled {
		return false, nil
	}

	return true, nil
}

// confirmGatewayControllerDeployed confirms the gateway controller is deployed,
// and if not, deploys it.
func confirmGatewayControllerDeployed(
	r *controller.Reconciler,
	kubernetesRuntimeInstance *v0.KubernetesRuntimeInstance,
	log *logr.Logger,
) error {

	// return if kubernetes runtime instance already has a gateway controller instance
	if kubernetesRuntimeInstance.GatewayControllerInstanceID != nil {
		return nil
	}

	// generate gloo edge manifest
	glooEdge, err := createGlooEdge()
	if err != nil {
		return fmt.Errorf("failed to create gloo edge resource: %w", err)
	}

	// generate support services collection manifest
	supportServicesCollection, err := createSupportServicesCollection()
	if err != nil {
		return fmt.Errorf("failed to create support services collection resource: %w", err)
	}

	// get infra provider
	infraProvider, err := client.GetInfraProviderByKubernetesRuntimeInstanceID(r.APIClient, r.APIServer, kubernetesRuntimeInstance.ID)
	if err != nil {
		return fmt.Errorf("failed to get infra provider: %w", err)
	}

	// generate cert manager manifest based on infra provider
	var certManager string
	switch *infraProvider {
	case v0.KubernetesRuntimeInfraProviderEKS:

		resourceInventory, err := client.GetResourceInventoryByK8sRuntimeInst(r.APIClient, r.APIServer, kubernetesRuntimeInstance.ID)
		if err != nil {
			return fmt.Errorf("failed to get dns management iam role arn: %w", err)
		}

		certManager, err = createCertManager(resourceInventory.Dns01ChallengeRole.RoleArn)
		if err != nil {
			return fmt.Errorf("failed to create cert manager resource: %w", err)
		}

	case v0.KubernetesRuntimeInfraProviderKind:

		certManager, err = createCertManager("")
		if err != nil {
			return fmt.Errorf("failed to create cert manager resource: %w", err)
		}

	default:
		break
	}

	// concatenate gloo edge, support services collection, and cert manager manifests
	manifest := fmt.Sprintf("---\n%s\n---\n%s\n---\n%s\n", supportServicesCollection, certManager, glooEdge)

	// create gateway controller workload definition
	workloadDefName := "gloo-edge"
	glooEdgeWorkloadDefinition := v0.WorkloadDefinition{
		Definition:   v0.Definition{Name: &workloadDefName},
		YAMLDocument: &manifest,
	}

	// create gateway controller workload definition
	createdWorkloadDef, err := client.CreateWorkloadDefinition(r.APIClient, r.APIServer, &glooEdgeWorkloadDefinition)
	if err != nil {
		return fmt.Errorf("failed to create gateway controller workload definition: %w", err)
	}

	// create gateway workload instance
	glooEdgeWorkloadInstance := v0.WorkloadInstance{
		Instance:                    v0.Instance{Name: &workloadDefName},
		KubernetesRuntimeInstanceID: kubernetesRuntimeInstance.ID,
		WorkloadDefinitionID:        createdWorkloadDef.ID,
	}
	createdGlooEdgeWorkloadInstance, err := client.CreateWorkloadInstance(r.APIClient, r.APIServer, &glooEdgeWorkloadInstance)
	if err != nil {
		return fmt.Errorf("failed to create gateway controller workload instance: %w", err)
	}

	// update kubernetes runtime instance with gateway controller instance id
	kubernetesRuntimeInstance.GatewayControllerInstanceID = createdGlooEdgeWorkloadInstance.ID
	_, err = client.UpdateKubernetesRuntimeInstance(r.APIClient, r.APIServer, kubernetesRuntimeInstance)
	if err != nil {
		return fmt.Errorf("failed to update kubernetes runtime instance with gateway controller instance id: %w", err)
	}

	log.V(1).Info(
		"gloo edge deployed",
		"workloadInstanceID", glooEdgeWorkloadInstance.ID,
	)

	return nil
}

// confirmGatewayPortExposed confirms whether the gateway controller has
// exposed the requested port
func confirmGatewayPortExposed(
	r *controller.Reconciler,
	gatewayInstance *v0.GatewayInstance,
	kubernetesRuntimeInstance *v0.KubernetesRuntimeInstance,
	gatewayDefinition *v0.GatewayDefinition,
	log *logr.Logger,
) error {

	// get gateway controller workload resource instances
	if kubernetesRuntimeInstance.GatewayControllerInstanceID == nil {
		return fmt.Errorf("gateway controller instance ID is nil")
	}
	workloadResourceInstances, err := client.GetWorkloadResourceInstancesByWorkloadInstanceID(r.APIClient, r.APIServer, *kubernetesRuntimeInstance.GatewayControllerInstanceID)
	if err != nil {
		return fmt.Errorf("failed to get workload resource instances: %w", err)
	}

	// unmarshal gloo edge custom resource
	gateway, err := workloadutil.UnmarshalUniqueWorkloadResourceInstance(workloadResourceInstances, "GlooEdge")
	if err != nil {
		return fmt.Errorf("failed to unmarshal gloo edge workload resource instance: %w", err)
	}

	// get ports from gloo edge custom resource
	ports, found, err := unstructured.NestedSlice(gateway, "spec", "ports")
	if err != nil || !found {
		return fmt.Errorf("failed to get tcp ports from from gloo edge custom resource: %v", err)
	}

	// check existing gateways for requested port
	var portFound = false
	for _, portSpec := range ports {
		spec := portSpec.(map[string]interface{})
		portNumber, portNumberFound, err := unstructured.NestedFloat64(spec, "port")
		if err != nil {
			return fmt.Errorf("failed to get port from from gloo edge custom resource: %v", err)
		}

		ssl, sslFound, err := unstructured.NestedBool(spec, "ssl")
		if err != nil {
			return fmt.Errorf("failed to get ssl from from gloo edge custom resource: %v", err)
		}

		if portNumberFound &&
			sslFound &&
			ssl == *gatewayDefinition.TLSEnabled &&
			int(portNumber) == *gatewayDefinition.TCPPort {
			portFound = true
			break
		}
	}

	// return if port is found
	if portFound {
		log.V(1).Info(
			"port already exposed",
			"port", fmt.Sprintf("%d", *gatewayDefinition.TCPPort),
		)
		return nil
	}

	// otherwise, update gloo edge configuration

	// create a new gloo edge port object
	portNumber := int64(*gatewayDefinition.TCPPort)
	portString := strconv.Itoa(int(*gatewayDefinition.TCPPort))
	glooEdgePort := createGlooEdgePort(portString, portNumber, *gatewayDefinition.TLSEnabled)

	// append the new port to the ports slice
	ports = append(ports, glooEdgePort)

	// set the ports slice on the gloo edge object
	err = unstructured.SetNestedSlice(gateway, ports, "spec", "ports")
	if err != nil {
		return fmt.Errorf("failed to set ports on gloo edge custom resource: %v", err)
	}

	jsonDefinition, err := util.MarshalJSON(gateway)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	// update the gloo edge workload resource object
	glooEdgeObject, err := workloadutil.GetUniqueWorkloadResourceInstance(workloadResourceInstances, "GlooEdge")
	if err != nil {
		return fmt.Errorf("failed to get gloo edge workload resource instance: %w", err)
	}
	gatewayObjectWorkloadResourceObjectReconciled := false
	glooEdgeObject.Reconciled = &gatewayObjectWorkloadResourceObjectReconciled
	glooEdgeObject.JSONDefinition = &jsonDefinition
	_, err = client.UpdateWorkloadResourceInstance(r.APIClient, r.APIServer, glooEdgeObject)
	if err != nil {
		return fmt.Errorf("failed to update gloo edge workload resource instance: %w", err)
	}

	// trigger a reconciliation of the gateway controller workload instance
	glooEdgeReconciled := false
	updatedGatewayControllerWorkloadInstance := v0.WorkloadInstance{
		Common:         v0.Common{ID: kubernetesRuntimeInstance.GatewayControllerInstanceID},
		Reconciliation: v0.Reconciliation{Reconciled: &glooEdgeReconciled},
	}
	_, err = client.UpdateWorkloadInstance(r.APIClient, r.APIServer, &updatedGatewayControllerWorkloadInstance)
	if err != nil {
		return fmt.Errorf("failed to update gateway controller workload instance: %w", err)
	}

	log.V(1).Info(
		"updated gateway controller instance to expose requested port",
		"port", fmt.Sprintf("%d", *gatewayDefinition.TCPPort),
	)

	return nil

}

// configureVirtualService configures a VirtualService custom resource
// based on the configuration of the gateway workload definition
func configureVirtualService(
	r *controller.Reconciler,
	gatewayDefinition *v0.GatewayDefinition,
	workloadInstance *v0.WorkloadInstance,
	kubernetesRuntimeInstance *v0.KubernetesRuntimeInstance,
) (*datatypes.JSON, error) {

	// get workload resource instances
	workloadResourceInstances, err := client.GetWorkloadResourceInstancesByWorkloadInstanceID(r.APIClient, r.APIServer, *workloadInstance.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workload resource instances: %w", err)
	}

	// unmarshal service
	var service map[string]interface{}
	if gatewayDefinition.ServiceName != nil && *gatewayDefinition.ServiceName != "" {
		service, err = workloadutil.UnmarshalWorkloadResourceInstance(workloadResourceInstances, "Service", *gatewayDefinition.ServiceName)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal service workload resource instance: %w", err)
		}
	} else {
		service, err = workloadutil.UnmarshalUniqueWorkloadResourceInstance(workloadResourceInstances, "Service")
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal service workload resource instance: %w", err)
		}
	}

	// unmarshal service namespace
	namespace, found, err := unstructured.NestedString(service, "metadata", "namespace")
	if err != nil || !found {
		return nil, fmt.Errorf("failed to unmarshal kubernetes service object's namespace field: %w", err)
	}

	// unmarshal service name
	name, found, err := unstructured.NestedString(service, "metadata", "name")
	if err != nil || !found {
		return nil, fmt.Errorf("failed to unmarshal kubernetes service object's name field: %w", err)
	}

	// get gateway workload resource definitions
	gatewayWorkloadResourceDefinitions, err := client.GetWorkloadResourceDefinitionsByWorkloadDefinitionID(r.APIClient, r.APIServer, *gatewayDefinition.WorkloadDefinitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway workload resource definitions: %w", err)
	}

	// unmarshal virtual service
	virtualService, err := workloadutil.UnmarshalUniqueWorkloadResourceDefinition(gatewayWorkloadResourceDefinitions, "VirtualService")
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual service workload resource definition: %w", err)
	}

	// get route array object
	routes, found, err := unstructured.NestedSlice(virtualService, "spec", "virtualHost", "routes")
	if err != nil || !found {
		return nil, fmt.Errorf("failed to get virtualservice route: %w", err)
	}
	if len(routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	// set virtual service upstream name field
	err = unstructured.SetNestedField(
		routes[0].(map[string]interface{}),
		fmt.Sprintf("%s-%s-80", namespace, name), // $namespace-$name-$port is convention for gloo edge upstream names
		"routeAction",
		"single",
		"upstream",
		"name",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set upstream name on virtual service: %w", err)
	}

	// get gloo edge namespace
	glooEdgeNamespace, err := getGlooEdgeNamespace(r, kubernetesRuntimeInstance.GatewayControllerInstanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gloo edge namespace: %w", err)
	}

	// set virtual service upstream namespace field
	err = unstructured.SetNestedField(
		routes[0].(map[string]interface{}),
		glooEdgeNamespace,
		"routeAction",
		"single",
		"upstream",
		"namespace",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set upstream name on virtual service: %w", err)
	}

	// set route field
	err = unstructured.SetNestedSlice(virtualService, routes, "spec", "virtualHost", "routes")
	if err != nil {
		return nil, fmt.Errorf("failed to set route on virtual service: %w", err)
	}

	if *gatewayDefinition.TLSEnabled {

		// set secret ref namespace
		err = unstructured.SetNestedField(
			virtualService,
			namespace,
			"spec",
			"sslConfig",
			"secretRef",
			"namespace",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to set secret ref name on virtual service: %w", err)
		}

	}

	// marshal virtual service
	virtualServiceMarshaled, err := util.MarshalJSON(virtualService)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal virtual service: %w", err)
	}

	return &virtualServiceMarshaled, nil

}

// configureIssuer configures an Issuer custom resource.
func configureIssuer(
	r *controller.Reconciler,
	gatewayDefinition *v0.GatewayDefinition,
	workloadInstance *v0.WorkloadInstance,
	kubernetesRuntimeInstance *v0.KubernetesRuntimeInstance,
	domainNameDefinition v0.DomainNameDefinition,
) (*datatypes.JSON, error) {

	// get gateway workload resource definitions
	gatewayWorkloadResourceDefinitions, err := client.GetWorkloadResourceDefinitionsByWorkloadDefinitionID(r.APIClient, r.APIServer, *gatewayDefinition.WorkloadDefinitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway workload resource definitions: %w", err)
	}

	// unmarshal virtual service
	issuer, err := workloadutil.UnmarshalUniqueWorkloadResourceDefinition(gatewayWorkloadResourceDefinitions, "Issuer")
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual service workload resource definition: %w", err)
	}

	// set issuer name
	kebabDomain := getKebabDomain(*domainNameDefinition.Name)
	err = unstructured.SetNestedField(issuer, kebabDomain, "metadata", "name")
	if err != nil {
		return nil, fmt.Errorf("failed to set name on issuer: %w", err)
	}

	// set solver
	solver := []interface{}{
		map[string]interface{}{
			"selector": map[string]interface{}{
				"dnsZones": []interface{}{
					*domainNameDefinition.Domain,
				},
			},
			"dns01": map[string]interface{}{
				"route53": map[string]interface{}{
					"region": "us-east-1",
				},
			},
		},
	}

	err = unstructured.SetNestedSlice(issuer, solver, "spec", "acme", "solvers")
	if err != nil {
		return nil, fmt.Errorf("failed to set solvers on issuer: %w", err)
	}

	// marshal into json
	issuerMarshaled, err := util.MarshalJSON(issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal issuer: %w", err)
	}

	return &issuerMarshaled, nil
}

// configureCertificate configures a Certificate custom resource.
func configureCertificate(
	r *controller.Reconciler,
	gatewayDefinition *v0.GatewayDefinition,
	workloadInstance *v0.WorkloadInstance,
	kubernetesRuntimeInstance *v0.KubernetesRuntimeInstance,
	domainNameDefinition v0.DomainNameDefinition,
) (*datatypes.JSON, error) {

	// get gateway workload resource definitions
	gatewayWorkloadResourceDefinitions, err := client.GetWorkloadResourceDefinitionsByWorkloadDefinitionID(r.APIClient, r.APIServer, *gatewayDefinition.WorkloadDefinitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway workload resource definitions: %w", err)
	}

	// unmarshal virtual service
	certificate, err := workloadutil.UnmarshalUniqueWorkloadResourceDefinition(gatewayWorkloadResourceDefinitions, "Certificate")
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual service workload resource definition: %w", err)
	}

	// set certificate name
	kebabDomain := getKebabDomain(*domainNameDefinition.Name)
	err = unstructured.SetNestedField(certificate, kebabDomain, "metadata", "name")
	if err != nil {
		return nil, fmt.Errorf("failed to set name on issuer: %w", err)
	}

	dnsNames := []interface{}{*domainNameDefinition.Domain}

	// set dns names
	err = unstructured.SetNestedSlice(certificate, dnsNames, "spec", "dnsNames")
	if err != nil {
		return nil, fmt.Errorf("failed to set dns names on certificate: %w", err)
	}

	// set issuerRef name
	err = unstructured.SetNestedField(certificate, kebabDomain, "spec", "issuerRef", "name")
	if err != nil {
		return nil, fmt.Errorf("failed to set issuerRef on certificate: %w", err)
	}

	// marshal into json
	certificateMarshaled, err := util.MarshalJSON(certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal certificate: %w", err)
	}

	return &certificateMarshaled, nil
}

// getGatewayInstanceObjects returns the objects that should be created.
func getGatewayInstanceObjects(r *controller.Reconciler, gatewayInstance *v0.GatewayInstance) ([]string, error) {

	gatewayDefinition, err := client.GetGatewayDefinitionByID(r.APIClient, r.APIServer, *gatewayInstance.GatewayDefinitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway definition: %w", err)
	}

	if *gatewayDefinition.TLSEnabled {
		return []string{"VirtualService", "Issuer", "Certificate"}, nil
	}
	return []string{"VirtualService"}, nil
}

func configureWorkloadResourceInstances(
	r *controller.Reconciler,
	gatewayDefinition *v0.GatewayDefinition,
	workloadInstance *v0.WorkloadInstance,
	kubernetesRuntimeInstance *v0.KubernetesRuntimeInstance,
) (*[]v0.WorkloadResourceInstance, error) {

	var jsonDefinitions []*datatypes.JSON

	// configure virtual service manifest
	virtualService, err := configureVirtualService(r, gatewayDefinition, workloadInstance, kubernetesRuntimeInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to configure virtual service: %w", err)
	}
	jsonDefinitions = append(jsonDefinitions, virtualService)

	if *gatewayDefinition.TLSEnabled {

		if gatewayDefinition.DomainNameDefinitionID == nil {
			return nil, fmt.Errorf("failed to create certificate, domain name definition ID is nil")
		}

		// get domain name definition
		domainNameDefinition, err := client.GetDomainNameDefinitionByID(r.APIClient, r.APIServer, *gatewayDefinition.DomainNameDefinitionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get domain name definition: %w", err)
		}

		// configure issuer manifest
		issuer, err := configureIssuer(r, gatewayDefinition, workloadInstance, kubernetesRuntimeInstance, *domainNameDefinition)
		if err != nil {
			return nil, fmt.Errorf("failed to configure issuer: %w", err)
		}
		jsonDefinitions = append(jsonDefinitions, issuer)

		// configure certificate manifest
		certificate, err := configureCertificate(r, gatewayDefinition, workloadInstance, kubernetesRuntimeInstance, *domainNameDefinition)
		if err != nil {
			return nil, fmt.Errorf("failed to configure certificate: %w", err)
		}
		jsonDefinitions = append(jsonDefinitions, certificate)
	}

	var workloadResourceInstances []v0.WorkloadResourceInstance
	for _, jsonDefinition := range jsonDefinitions {
		workloadResourceInstanceReconciled := false
		workloadResourceInstance := v0.WorkloadResourceInstance{
			JSONDefinition:     jsonDefinition,
			WorkloadInstanceID: workloadInstance.ID,
			Reconciled:         &workloadResourceInstanceReconciled,
		}
		workloadResourceInstances = append(workloadResourceInstances, workloadResourceInstance)
	}

	return &workloadResourceInstances, nil
}
