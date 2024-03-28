// generated by 'threeport-sdk gen' for API model boilerplate' - do not edit

package handlers

import (
	"errors"
	"fmt"
	echo "github.com/labstack/echo/v4"
	api "github.com/threeport/threeport/pkg/api"
	iapi "github.com/threeport/threeport/pkg/api-server/v0"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	notifications "github.com/threeport/threeport/pkg/notifications/v0"
	gorm "gorm.io/gorm"
	clause "gorm.io/gorm/clause"
	"net/http"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// ControlPlaneDefinition
///////////////////////////////////////////////////////////////////////////////

// @Summary GetControlPlaneDefinitionVersions gets the supported versions for the control plane definition API.
// @Description Get the supported API versions for control plane definitions.
// @ID controlPlaneDefinition-get-versions
// @Produce json
// @Success 200 {object} api.RESTAPIVersions "OK"
// @Router /control-plane-definitions/versions [GET]
func (h Handler) GetControlPlaneDefinitionVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeControlPlaneDefinition)])
}

// @Summary adds a new control plane definition.
// @Description Add a new control plane definition to the Threeport database.
// @ID add-v0-controlPlaneDefinition
// @Accept json
// @Produce json
// @Param controlPlaneDefinition body v0.ControlPlaneDefinition true "ControlPlaneDefinition object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-definitions [POST]
func (h Handler) AddControlPlaneDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneDefinition
	var controlPlaneDefinition v0.ControlPlaneDefinition

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType, controlPlaneDefinition); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&controlPlaneDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, controlPlaneDefinition, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// check for duplicate names
	var existingControlPlaneDefinition v0.ControlPlaneDefinition
	nameUsed := true
	result := h.DB.Where("name = ?", controlPlaneDefinition.Name).First(&existingControlPlaneDefinition)
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
	if result := h.DB.Create(&controlPlaneDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller if reconciliation is required
	if !*controlPlaneDefinition.Reconciled {
		notifPayload, err := controlPlaneDefinition.NotificationPayload(
			notifications.NotificationOperationCreated,
			false,
			time.Now().Unix(),
		)
		if err != nil {
			return iapi.ResponseStatus500(c, nil, err, objectType)
		}
		h.JS.Publish(v0.ControlPlaneDefinitionCreateSubject, *notifPayload)
	}

	response, err := v0.CreateResponse(nil, controlPlaneDefinition, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all control plane definitions.
// @Description Get all control plane definitions from the Threeport database.
// @ID get-v0-controlPlaneDefinitions
// @Accept json
// @Produce json
// @Param name query string false "control plane definition search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-definitions [GET]
func (h Handler) GetControlPlaneDefinitions(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneDefinition
	params, err := c.(*iapi.CustomContext).GetPaginationParams()
	if err != nil {
		return iapi.ResponseStatus400(c, &params, err, objectType)
	}

	var filter v0.ControlPlaneDefinition
	if err := c.Bind(&filter); err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	var totalCount int64
	if result := h.DB.Model(&v0.ControlPlaneDefinition{}).Where(&filter).Count(&totalCount); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	records := &[]v0.ControlPlaneDefinition{}
	if result := h.DB.Order("ID asc").Where(&filter).Limit(params.Size).Offset((params.Page - 1) * params.Size).Find(records); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a control plane definition.
// @Description Get a particular control plane definition from the database.
// @ID get-v0-controlPlaneDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-definitions/{id} [GET]
func (h Handler) GetControlPlaneDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneDefinition
	controlPlaneDefinitionID := c.Param("id")
	var controlPlaneDefinition v0.ControlPlaneDefinition
	if result := h.DB.First(&controlPlaneDefinition, controlPlaneDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, controlPlaneDefinition, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates specific fields for an existing control plane definition.
// @Description Update a control plane definition in the database.  Provide one or more fields to update.
// @Description Note: This API endpint is for updating control plane definition objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID update-v0-controlPlaneDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param controlPlaneDefinition body v0.ControlPlaneDefinition true "ControlPlaneDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-definitions/{id} [PATCH]
func (h Handler) UpdateControlPlaneDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneDefinition
	controlPlaneDefinitionID := c.Param("id")
	var existingControlPlaneDefinition v0.ControlPlaneDefinition
	if result := h.DB.First(&existingControlPlaneDefinition, controlPlaneDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// check for empty payload, invalid or unsupported fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, true, objectType, existingControlPlaneDefinition); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedControlPlaneDefinition v0.ControlPlaneDefinition
	if err := c.Bind(&updatedControlPlaneDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// update object in database
	if result := h.DB.Model(&existingControlPlaneDefinition).Updates(updatedControlPlaneDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller if reconciliation is required
	if !*existingControlPlaneDefinition.Reconciled {
		notifPayload, err := existingControlPlaneDefinition.NotificationPayload(
			notifications.NotificationOperationUpdated,
			false,
			time.Now().Unix(),
		)
		if err != nil {
			return iapi.ResponseStatus500(c, nil, err, objectType)
		}
		h.JS.Publish(v0.ControlPlaneDefinitionUpdateSubject, *notifPayload)
	}

	response, err := v0.CreateResponse(nil, existingControlPlaneDefinition, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates an existing control plane definition by replacing the entire object.
// @Description Replace a control plane definition in the database.  All required fields must be provided.
// @Description If any optional fields are not provided, they will be null post-update.
// @Description Note: This API endpint is for updating control plane definition objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID replace-v0-controlPlaneDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param controlPlaneDefinition body v0.ControlPlaneDefinition true "ControlPlaneDefinition object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-definitions/{id} [PUT]
func (h Handler) ReplaceControlPlaneDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneDefinition
	controlPlaneDefinitionID := c.Param("id")
	var existingControlPlaneDefinition v0.ControlPlaneDefinition
	if result := h.DB.First(&existingControlPlaneDefinition, controlPlaneDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// check for empty payload, invalid or unsupported fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, true, objectType, existingControlPlaneDefinition); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedControlPlaneDefinition v0.ControlPlaneDefinition
	if err := c.Bind(&updatedControlPlaneDefinition); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, updatedControlPlaneDefinition, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// persist provided data
	updatedControlPlaneDefinition.ID = existingControlPlaneDefinition.ID
	if result := h.DB.Session(&gorm.Session{FullSaveAssociations: false}).Omit("CreatedAt", "DeletedAt").Save(&updatedControlPlaneDefinition); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// reload updated data from DB
	if result := h.DB.First(&existingControlPlaneDefinition, controlPlaneDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingControlPlaneDefinition, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a control plane definition.
// @Description Delete a control plane definition by ID from the database.
// @ID delete-v0-controlPlaneDefinition
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 409 {object} v0.Response "Conflict"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-definitions/{id} [DELETE]
func (h Handler) DeleteControlPlaneDefinition(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneDefinition
	controlPlaneDefinitionID := c.Param("id")
	var controlPlaneDefinition v0.ControlPlaneDefinition
	if result := h.DB.Preload("ControlPlaneInstances").First(&controlPlaneDefinition, controlPlaneDefinitionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// check to make sure no dependent instances exist for this definition
	if len(controlPlaneDefinition.ControlPlaneInstances) != 0 {
		err := errors.New("control plane definition has related control plane instances - cannot be deleted")
		return iapi.ResponseStatus409(c, nil, err, objectType)
	}

	// schedule for deletion if not already scheduled
	// if scheduled and reconciled, delete object from DB
	// if scheduled but not reconciled, return 409 (controller is working on it)
	if controlPlaneDefinition.DeletionScheduled == nil {
		// schedule for deletion
		reconciled := false
		timestamp := time.Now().UTC()
		scheduledControlPlaneDefinition := v0.ControlPlaneDefinition{
			Reconciliation: v0.Reconciliation{
				DeletionScheduled: &timestamp,
				Reconciled:        &reconciled,
			}}
		if result := h.DB.Model(&controlPlaneDefinition).Updates(scheduledControlPlaneDefinition); result.Error != nil {
			return iapi.ResponseStatus500(c, nil, result.Error, objectType)
		}
		// notify controller
		notifPayload, err := controlPlaneDefinition.NotificationPayload(
			notifications.NotificationOperationDeleted,
			false,
			time.Now().Unix(),
		)
		if err != nil {
			return iapi.ResponseStatus500(c, nil, err, objectType)
		}
		h.JS.Publish(v0.ControlPlaneDefinitionDeleteSubject, *notifPayload)
	} else {
		if controlPlaneDefinition.DeletionConfirmed == nil {
			// if deletion scheduled but not reconciled, return 409 - deletion
			// already underway
			return iapi.ResponseStatus409(c, nil, errors.New(fmt.Sprintf(
				"object with ID %d already being deleted",
				*controlPlaneDefinition.ID,
			)), objectType)
		} else {
			// object scheduled for deletion and confirmed - it can be deleted
			// from DB
			if result := h.DB.Delete(&controlPlaneDefinition); result.Error != nil {
				return iapi.ResponseStatus500(c, nil, result.Error, objectType)
			}
		}
	}

	response, err := v0.CreateResponse(nil, controlPlaneDefinition, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

///////////////////////////////////////////////////////////////////////////////
// ControlPlaneInstance
///////////////////////////////////////////////////////////////////////////////

// @Summary GetControlPlaneInstanceVersions gets the supported versions for the control plane instance API.
// @Description Get the supported API versions for control plane instances.
// @ID controlPlaneInstance-get-versions
// @Produce json
// @Success 200 {object} api.RESTAPIVersions "OK"
// @Router /control-plane-instances/versions [GET]
func (h Handler) GetControlPlaneInstanceVersions(c echo.Context) error {
	return c.JSON(http.StatusOK, api.RestapiVersions[string(v0.ObjectTypeControlPlaneInstance)])
}

// @Summary adds a new control plane instance.
// @Description Add a new control plane instance to the Threeport database.
// @ID add-v0-controlPlaneInstance
// @Accept json
// @Produce json
// @Param controlPlaneInstance body v0.ControlPlaneInstance true "ControlPlaneInstance object"
// @Success 201 {object} v0.Response "Created"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-instances [POST]
func (h Handler) AddControlPlaneInstance(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneInstance
	var controlPlaneInstance v0.ControlPlaneInstance

	// check for empty payload, unsupported fields, GORM Model fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, false, objectType, controlPlaneInstance); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	if err := c.Bind(&controlPlaneInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, controlPlaneInstance, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// check for duplicate names
	var existingControlPlaneInstance v0.ControlPlaneInstance
	nameUsed := true
	result := h.DB.Where("name = ?", controlPlaneInstance.Name).First(&existingControlPlaneInstance)
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
	if result := h.DB.Create(&controlPlaneInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller if reconciliation is required
	if !*controlPlaneInstance.Reconciled {
		notifPayload, err := controlPlaneInstance.NotificationPayload(
			notifications.NotificationOperationCreated,
			false,
			time.Now().Unix(),
		)
		if err != nil {
			return iapi.ResponseStatus500(c, nil, err, objectType)
		}
		h.JS.Publish(v0.ControlPlaneInstanceCreateSubject, *notifPayload)
	}

	response, err := v0.CreateResponse(nil, controlPlaneInstance, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus201(c, *response)
}

// @Summary gets all control plane instances.
// @Description Get all control plane instances from the Threeport database.
// @ID get-v0-controlPlaneInstances
// @Accept json
// @Produce json
// @Param name query string false "control plane instance search by name"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-instances [GET]
func (h Handler) GetControlPlaneInstances(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneInstance
	params, err := c.(*iapi.CustomContext).GetPaginationParams()
	if err != nil {
		return iapi.ResponseStatus400(c, &params, err, objectType)
	}

	var filter v0.ControlPlaneInstance
	if err := c.Bind(&filter); err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	var totalCount int64
	if result := h.DB.Preload(clause.Associations).Model(&v0.ControlPlaneInstance{}).Where(&filter).Count(&totalCount); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	records := &[]v0.ControlPlaneInstance{}
	if result := h.DB.Preload(clause.Associations).Order("ID asc").Where(&filter).Limit(params.Size).Offset((params.Page - 1) * params.Size).Find(records); result.Error != nil {
		return iapi.ResponseStatus500(c, &params, result.Error, objectType)
	}

	response, err := v0.CreateResponse(v0.CreateMeta(params, totalCount), *records, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, &params, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary gets a control plane instance.
// @Description Get a particular control plane instance from the database.
// @ID get-v0-controlPlaneInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-instances/{id} [GET]
func (h Handler) GetControlPlaneInstance(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneInstance
	controlPlaneInstanceID := c.Param("id")
	var controlPlaneInstance v0.ControlPlaneInstance
	if result := h.DB.Preload(clause.Associations).First(&controlPlaneInstance, controlPlaneInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, controlPlaneInstance, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates specific fields for an existing control plane instance.
// @Description Update a control plane instance in the database.  Provide one or more fields to update.
// @Description Note: This API endpint is for updating control plane instance objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID update-v0-controlPlaneInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param controlPlaneInstance body v0.ControlPlaneInstance true "ControlPlaneInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-instances/{id} [PATCH]
func (h Handler) UpdateControlPlaneInstance(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneInstance
	controlPlaneInstanceID := c.Param("id")
	var existingControlPlaneInstance v0.ControlPlaneInstance
	if result := h.DB.First(&existingControlPlaneInstance, controlPlaneInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// check for empty payload, invalid or unsupported fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, true, objectType, existingControlPlaneInstance); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedControlPlaneInstance v0.ControlPlaneInstance
	if err := c.Bind(&updatedControlPlaneInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// update object in database
	if result := h.DB.Model(&existingControlPlaneInstance).Updates(updatedControlPlaneInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// notify controller if reconciliation is required
	if !*existingControlPlaneInstance.Reconciled {
		notifPayload, err := existingControlPlaneInstance.NotificationPayload(
			notifications.NotificationOperationUpdated,
			false,
			time.Now().Unix(),
		)
		if err != nil {
			return iapi.ResponseStatus500(c, nil, err, objectType)
		}
		h.JS.Publish(v0.ControlPlaneInstanceUpdateSubject, *notifPayload)
	}

	response, err := v0.CreateResponse(nil, existingControlPlaneInstance, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary updates an existing control plane instance by replacing the entire object.
// @Description Replace a control plane instance in the database.  All required fields must be provided.
// @Description If any optional fields are not provided, they will be null post-update.
// @Description Note: This API endpint is for updating control plane instance objects only.
// @Description Request bodies that include related objects will be accepted, however
// @Description the related objects will not be changed.  Call the patch or put method for
// @Description each particular existing object to change them.
// @ID replace-v0-controlPlaneInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param controlPlaneInstance body v0.ControlPlaneInstance true "ControlPlaneInstance object"
// @Success 200 {object} v0.Response "OK"
// @Failure 400 {object} v0.Response "Bad Request"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-instances/{id} [PUT]
func (h Handler) ReplaceControlPlaneInstance(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneInstance
	controlPlaneInstanceID := c.Param("id")
	var existingControlPlaneInstance v0.ControlPlaneInstance
	if result := h.DB.First(&existingControlPlaneInstance, controlPlaneInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// check for empty payload, invalid or unsupported fields, optional associations, etc.
	if id, err := iapi.PayloadCheck(c, true, objectType, existingControlPlaneInstance); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// bind payload
	var updatedControlPlaneInstance v0.ControlPlaneInstance
	if err := c.Bind(&updatedControlPlaneInstance); err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	// check for missing required fields
	if id, err := iapi.ValidateBoundData(c, updatedControlPlaneInstance, objectType); err != nil {
		return iapi.ResponseStatusErr(id, c, nil, errors.New(err.Error()), objectType)
	}

	// persist provided data
	updatedControlPlaneInstance.ID = existingControlPlaneInstance.ID
	if result := h.DB.Session(&gorm.Session{FullSaveAssociations: false}).Omit("CreatedAt", "DeletedAt").Save(&updatedControlPlaneInstance); result.Error != nil {
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// reload updated data from DB
	if result := h.DB.First(&existingControlPlaneInstance, controlPlaneInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	response, err := v0.CreateResponse(nil, existingControlPlaneInstance, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}

// @Summary deletes a control plane instance.
// @Description Delete a control plane instance by ID from the database.
// @ID delete-v0-controlPlaneInstance
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} v0.Response "OK"
// @Failure 404 {object} v0.Response "Not Found"
// @Failure 409 {object} v0.Response "Conflict"
// @Failure 500 {object} v0.Response "Internal Server Error"
// @Router /v0/control-plane-instances/{id} [DELETE]
func (h Handler) DeleteControlPlaneInstance(c echo.Context) error {
	objectType := v0.ObjectTypeControlPlaneInstance
	controlPlaneInstanceID := c.Param("id")
	var controlPlaneInstance v0.ControlPlaneInstance
	if result := h.DB.First(&controlPlaneInstance, controlPlaneInstanceID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return iapi.ResponseStatus404(c, nil, result.Error, objectType)
		}
		return iapi.ResponseStatus500(c, nil, result.Error, objectType)
	}

	// schedule for deletion if not already scheduled
	// if scheduled and reconciled, delete object from DB
	// if scheduled but not reconciled, return 409 (controller is working on it)
	if controlPlaneInstance.DeletionScheduled == nil {
		// schedule for deletion
		reconciled := false
		timestamp := time.Now().UTC()
		scheduledControlPlaneInstance := v0.ControlPlaneInstance{
			Reconciliation: v0.Reconciliation{
				DeletionScheduled: &timestamp,
				Reconciled:        &reconciled,
			}}
		if result := h.DB.Model(&controlPlaneInstance).Updates(scheduledControlPlaneInstance); result.Error != nil {
			return iapi.ResponseStatus500(c, nil, result.Error, objectType)
		}
		// notify controller
		notifPayload, err := controlPlaneInstance.NotificationPayload(
			notifications.NotificationOperationDeleted,
			false,
			time.Now().Unix(),
		)
		if err != nil {
			return iapi.ResponseStatus500(c, nil, err, objectType)
		}
		h.JS.Publish(v0.ControlPlaneInstanceDeleteSubject, *notifPayload)
	} else {
		if controlPlaneInstance.DeletionConfirmed == nil {
			// if deletion scheduled but not reconciled, return 409 - deletion
			// already underway
			return iapi.ResponseStatus409(c, nil, errors.New(fmt.Sprintf(
				"object with ID %d already being deleted",
				*controlPlaneInstance.ID,
			)), objectType)
		} else {
			// object scheduled for deletion and confirmed - it can be deleted
			// from DB
			if result := h.DB.Delete(&controlPlaneInstance); result.Error != nil {
				return iapi.ResponseStatus500(c, nil, result.Error, objectType)
			}
		}
	}

	response, err := v0.CreateResponse(nil, controlPlaneInstance, objectType)
	if err != nil {
		return iapi.ResponseStatus500(c, nil, err, objectType)
	}

	return iapi.ResponseStatus200(c, *response)
}
