// generated by 'threeport-codegen api-model' - do not edit

package handlers

import (
	"errors"
	echo "github.com/labstack/echo/v4"
	iapi "github.com/threeport/threeport/internal/api"
	api "github.com/threeport/threeport/pkg/api"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	gorm "gorm.io/gorm"
	"net/http"
)

///////////////////////////////////////////////////////////////////////////////
// EthereumNodeDefinition
///////////////////////////////////////////////////////////////////////////////

// @Summary GetEthereumNodeDefinitionVersions gets the supported versions for the ethereum node definition API.
// @Description Get the supported API versions for ethereum node definitions.
// @ID ethereumNodeDefinition-get-versions
// @Produce json
// @Success 200 {object} api.RESTAPIVersions "OK"
// @Router /ethereum_node_definitions/versions [get]
func (h Handler) GetEthereumNodeDefinitionVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeEthereumNodeDefinition)])
}

// @Summary adds a new ethereum node definition.
// @Description Add a new ethereum node definition to the Threeport database.
// @ID add-ethereumNodeDefinition
// @Accept json
// @Produce json
// @Param ethereumNodeDefinition body v0.EthereumNodeDefinition true "EthereumNodeDefinition object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_definitions [post]
func (h Handler) AddEthereumNodeDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeDefinition
	var ethereumNodeDefinition v0.EthereumNodeDefinition

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&ethereumNodeDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, ethereumNodeDefinition, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if result := h.DB.Create(&ethereumNodeDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller
	notifPayload, err := ethereumNodeDefinition.NotificationPayload(false, 0)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}
	h.JS.Publish(v0.EthereumNodeDefinitionCreateSubject, *notifPayload)

	response, err := v0.CreateResponse(nil, ethereumNodeDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all ethereum node definitions.
// @Description Get all ethereum node definitions from the Threeport database.
// @ID get-ethereumNodeDefinitions
// @Accept json
// @Produce json
// @Param name query string false "ethereum node definition search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_definitions [get]
func (h Handler) GetEthereumNodeDefinitions(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeDefinition
	params, err := c.(*iapi.CustomContext).GetPaginationParams()
	if err != nil {
		return iapi.ResponseStatus400(c, &params, err, objectType)
	}

	var filter v0.EthereumNodeDefinition
	if err := c.Bind(&filter); err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	var totalCount int64
	if result := h.DB.Model(&v0.EthereumNodeDefinition{}).Where(&filter).Count(&totalCount); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	records := &[]v0.EthereumNodeDefinition{}
	if result := h.DB.Order("ID asc").Where(&filter).Limit(params.Size).Offset((params.Page - 1) * params.Size).Find(records); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a ethereum node definition.
// @Description Get a particular ethereum node definition from the database.
// @ID get-ethereumNodeDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_definitions/{id} [get]
func (h Handler) GetEthereumNodeDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeDefinition
	ethereumNodeDefinitionID := c.Param("id")
	var ethereumNodeDefinition v0.EthereumNodeDefinition
	if result := h.DB.First(&ethereumNodeDefinition, ethereumNodeDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, ethereumNodeDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates specific fields for an existing ethereum node definition.
// @Description Update a ethereum node definition in the database.  Provide one or more fields to update.
// @Description Note: This API endpint is for updating ethereum node definition objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID update-ethereumNodeDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param ethereumNodeDefinition body v0.EthereumNodeDefinition true "EthereumNodeDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_definitions/{id} [patch]
func (h Handler) UpdateEthereumNodeDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeDefinition
	ethereumNodeDefinitionID := c.Param("id")
	var existingEthereumNodeDefinition v0.EthereumNodeDefinition
	if result := h.DB.First(&existingEthereumNodeDefinition, ethereumNodeDefinitionID); result.Error != nil {
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
	var updatedEthereumNodeDefinition v0.EthereumNodeDefinition
	if err := c.Bind(&updatedEthereumNodeDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	if result := h.DB.Model(&existingEthereumNodeDefinition).Updates(updatedEthereumNodeDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingEthereumNodeDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates an existing ethereum node definition by replacing the entire object.
// @Description Replace a ethereum node definition in the database.  All required fields must be provided.
// @Description If any optional fields are not provided, they will be null post-update.
// @Description Note: This API endpint is for updating ethereum node definition objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID replace-ethereumNodeDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param ethereumNodeDefinition body v0.EthereumNodeDefinition true "EthereumNodeDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_definitions/{id} [put]
func (h Handler) ReplaceEthereumNodeDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeDefinition
	ethereumNodeDefinitionID := c.Param("id")
	var existingEthereumNodeDefinition v0.EthereumNodeDefinition
	if result := h.DB.First(&existingEthereumNodeDefinition, ethereumNodeDefinitionID); result.Error != nil {
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
	var updatedEthereumNodeDefinition v0.EthereumNodeDefinition
	if err := c.Bind(&updatedEthereumNodeDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, updatedEthereumNodeDefinition, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// persist provided data
	updatedEthereumNodeDefinition.ID = existingEthereumNodeDefinition.ID
	if result := h.DB.Session(&gorm.Session{FullSaveAssociations: false}).Omit("CreatedAt", "DeletedAt").Save(&updatedEthereumNodeDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// reload updated data from DB
	if result := h.DB.First(&existingEthereumNodeDefinition, ethereumNodeDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingEthereumNodeDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a ethereum node definition.
// @Description Delete a ethereum node definition by from the database.
// @ID delete-ethereumNodeDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_definitions/{id} [delete]
func (h Handler) DeleteEthereumNodeDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeDefinition
	ethereumNodeDefinitionID := c.Param("id")
	var ethereumNodeDefinition v0.EthereumNodeDefinition
	if result := h.DB.First(&ethereumNodeDefinition, ethereumNodeDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	if result := h.DB.Delete(&ethereumNodeDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, ethereumNodeDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

///////////////////////////////////////////////////////////////////////////////
// EthereumNodeInstance
///////////////////////////////////////////////////////////////////////////////

// @Summary GetEthereumNodeInstanceVersions gets the supported versions for the ethereum node instance API.
// @Description Get the supported API versions for ethereum node instances.
// @ID ethereumNodeInstance-get-versions
// @Produce json
// @Success 200 {object} api.RESTAPIVersions "OK"
// @Router /ethereum_node_instances/versions [get]
func (h Handler) GetEthereumNodeInstanceVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeEthereumNodeInstance)])
}

// @Summary adds a new ethereum node instance.
// @Description Add a new ethereum node instance to the Threeport database.
// @ID add-ethereumNodeInstance
// @Accept json
// @Produce json
// @Param ethereumNodeInstance body v0.EthereumNodeInstance true "EthereumNodeInstance object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_instances [post]
func (h Handler) AddEthereumNodeInstance(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeInstance
	var ethereumNodeInstance v0.EthereumNodeInstance

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&ethereumNodeInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, ethereumNodeInstance, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if result := h.DB.Create(&ethereumNodeInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller
	notifPayload, err := ethereumNodeInstance.NotificationPayload(false, 0)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}
	h.JS.Publish(v0.EthereumNodeInstanceCreateSubject, *notifPayload)

	response, err := v0.CreateResponse(nil, ethereumNodeInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all ethereum node instances.
// @Description Get all ethereum node instances from the Threeport database.
// @ID get-ethereumNodeInstances
// @Accept json
// @Produce json
// @Param name query string false "ethereum node instance search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_instances [get]
func (h Handler) GetEthereumNodeInstances(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeInstance
	params, err := c.(*iapi.CustomContext).GetPaginationParams()
	if err != nil {
		return iapi.ResponseStatus400(c, &params, err, objectType)
	}

	var filter v0.EthereumNodeInstance
	if err := c.Bind(&filter); err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	var totalCount int64
	if result := h.DB.Model(&v0.EthereumNodeInstance{}).Where(&filter).Count(&totalCount); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	records := &[]v0.EthereumNodeInstance{}
	if result := h.DB.Order("ID asc").Where(&filter).Limit(params.Size).Offset((params.Page - 1) * params.Size).Find(records); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a ethereum node instance.
// @Description Get a particular ethereum node instance from the database.
// @ID get-ethereumNodeInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_instances/{id} [get]
func (h Handler) GetEthereumNodeInstance(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeInstance
	ethereumNodeInstanceID := c.Param("id")
	var ethereumNodeInstance v0.EthereumNodeInstance
	if result := h.DB.First(&ethereumNodeInstance, ethereumNodeInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, ethereumNodeInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates specific fields for an existing ethereum node instance.
// @Description Update a ethereum node instance in the database.  Provide one or more fields to update.
// @Description Note: This API endpint is for updating ethereum node instance objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID update-ethereumNodeInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param ethereumNodeInstance body v0.EthereumNodeInstance true "EthereumNodeInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_instances/{id} [patch]
func (h Handler) UpdateEthereumNodeInstance(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeInstance
	ethereumNodeInstanceID := c.Param("id")
	var existingEthereumNodeInstance v0.EthereumNodeInstance
	if result := h.DB.First(&existingEthereumNodeInstance, ethereumNodeInstanceID); result.Error != nil {
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
	var updatedEthereumNodeInstance v0.EthereumNodeInstance
	if err := c.Bind(&updatedEthereumNodeInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	if result := h.DB.Model(&existingEthereumNodeInstance).Updates(updatedEthereumNodeInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingEthereumNodeInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates an existing ethereum node instance by replacing the entire object.
// @Description Replace a ethereum node instance in the database.  All required fields must be provided.
// @Description If any optional fields are not provided, they will be null post-update.
// @Description Note: This API endpint is for updating ethereum node instance objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID replace-ethereumNodeInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param ethereumNodeInstance body v0.EthereumNodeInstance true "EthereumNodeInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_instances/{id} [put]
func (h Handler) ReplaceEthereumNodeInstance(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeInstance
	ethereumNodeInstanceID := c.Param("id")
	var existingEthereumNodeInstance v0.EthereumNodeInstance
	if result := h.DB.First(&existingEthereumNodeInstance, ethereumNodeInstanceID); result.Error != nil {
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
	var updatedEthereumNodeInstance v0.EthereumNodeInstance
	if err := c.Bind(&updatedEthereumNodeInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, updatedEthereumNodeInstance, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// persist provided data
	updatedEthereumNodeInstance.ID = existingEthereumNodeInstance.ID
	if result := h.DB.Session(&gorm.Session{FullSaveAssociations: false}).Omit("CreatedAt", "DeletedAt").Save(&updatedEthereumNodeInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// reload updated data from DB
	if result := h.DB.First(&existingEthereumNodeInstance, ethereumNodeInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingEthereumNodeInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a ethereum node instance.
// @Description Delete a ethereum node instance by from the database.
// @ID delete-ethereumNodeInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/ethereum_node_instances/{id} [delete]
func (h Handler) DeleteEthereumNodeInstance(c echo.Context) error {
	objectType := v0.ObjectTypeEthereumNodeInstance
	ethereumNodeInstanceID := c.Param("id")
	var ethereumNodeInstance v0.EthereumNodeInstance
	if result := h.DB.First(&ethereumNodeInstance, ethereumNodeInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	if result := h.DB.Delete(&ethereumNodeInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, ethereumNodeInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}