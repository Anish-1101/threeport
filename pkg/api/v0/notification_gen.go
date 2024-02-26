// generated by 'threeport-sdk gen api-version' - do not edit
// +threeport-sdk route-exclude
// +threeport-sdk database-exclude

package v0

import "errors"

// GetSubjectByReconcilerName returns the subject for a reconciler's name.
func GetSubjectByReconcilerName(name string) (string, error) {
	switch name {
	case "AwsEksKubernetesRuntimeInstanceReconciler":
		return AwsEksKubernetesRuntimeInstanceSubject, nil
	case "AwsObjectStorageBucketInstanceReconciler":
		return AwsObjectStorageBucketInstanceSubject, nil
	case "AwsRelationalDatabaseInstanceReconciler":
		return AwsRelationalDatabaseInstanceSubject, nil
	case "ControlPlaneDefinitionReconciler":
		return ControlPlaneDefinitionSubject, nil
	case "ControlPlaneInstanceReconciler":
		return ControlPlaneInstanceSubject, nil
	case "DomainNameInstanceReconciler":
		return DomainNameInstanceSubject, nil
	case "GatewayDefinitionReconciler":
		return GatewayDefinitionSubject, nil
	case "GatewayInstanceReconciler":
		return GatewayInstanceSubject, nil
	case "HelmWorkloadDefinitionReconciler":
		return HelmWorkloadDefinitionSubject, nil
	case "HelmWorkloadInstanceReconciler":
		return HelmWorkloadInstanceSubject, nil
	case "KubernetesRuntimeDefinitionReconciler":
		return KubernetesRuntimeDefinitionSubject, nil
	case "KubernetesRuntimeInstanceReconciler":
		return KubernetesRuntimeInstanceSubject, nil
	case "LoggingDefinitionReconciler":
		return LoggingDefinitionSubject, nil
	case "LoggingInstanceReconciler":
		return LoggingInstanceSubject, nil
	case "MetricsDefinitionReconciler":
		return MetricsDefinitionSubject, nil
	case "MetricsInstanceReconciler":
		return MetricsInstanceSubject, nil
	case "ObservabilityDashboardDefinitionReconciler":
		return ObservabilityDashboardDefinitionSubject, nil
	case "ObservabilityDashboardInstanceReconciler":
		return ObservabilityDashboardInstanceSubject, nil
	case "ObservabilityStackDefinitionReconciler":
		return ObservabilityStackDefinitionSubject, nil
	case "ObservabilityStackInstanceReconciler":
		return ObservabilityStackInstanceSubject, nil
	case "TerraformDefinitionReconciler":
		return TerraformDefinitionSubject, nil
	case "TerraformInstanceReconciler":
		return TerraformInstanceSubject, nil
	case "WorkloadDefinitionReconciler":
		return WorkloadDefinitionSubject, nil
	case "WorkloadInstanceReconciler":
		return WorkloadInstanceSubject, nil

	default:
		return "", errors.New("unrecognized reconciler name")
	}

}
