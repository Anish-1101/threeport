// generated by 'threeport-sdk gen api-model' - do not edit

package handlers

import (
	"errors"
	echo "github.com/labstack/echo/v4"
	api "github.com/threeport/threeport/pkg/api"
	iapi "github.com/threeport/threeport/pkg/api-server/v0"
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
// @Router /log-backends/versions [GET]
func (h Handler) GetLogBackendVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeLogBackend)])
}

// @Summary adds a new log backend.
// @Description Add a new log backend to the Threeport database.
// @ID add-v0-logBackend
// @Accept json
// @Produce json
// @Param logBackend body v0.LogBackend true "LogBackend object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-backends [POST]
func (h Handler) AddLogBackend(c echo.Context) error {
	objectType := v0.ObjectTypeLogBackend
	var logBackend v0.LogBackend

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType, logBackend); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&logBackend); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, logBackend, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// check for duplicate names
	var existingLogBackend v0.LogBackend
	nameUsed := true
	result := h.DB.Where("name = ?", logBackend.Name).First(&existingLogBackend)
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
	if result := h.DB.Create(&logBackend); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logBackend, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all log backends.
// @Description Get all log backends from the Threeport database.
// @ID get-v0-logBackends
// @Accept json
// @Produce json
// @Param name query string false "log backend search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-backends [GET]
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

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a log backend.
// @Description Get a particular log backend from the database.
// @ID get-v0-logBackend
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-backends/{id} [GET]
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

	response, err := v0.CreateResponse(nil, logBackend, objectType)
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
// @ID update-v0-logBackend
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logBackend body v0.LogBackend true "LogBackend object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-backends/{id} [PATCH]
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
	if id, err := iapi.PayloadCheck(c, true, objectType, existingLogBackend); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedLogBackend v0.LogBackend
	if err := c.Bind(&updatedLogBackend); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// update object in database
	if result := h.DB.Model(&existingLogBackend).Updates(updatedLogBackend); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingLogBackend, objectType)
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
// @ID replace-v0-logBackend
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logBackend body v0.LogBackend true "LogBackend object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-backends/{id} [PUT]
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
	if id, err := iapi.PayloadCheck(c, true, objectType, existingLogBackend); err != nil {
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

	response, err := v0.CreateResponse(nil, existingLogBackend, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a log backend.
// @Description Delete a log backend by ID from the database.
// @ID delete-v0-logBackend
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 409 {object} v0.Response "Conflict"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-backends/{id} [DELETE]
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

	// delete object
	if result := h.DB.Delete(&logBackend); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logBackend, objectType)
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
// @Router /log-storage-definitions/versions [GET]
func (h Handler) GetLogStorageDefinitionVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeLogStorageDefinition)])
}

// @Summary adds a new log storage definition.
// @Description Add a new log storage definition to the Threeport database.
// @ID add-v0-logStorageDefinition
// @Accept json
// @Produce json
// @Param logStorageDefinition body v0.LogStorageDefinition true "LogStorageDefinition object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-definitions [POST]
func (h Handler) AddLogStorageDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageDefinition
	var logStorageDefinition v0.LogStorageDefinition

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType, logStorageDefinition); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&logStorageDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, logStorageDefinition, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// check for duplicate names
	var existingLogStorageDefinition v0.LogStorageDefinition
	nameUsed := true
	result := h.DB.Where("name = ?", logStorageDefinition.Name).First(&existingLogStorageDefinition)
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
	if result := h.DB.Create(&logStorageDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logStorageDefinition, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all log storage definitions.
// @Description Get all log storage definitions from the Threeport database.
// @ID get-v0-logStorageDefinitions
// @Accept json
// @Produce json
// @Param name query string false "log storage definition search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-definitions [GET]
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

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a log storage definition.
// @Description Get a particular log storage definition from the database.
// @ID get-v0-logStorageDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-definitions/{id} [GET]
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

	response, err := v0.CreateResponse(nil, logStorageDefinition, objectType)
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
// @ID update-v0-logStorageDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logStorageDefinition body v0.LogStorageDefinition true "LogStorageDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-definitions/{id} [PATCH]
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
	if id, err := iapi.PayloadCheck(c, true, objectType, existingLogStorageDefinition); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedLogStorageDefinition v0.LogStorageDefinition
	if err := c.Bind(&updatedLogStorageDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// update object in database
	if result := h.DB.Model(&existingLogStorageDefinition).Updates(updatedLogStorageDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingLogStorageDefinition, objectType)
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
// @ID replace-v0-logStorageDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logStorageDefinition body v0.LogStorageDefinition true "LogStorageDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-definitions/{id} [PUT]
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
	if id, err := iapi.PayloadCheck(c, true, objectType, existingLogStorageDefinition); err != nil {
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

	response, err := v0.CreateResponse(nil, existingLogStorageDefinition, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a log storage definition.
// @Description Delete a log storage definition by ID from the database.
// @ID delete-v0-logStorageDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 409 {object} v0.Response "Conflict"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-definitions/{id} [DELETE]
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

	// delete object
	if result := h.DB.Delete(&logStorageDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logStorageDefinition, objectType)
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
// @Router /log-storage-instances/versions [GET]
func (h Handler) GetLogStorageInstanceVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeLogStorageInstance)])
}

// @Summary adds a new log storage instance.
// @Description Add a new log storage instance to the Threeport database.
// @ID add-v0-logStorageInstance
// @Accept json
// @Produce json
// @Param logStorageInstance body v0.LogStorageInstance true "LogStorageInstance object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-instances [POST]
func (h Handler) AddLogStorageInstance(c echo.Context) error {
	objectType := v0.ObjectTypeLogStorageInstance
	var logStorageInstance v0.LogStorageInstance

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType, logStorageInstance); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&logStorageInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, logStorageInstance, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// check for duplicate names
	var existingLogStorageInstance v0.LogStorageInstance
	nameUsed := true
	result := h.DB.Where("name = ?", logStorageInstance.Name).First(&existingLogStorageInstance)
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
	if result := h.DB.Create(&logStorageInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logStorageInstance, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all log storage instances.
// @Description Get all log storage instances from the Threeport database.
// @ID get-v0-logStorageInstances
// @Accept json
// @Produce json
// @Param name query string false "log storage instance search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-instances [GET]
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

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a log storage instance.
// @Description Get a particular log storage instance from the database.
// @ID get-v0-logStorageInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-instances/{id} [GET]
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

	response, err := v0.CreateResponse(nil, logStorageInstance, objectType)
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
// @ID update-v0-logStorageInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logStorageInstance body v0.LogStorageInstance true "LogStorageInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-instances/{id} [PATCH]
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
	if id, err := iapi.PayloadCheck(c, true, objectType, existingLogStorageInstance); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedLogStorageInstance v0.LogStorageInstance
	if err := c.Bind(&updatedLogStorageInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// update object in database
	if result := h.DB.Model(&existingLogStorageInstance).Updates(updatedLogStorageInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingLogStorageInstance, objectType)
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
// @ID replace-v0-logStorageInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param logStorageInstance body v0.LogStorageInstance true "LogStorageInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-instances/{id} [PUT]
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
	if id, err := iapi.PayloadCheck(c, true, objectType, existingLogStorageInstance); err != nil {
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

	response, err := v0.CreateResponse(nil, existingLogStorageInstance, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a log storage instance.
// @Description Delete a log storage instance by ID from the database.
// @ID delete-v0-logStorageInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 409 {object} v0.Response "Conflict"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/log-storage-instances/{id} [DELETE]
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

	// delete object
	if result := h.DB.Delete(&logStorageInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, logStorageInstance, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}
