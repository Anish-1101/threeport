// generated by 'threeport-codegen api-version' - do not edit

package routes

import (
	echo "github.com/labstack/echo/v4"
	handlers "github.com/threeport/threeport/internal/api/handlers"
)

func AddRoutes(e *echo.Echo, h *handlers.Handler) {
	ProfileRoutes(e, h)
	TierRoutes(e, h)
	AwsAccountRoutes(e, h)
	AwsEksClusterDefinitionRoutes(e, h)
	AwsEksClusterInstanceRoutes(e, h)
	AwsRelationalDatabaseDefinitionRoutes(e, h)
	AwsRelationalDatabaseInstanceRoutes(e, h)
	ClusterDefinitionRoutes(e, h)
	ClusterInstanceRoutes(e, h)
	DomainNameDefinitionRoutes(e, h)
	DomainNameInstanceRoutes(e, h)
	ForwardProxyDefinitionRoutes(e, h)
	ForwardProxyInstanceRoutes(e, h)
	LogBackendRoutes(e, h)
	LogStorageDefinitionRoutes(e, h)
	LogStorageInstanceRoutes(e, h)
	NetworkIngressDefinitionRoutes(e, h)
	NetworkIngressInstanceRoutes(e, h)
	WorkloadDefinitionRoutes(e, h)
	WorkloadResourceDefinitionRoutes(e, h)
	WorkloadInstanceRoutes(e, h)
	WorkloadResourceInstanceRoutes(e, h)

}
