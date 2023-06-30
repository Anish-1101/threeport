// generated by 'threeport-codegen api-model' - do not edit

package handlers

import (
	"errors"
	echo "github.com/labstack/echo/v4"
	iapi "github.com/threeport/threeport/internal/api"
	api "github.com/threeport/threeport/pkg/api"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	notifications "github.com/threeport/threeport/pkg/notifications/v0"
	gorm "gorm.io/gorm"
	"net/http"
)

///////////////////////////////////////////////////////////////////////////////
// NetworkIngressDefinition
///////////////////////////////////////////////////////////////////////////////

// @Summary GetNetworkIngressDefinitionVersions gets the supported versions for the network ingress definition API.
// @Description Get the supported API versions for network ingress definitions.
// @ID networkIngressDefinition-get-versions
// @Produce json
// @Success 200 {object} api.RESTAPIVersions "OK"
// @Router /network-ingress-definitions/versions [get]
func (h Handler) GetNetworkIngressDefinitionVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeNetworkIngressDefinition)])
}

// @Summary adds a new network ingress definition.
// @Description Add a new network ingress definition to the Threeport database.
// @ID add-networkIngressDefinition
// @Accept json
// @Produce json
// @Param networkIngressDefinition body v0.NetworkIngressDefinition true "NetworkIngressDefinition object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-definitions [post]
func (h Handler) AddNetworkIngressDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressDefinition
	var networkIngressDefinition v0.NetworkIngressDefinition

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&networkIngressDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, networkIngressDefinition, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// check for duplicate names
	var existingNetworkIngressDefinition v0.NetworkIngressDefinition
	nameUsed := true
	result := h.DB.Where("name = ?", networkIngressDefinition.Name).First(&existingNetworkIngressDefinition)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			nameUsed = false
		} else {
			return iapi.ResponseStatus500(c, nil, result.Error, objectType)
		}
	}
	if nameUsed {
		return iapi.ResponseStatus409(c, nil, errors.New("object with provided name already exists"), objectType)
	}

	// persist to DB
	if result := h.DB.Create(&networkIngressDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller
	notifPayload, err := networkIngressDefinition.NotificationPayload(
		notifications.NotificationOperationCreated,
		false,
		0,
	)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}
	h.JS.Publish(v0.NetworkIngressDefinitionCreateSubject, *notifPayload)

	response, err := v0.CreateResponse(nil, networkIngressDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all network ingress definitions.
// @Description Get all network ingress definitions from the Threeport database.
// @ID get-networkIngressDefinitions
// @Accept json
// @Produce json
// @Param name query string false "network ingress definition search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-definitions [get]
func (h Handler) GetNetworkIngressDefinitions(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressDefinition
	params, err := c.(*iapi.CustomContext).GetPaginationParams()
	if err != nil {
		return iapi.ResponseStatus400(c, &params, err, objectType)
	}

	var filter v0.NetworkIngressDefinition
	if err := c.Bind(&filter); err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	var totalCount int64
	if result := h.DB.Model(&v0.NetworkIngressDefinition{}).Where(&filter).Count(&totalCount); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	records := &[]v0.NetworkIngressDefinition{}
	if result := h.DB.Order("ID asc").Where(&filter).Limit(params.Size).Offset((params.Page - 1) * params.Size).Find(records); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a network ingress definition.
// @Description Get a particular network ingress definition from the database.
// @ID get-networkIngressDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-definitions/{id} [get]
func (h Handler) GetNetworkIngressDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressDefinition
	networkIngressDefinitionID := c.Param("id")
	var networkIngressDefinition v0.NetworkIngressDefinition
	if result := h.DB.First(&networkIngressDefinition, networkIngressDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, networkIngressDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates specific fields for an existing network ingress definition.
// @Description Update a network ingress definition in the database.  Provide one or more fields to update.
// @Description Note: This API endpint is for updating network ingress definition objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID update-networkIngressDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param networkIngressDefinition body v0.NetworkIngressDefinition true "NetworkIngressDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-definitions/{id} [patch]
func (h Handler) UpdateNetworkIngressDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressDefinition
	networkIngressDefinitionID := c.Param("id")
	var existingNetworkIngressDefinition v0.NetworkIngressDefinition
	if result := h.DB.First(&existingNetworkIngressDefinition, networkIngressDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// check for empty payload, invalid or unsupported fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, true, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedNetworkIngressDefinition v0.NetworkIngressDefinition
	if err := c.Bind(&updatedNetworkIngressDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	if result := h.DB.Model(&existingNetworkIngressDefinition).Updates(updatedNetworkIngressDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingNetworkIngressDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates an existing network ingress definition by replacing the entire object.
// @Description Replace a network ingress definition in the database.  All required fields must be provided.
// @Description If any optional fields are not provided, they will be null post-update.
// @Description Note: This API endpint is for updating network ingress definition objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID replace-networkIngressDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param networkIngressDefinition body v0.NetworkIngressDefinition true "NetworkIngressDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-definitions/{id} [put]
func (h Handler) ReplaceNetworkIngressDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressDefinition
	networkIngressDefinitionID := c.Param("id")
	var existingNetworkIngressDefinition v0.NetworkIngressDefinition
	if result := h.DB.First(&existingNetworkIngressDefinition, networkIngressDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// check for empty payload, invalid or unsupported fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, true, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedNetworkIngressDefinition v0.NetworkIngressDefinition
	if err := c.Bind(&updatedNetworkIngressDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, updatedNetworkIngressDefinition, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// persist provided data
	updatedNetworkIngressDefinition.ID = existingNetworkIngressDefinition.ID
	if result := h.DB.Session(&gorm.Session{FullSaveAssociations: false}).Omit("CreatedAt", "DeletedAt").Save(&updatedNetworkIngressDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// reload updated data from DB
	if result := h.DB.First(&existingNetworkIngressDefinition, networkIngressDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingNetworkIngressDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a network ingress definition.
// @Description Delete a network ingress definition by ID from the database.
// @ID delete-networkIngressDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 409 {object} v0.Response "Conflict"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-definitions/{id} [delete]
func (h Handler) DeleteNetworkIngressDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressDefinition
	networkIngressDefinitionID := c.Param("id")
	var networkIngressDefinition v0.NetworkIngressDefinition
	if result := h.DB.First(&networkIngressDefinition, networkIngressDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	if result := h.DB.Delete(&networkIngressDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller
	notifPayload, err := networkIngressDefinition.NotificationPayload(
		notifications.NotificationOperationDeleted,
		false,
		0,
	)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}
	h.JS.Publish(v0.NetworkIngressDefinitionDeleteSubject, *notifPayload)

	response, err := v0.CreateResponse(nil, networkIngressDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

///////////////////////////////////////////////////////////////////////////////
// NetworkIngressInstance
///////////////////////////////////////////////////////////////////////////////

// @Summary GetNetworkIngressInstanceVersions gets the supported versions for the network ingress instance API.
// @Description Get the supported API versions for network ingress instances.
// @ID networkIngressInstance-get-versions
// @Produce json
// @Success 200 {object} api.RESTAPIVersions "OK"
// @Router /network-ingress-instances/versions [get]
func (h Handler) GetNetworkIngressInstanceVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeNetworkIngressInstance)])
}

// @Summary adds a new network ingress instance.
// @Description Add a new network ingress instance to the Threeport database.
// @ID add-networkIngressInstance
// @Accept json
// @Produce json
// @Param networkIngressInstance body v0.NetworkIngressInstance true "NetworkIngressInstance object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-instances [post]
func (h Handler) AddNetworkIngressInstance(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressInstance
	var networkIngressInstance v0.NetworkIngressInstance

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&networkIngressInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, networkIngressInstance, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// check for duplicate names
	var existingNetworkIngressInstance v0.NetworkIngressInstance
	nameUsed := true
	result := h.DB.Where("name = ?", networkIngressInstance.Name).First(&existingNetworkIngressInstance)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			nameUsed = false
		} else {
			return iapi.ResponseStatus500(c, nil, result.Error, objectType)
		}
	}
	if nameUsed {
		return iapi.ResponseStatus409(c, nil, errors.New("object with provided name already exists"), objectType)
	}

	// persist to DB
	if result := h.DB.Create(&networkIngressInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller
	notifPayload, err := networkIngressInstance.NotificationPayload(
		notifications.NotificationOperationCreated,
		false,
		0,
	)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}
	h.JS.Publish(v0.NetworkIngressInstanceCreateSubject, *notifPayload)

	response, err := v0.CreateResponse(nil, networkIngressInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all network ingress instances.
// @Description Get all network ingress instances from the Threeport database.
// @ID get-networkIngressInstances
// @Accept json
// @Produce json
// @Param name query string false "network ingress instance search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-instances [get]
func (h Handler) GetNetworkIngressInstances(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressInstance
	params, err := c.(*iapi.CustomContext).GetPaginationParams()
	if err != nil {
		return iapi.ResponseStatus400(c, &params, err, objectType)
	}

	var filter v0.NetworkIngressInstance
	if err := c.Bind(&filter); err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	var totalCount int64
	if result := h.DB.Model(&v0.NetworkIngressInstance{}).Where(&filter).Count(&totalCount); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	records := &[]v0.NetworkIngressInstance{}
	if result := h.DB.Order("ID asc").Where(&filter).Limit(params.Size).Offset((params.Page - 1) * params.Size).Find(records); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a network ingress instance.
// @Description Get a particular network ingress instance from the database.
// @ID get-networkIngressInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-instances/{id} [get]
func (h Handler) GetNetworkIngressInstance(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressInstance
	networkIngressInstanceID := c.Param("id")
	var networkIngressInstance v0.NetworkIngressInstance
	if result := h.DB.First(&networkIngressInstance, networkIngressInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, networkIngressInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates specific fields for an existing network ingress instance.
// @Description Update a network ingress instance in the database.  Provide one or more fields to update.
// @Description Note: This API endpint is for updating network ingress instance objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID update-networkIngressInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param networkIngressInstance body v0.NetworkIngressInstance true "NetworkIngressInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-instances/{id} [patch]
func (h Handler) UpdateNetworkIngressInstance(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressInstance
	networkIngressInstanceID := c.Param("id")
	var existingNetworkIngressInstance v0.NetworkIngressInstance
	if result := h.DB.First(&existingNetworkIngressInstance, networkIngressInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// check for empty payload, invalid or unsupported fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, true, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedNetworkIngressInstance v0.NetworkIngressInstance
	if err := c.Bind(&updatedNetworkIngressInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	if result := h.DB.Model(&existingNetworkIngressInstance).Updates(updatedNetworkIngressInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingNetworkIngressInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates an existing network ingress instance by replacing the entire object.
// @Description Replace a network ingress instance in the database.  All required fields must be provided.
// @Description If any optional fields are not provided, they will be null post-update.
// @Description Note: This API endpint is for updating network ingress instance objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID replace-networkIngressInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param networkIngressInstance body v0.NetworkIngressInstance true "NetworkIngressInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-instances/{id} [put]
func (h Handler) ReplaceNetworkIngressInstance(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressInstance
	networkIngressInstanceID := c.Param("id")
	var existingNetworkIngressInstance v0.NetworkIngressInstance
	if result := h.DB.First(&existingNetworkIngressInstance, networkIngressInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// check for empty payload, invalid or unsupported fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, true, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedNetworkIngressInstance v0.NetworkIngressInstance
	if err := c.Bind(&updatedNetworkIngressInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, updatedNetworkIngressInstance, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// persist provided data
	updatedNetworkIngressInstance.ID = existingNetworkIngressInstance.ID
	if result := h.DB.Session(&gorm.Session{FullSaveAssociations: false}).Omit("CreatedAt", "DeletedAt").Save(&updatedNetworkIngressInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// reload updated data from DB
	if result := h.DB.First(&existingNetworkIngressInstance, networkIngressInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingNetworkIngressInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a network ingress instance.
// @Description Delete a network ingress instance by ID from the database.
// @ID delete-networkIngressInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 409 {object} v0.Response "Conflict"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/network-ingress-instances/{id} [delete]
func (h Handler) DeleteNetworkIngressInstance(c echo.Context) error {
	objectType := v0.ObjectTypeNetworkIngressInstance
	networkIngressInstanceID := c.Param("id")
	var networkIngressInstance v0.NetworkIngressInstance
	if result := h.DB.First(&networkIngressInstance, networkIngressInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	if result := h.DB.Delete(&networkIngressInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller
	notifPayload, err := networkIngressInstance.NotificationPayload(
		notifications.NotificationOperationDeleted,
		false,
		0,
	)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}
	h.JS.Publish(v0.NetworkIngressInstanceDeleteSubject, *notifPayload)

	response, err := v0.CreateResponse(nil, networkIngressInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}
