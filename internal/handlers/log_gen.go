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
// LogBackend
///////////////////////////////////////////////////////////////////////////////

// @Summary GetLogBackendVersions gets the supported versions for the log backend API.
// @Description Get the supported API versions for log backends.
// @ID logBackend-get-versions
// @Produce json
// @Success 200 {object} api.RESTAPIVersions "OK"
// @Router /log_backends/versions [get]
func (h Handler) GetLogBackendVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeLogBackend)])
}

// @Summary adds a new log backend.
// @Description Add a new log backend to the Threeport database.
// @ID add-logBackend
// @Accept json
// @Produce json
// @Param logBackend body v0.LogBackend true "LogBackend object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_backends [post]
func (h Handler) AddLogBackend(c echo.Context) error {
	objectType := v0.ObjectTypeLogBackend
	var logBackend v0.LogBackend

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&logBackend); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, logBackend, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if result := h.DB.Create(&logBackend); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller
	notifPayload, err := logBackend.NotificationPayload(false, 0)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}
	h.JS.Publish(v0.LogBackendCreateSubject, *notifPayload)

	response, err := v0.CreateResponse(nil, logBackend)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all log backends.
// @Description Get all log backends from the Threeport database.
// @ID get-logBackends
// @Accept json
// @Produce json
// @Param name query string false "log backend search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_backends [get]
func (h Handler) GetLogBackends(c echo.Context) error {
	objectType := v0.ObjectTypeLogBackend
	params, err := c.(*iapi.CustomContext).GetPaginationParams()
	if err != nil {
		return iapi.ResponseStatus400(c, &params, err, objectType)
	}

	var filter v0.LogBackend
	if err := c.Bind(&filter); err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	var totalCount int64
	if result := h.DB.Model(&v0.LogBackend{}).Where(&filter).Count(&totalCount); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	records := &[]v0.LogBackend{}
	if result := h.DB.Order("ID asc").Where(&filter).Limit(params.Size).Offset((params.Page - 1) * params.Size).Find(records); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a log backend.
// @Description Get a particular log backend from the database.
// @ID get-logBackend
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_backends/{id} [get]
func (h Handler) GetLogBackend(c echo.Context) error {
	objectType := v0.ObjectTypeLogBackend
	logBackendID := c.Param("id")
	var logBackend v0.LogBackend
	if result := h.DB.First(&logBackend, logBackendID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logBackend)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates specific fields for an existing log backend.
// @Description Update a log backend in the database.  Provide one or more fields to update.
// @Description Note: This API endpint is for updating log backend objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID update-logBackend
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logBackend body v0.LogBackend true "LogBackend object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_backends/{id} [patch]
func (h Handler) UpdateLogBackend(c echo.Context) error {
	objectType := v0.ObjectTypeLogBackend
	logBackendID := c.Param("id")
	var existingLogBackend v0.LogBackend
	if result := h.DB.First(&existingLogBackend, logBackendID); result.Error != nil {
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
	var updatedLogBackend v0.LogBackend
	if err := c.Bind(&updatedLogBackend); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	if result := h.DB.Model(&existingLogBackend).Updates(updatedLogBackend); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingLogBackend)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates an existing log backend by replacing the entire object.
// @Description Replace a log backend in the database.  All required fields must be provided.
// @Description If any optional fields are not provided, they will be null post-update.
// @Description Note: This API endpint is for updating log backend objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID replace-logBackend
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logBackend body v0.LogBackend true "LogBackend object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_backends/{id} [put]
func (h Handler) ReplaceLogBackend(c echo.Context) error {
	objectType := v0.ObjectTypeLogBackend
	logBackendID := c.Param("id")
	var existingLogBackend v0.LogBackend
	if result := h.DB.First(&existingLogBackend, logBackendID); result.Error != nil {
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
	var updatedLogBackend v0.LogBackend
	if err := c.Bind(&updatedLogBackend); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, updatedLogBackend, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// persist provided data
	updatedLogBackend.ID = existingLogBackend.ID
	if result := h.DB.Session(&gorm.Session{FullSaveAssociations: false}).Omit("CreatedAt", "DeletedAt").Save(&updatedLogBackend); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// reload updated data from DB
	if result := h.DB.First(&existingLogBackend, logBackendID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingLogBackend)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a log backend.
// @Description Delete a log backend by from the database.
// @ID delete-logBackend
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_backends/{id} [delete]
func (h Handler) DeleteLogBackend(c echo.Context) error {
	objectType := v0.ObjectTypeLogBackend
	logBackendID := c.Param("id")
	var logBackend v0.LogBackend
	if result := h.DB.First(&logBackend, logBackendID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	if result := h.DB.Delete(&logBackend); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logBackend)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

///////////////////////////////////////////////////////////////////////////////
// LogStorageDefinition
///////////////////////////////////////////////////////////////////////////////

// @Summary GetLogStorageDefinitionVersions gets the supported versions for the log storage definition API.
// @Description Get the supported API versions for log storage definitions.
// @ID logStorageDefinition-get-versions
// @Produce json
// @Success 200 {object} api.RESTAPIVersions "OK"
// @Router /log_storage_definitions/versions [get]
func (h Handler) GetLogStorageDefinitionVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeLogStorageDefinition)])
}

// @Summary adds a new log storage definition.
// @Description Add a new log storage definition to the Threeport database.
// @ID add-logStorageDefinition
// @Accept json
// @Produce json
// @Param logStorageDefinition body v0.LogStorageDefinition true "LogStorageDefinition object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_definitions [post]
func (h Handler) AddLogStorageDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageDefinition
	var logStorageDefinition v0.LogStorageDefinition

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&logStorageDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, logStorageDefinition, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if result := h.DB.Create(&logStorageDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller
	notifPayload, err := logStorageDefinition.NotificationPayload(false, 0)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}
	h.JS.Publish(v0.LogStorageDefinitionCreateSubject, *notifPayload)

	response, err := v0.CreateResponse(nil, logStorageDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all log storage definitions.
// @Description Get all log storage definitions from the Threeport database.
// @ID get-logStorageDefinitions
// @Accept json
// @Produce json
// @Param name query string false "log storage definition search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_definitions [get]
func (h Handler) GetLogStorageDefinitions(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageDefinition
	params, err := c.(*iapi.CustomContext).GetPaginationParams()
	if err != nil {
		return iapi.ResponseStatus400(c, &params, err, objectType)
	}

	var filter v0.LogStorageDefinition
	if err := c.Bind(&filter); err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	var totalCount int64
	if result := h.DB.Model(&v0.LogStorageDefinition{}).Where(&filter).Count(&totalCount); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	records := &[]v0.LogStorageDefinition{}
	if result := h.DB.Order("ID asc").Where(&filter).Limit(params.Size).Offset((params.Page - 1) * params.Size).Find(records); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a log storage definition.
// @Description Get a particular log storage definition from the database.
// @ID get-logStorageDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_definitions/{id} [get]
func (h Handler) GetLogStorageDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageDefinition
	logStorageDefinitionID := c.Param("id")
	var logStorageDefinition v0.LogStorageDefinition
	if result := h.DB.First(&logStorageDefinition, logStorageDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logStorageDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates specific fields for an existing log storage definition.
// @Description Update a log storage definition in the database.  Provide one or more fields to update.
// @Description Note: This API endpint is for updating log storage definition objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID update-logStorageDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logStorageDefinition body v0.LogStorageDefinition true "LogStorageDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_definitions/{id} [patch]
func (h Handler) UpdateLogStorageDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageDefinition
	logStorageDefinitionID := c.Param("id")
	var existingLogStorageDefinition v0.LogStorageDefinition
	if result := h.DB.First(&existingLogStorageDefinition, logStorageDefinitionID); result.Error != nil {
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
	var updatedLogStorageDefinition v0.LogStorageDefinition
	if err := c.Bind(&updatedLogStorageDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	if result := h.DB.Model(&existingLogStorageDefinition).Updates(updatedLogStorageDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingLogStorageDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates an existing log storage definition by replacing the entire object.
// @Description Replace a log storage definition in the database.  All required fields must be provided.
// @Description If any optional fields are not provided, they will be null post-update.
// @Description Note: This API endpint is for updating log storage definition objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID replace-logStorageDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logStorageDefinition body v0.LogStorageDefinition true "LogStorageDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_definitions/{id} [put]
func (h Handler) ReplaceLogStorageDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageDefinition
	logStorageDefinitionID := c.Param("id")
	var existingLogStorageDefinition v0.LogStorageDefinition
	if result := h.DB.First(&existingLogStorageDefinition, logStorageDefinitionID); result.Error != nil {
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
	var updatedLogStorageDefinition v0.LogStorageDefinition
	if err := c.Bind(&updatedLogStorageDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, updatedLogStorageDefinition, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// persist provided data
	updatedLogStorageDefinition.ID = existingLogStorageDefinition.ID
	if result := h.DB.Session(&gorm.Session{FullSaveAssociations: false}).Omit("CreatedAt", "DeletedAt").Save(&updatedLogStorageDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// reload updated data from DB
	if result := h.DB.First(&existingLogStorageDefinition, logStorageDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingLogStorageDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a log storage definition.
// @Description Delete a log storage definition by from the database.
// @ID delete-logStorageDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_definitions/{id} [delete]
func (h Handler) DeleteLogStorageDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageDefinition
	logStorageDefinitionID := c.Param("id")
	var logStorageDefinition v0.LogStorageDefinition
	if result := h.DB.First(&logStorageDefinition, logStorageDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	if result := h.DB.Delete(&logStorageDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logStorageDefinition)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

///////////////////////////////////////////////////////////////////////////////
// LogStorageInstance
///////////////////////////////////////////////////////////////////////////////

// @Summary GetLogStorageInstanceVersions gets the supported versions for the log storage instance API.
// @Description Get the supported API versions for log storage instances.
// @ID logStorageInstance-get-versions
// @Produce json
// @Success 200 {object} api.RESTAPIVersions "OK"
// @Router /log_storage_instances/versions [get]
func (h Handler) GetLogStorageInstanceVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeLogStorageInstance)])
}

// @Summary adds a new log storage instance.
// @Description Add a new log storage instance to the Threeport database.
// @ID add-logStorageInstance
// @Accept json
// @Produce json
// @Param logStorageInstance body v0.LogStorageInstance true "LogStorageInstance object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_instances [post]
func (h Handler) AddLogStorageInstance(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageInstance
	var logStorageInstance v0.LogStorageInstance

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&logStorageInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, logStorageInstance, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if result := h.DB.Create(&logStorageInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller
	notifPayload, err := logStorageInstance.NotificationPayload(false, 0)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}
	h.JS.Publish(v0.LogStorageInstanceCreateSubject, *notifPayload)

	response, err := v0.CreateResponse(nil, logStorageInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all log storage instances.
// @Description Get all log storage instances from the Threeport database.
// @ID get-logStorageInstances
// @Accept json
// @Produce json
// @Param name query string false "log storage instance search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_instances [get]
func (h Handler) GetLogStorageInstances(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageInstance
	params, err := c.(*iapi.CustomContext).GetPaginationParams()
	if err != nil {
		return iapi.ResponseStatus400(c, &params, err, objectType)
	}

	var filter v0.LogStorageInstance
	if err := c.Bind(&filter); err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	var totalCount int64
	if result := h.DB.Model(&v0.LogStorageInstance{}).Where(&filter).Count(&totalCount); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	records := &[]v0.LogStorageInstance{}
	if result := h.DB.Order("ID asc").Where(&filter).Limit(params.Size).Offset((params.Page - 1) * params.Size).Find(records); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a log storage instance.
// @Description Get a particular log storage instance from the database.
// @ID get-logStorageInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_instances/{id} [get]
func (h Handler) GetLogStorageInstance(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageInstance
	logStorageInstanceID := c.Param("id")
	var logStorageInstance v0.LogStorageInstance
	if result := h.DB.First(&logStorageInstance, logStorageInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logStorageInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates specific fields for an existing log storage instance.
// @Description Update a log storage instance in the database.  Provide one or more fields to update.
// @Description Note: This API endpint is for updating log storage instance objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID update-logStorageInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logStorageInstance body v0.LogStorageInstance true "LogStorageInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_instances/{id} [patch]
func (h Handler) UpdateLogStorageInstance(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageInstance
	logStorageInstanceID := c.Param("id")
	var existingLogStorageInstance v0.LogStorageInstance
	if result := h.DB.First(&existingLogStorageInstance, logStorageInstanceID); result.Error != nil {
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
	var updatedLogStorageInstance v0.LogStorageInstance
	if err := c.Bind(&updatedLogStorageInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	if result := h.DB.Model(&existingLogStorageInstance).Updates(updatedLogStorageInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingLogStorageInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates an existing log storage instance by replacing the entire object.
// @Description Replace a log storage instance in the database.  All required fields must be provided.
// @Description If any optional fields are not provided, they will be null post-update.
// @Description Note: This API endpint is for updating log storage instance objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID replace-logStorageInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logStorageInstance body v0.LogStorageInstance true "LogStorageInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_instances/{id} [put]
func (h Handler) ReplaceLogStorageInstance(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageInstance
	logStorageInstanceID := c.Param("id")
	var existingLogStorageInstance v0.LogStorageInstance
	if result := h.DB.First(&existingLogStorageInstance, logStorageInstanceID); result.Error != nil {
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
	var updatedLogStorageInstance v0.LogStorageInstance
	if err := c.Bind(&updatedLogStorageInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, updatedLogStorageInstance, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// persist provided data
	updatedLogStorageInstance.ID = existingLogStorageInstance.ID
	if result := h.DB.Session(&gorm.Session{FullSaveAssociations: false}).Omit("CreatedAt", "DeletedAt").Save(&updatedLogStorageInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// reload updated data from DB
	if result := h.DB.First(&existingLogStorageInstance, logStorageInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingLogStorageInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a log storage instance.
// @Description Delete a log storage instance by from the database.
// @ID delete-logStorageInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log_storage_instances/{id} [delete]
func (h Handler) DeleteLogStorageInstance(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageInstance
	logStorageInstanceID := c.Param("id")
	var logStorageInstance v0.LogStorageInstance
	if result := h.DB.First(&logStorageInstance, logStorageInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	if result := h.DB.Delete(&logStorageInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logStorageInstance)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}
