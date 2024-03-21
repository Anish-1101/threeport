// generated by 'threeport-sdk gen api-model' - do not edit

package routes

import (
	echo "github.com/labstack/echo/v4"
	handlers "github.com/threeport/threeport/pkg/api-server/v0/handlers"
	v0 "github.com/threeport/threeport/pkg/api/v0"
)

// AttachedObjectReferenceRoutes sets up all routes for the AttachedObjectReference handlers.
func AttachedObjectReferenceRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET("/attached-object-references/versions", h.GetAttachedObjectReferenceVersions)

	e.POST(v0.PathAttachedObjectReferences, h.AddAttachedObjectReference)
	e.GET(v0.PathAttachedObjectReferences, h.GetAttachedObjectReferences)
	e.GET(v0.PathAttachedObjectReferences+"/:id", h.GetAttachedObjectReference)
	e.PATCH(v0.PathAttachedObjectReferences+"/:id", h.UpdateAttachedObjectReference)
	e.PUT(v0.PathAttachedObjectReferences+"/:id", h.ReplaceAttachedObjectReference)
	e.DELETE(v0.PathAttachedObjectReferences+"/:id", h.DeleteAttachedObjectReference)
}

// WorkloadDefinitionRoutes sets up all routes for the WorkloadDefinition handlers.
func WorkloadDefinitionRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET("/workload-definitions/versions", h.GetWorkloadDefinitionVersions)

	e.POST(v0.PathWorkloadDefinitions, h.AddWorkloadDefinition)
	e.GET(v0.PathWorkloadDefinitions, h.GetWorkloadDefinitions)
	e.GET(v0.PathWorkloadDefinitions+"/:id", h.GetWorkloadDefinition)
	e.PATCH(v0.PathWorkloadDefinitions+"/:id", h.UpdateWorkloadDefinition)
	e.PUT(v0.PathWorkloadDefinitions+"/:id", h.ReplaceWorkloadDefinition)
	e.DELETE(v0.PathWorkloadDefinitions+"/:id", h.DeleteWorkloadDefinition)
}

// WorkloadEventRoutes sets up all routes for the WorkloadEvent handlers.
func WorkloadEventRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET("/workload-events/versions", h.GetWorkloadEventVersions)

	e.POST(v0.PathWorkloadEvents, h.AddWorkloadEvent)
	e.GET(v0.PathWorkloadEvents, h.GetWorkloadEvents)
	e.GET(v0.PathWorkloadEvents+"/:id", h.GetWorkloadEvent)
	e.PATCH(v0.PathWorkloadEvents+"/:id", h.UpdateWorkloadEvent)
	e.PUT(v0.PathWorkloadEvents+"/:id", h.ReplaceWorkloadEvent)
	e.DELETE(v0.PathWorkloadEvents+"/:id", h.DeleteWorkloadEvent)
}

// WorkloadInstanceRoutes sets up all routes for the WorkloadInstance handlers.
func WorkloadInstanceRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET("/workload-instances/versions", h.GetWorkloadInstanceVersions)

	e.POST(v0.PathWorkloadInstances, h.AddWorkloadInstance)
	e.GET(v0.PathWorkloadInstances, h.GetWorkloadInstances)
	e.GET(v0.PathWorkloadInstances+"/:id", h.GetWorkloadInstance)
	e.PATCH(v0.PathWorkloadInstances+"/:id", h.UpdateWorkloadInstance)
	e.PUT(v0.PathWorkloadInstances+"/:id", h.ReplaceWorkloadInstance)
	e.DELETE(v0.PathWorkloadInstances+"/:id", h.DeleteWorkloadInstance)
}

// WorkloadResourceDefinitionRoutes sets up all routes for the WorkloadResourceDefinition handlers.
func WorkloadResourceDefinitionRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET("/workload-resource-definitions/versions", h.GetWorkloadResourceDefinitionVersions)

	e.POST(v0.PathWorkloadResourceDefinitions, h.AddWorkloadResourceDefinition)
	e.GET(v0.PathWorkloadResourceDefinitions, h.GetWorkloadResourceDefinitions)
	e.GET(v0.PathWorkloadResourceDefinitions+"/:id", h.GetWorkloadResourceDefinition)
	e.PATCH(v0.PathWorkloadResourceDefinitions+"/:id", h.UpdateWorkloadResourceDefinition)
	e.PUT(v0.PathWorkloadResourceDefinitions+"/:id", h.ReplaceWorkloadResourceDefinition)
	e.DELETE(v0.PathWorkloadResourceDefinitions+"/:id", h.DeleteWorkloadResourceDefinition)
}

// WorkloadResourceInstanceRoutes sets up all routes for the WorkloadResourceInstance handlers.
func WorkloadResourceInstanceRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET("/workload-resource-instances/versions", h.GetWorkloadResourceInstanceVersions)

	e.POST(v0.PathWorkloadResourceInstances, h.AddWorkloadResourceInstance)
	e.GET(v0.PathWorkloadResourceInstances, h.GetWorkloadResourceInstances)
	e.GET(v0.PathWorkloadResourceInstances+"/:id", h.GetWorkloadResourceInstance)
	e.PATCH(v0.PathWorkloadResourceInstances+"/:id", h.UpdateWorkloadResourceInstance)
	e.PUT(v0.PathWorkloadResourceInstances+"/:id", h.ReplaceWorkloadResourceInstance)
	e.DELETE(v0.PathWorkloadResourceInstances+"/:id", h.DeleteWorkloadResourceInstance)
}
