// generated by 'threeport-sdk gen api-model' - do not edit

package v0

import (
	"encoding/json"
	"fmt"
	notifications "github.com/threeport/threeport/pkg/notifications/v0"
)

const (
	ObjectTypeTerraformDefinition ObjectType = "TerraformDefinition"
	ObjectTypeTerraformInstance   ObjectType = "TerraformInstance"

	TerraformStreamName = "terraformStream"

	TerraformDefinitionSubject       = "terraformDefinition.*"
	TerraformDefinitionCreateSubject = "terraformDefinition.create"
	TerraformDefinitionUpdateSubject = "terraformDefinition.update"
	TerraformDefinitionDeleteSubject = "terraformDefinition.delete"

	TerraformInstanceSubject       = "terraformInstance.*"
	TerraformInstanceCreateSubject = "terraformInstance.create"
	TerraformInstanceUpdateSubject = "terraformInstance.update"
	TerraformInstanceDeleteSubject = "terraformInstance.delete"

	PathTerraformDefinitions = "/v0/terraform-definitions"
	PathTerraformInstances   = "/v0/terraform-instances"
)

// GetTerraformDefinitionSubjects returns the NATS subjects
// for terraform definitions.
func GetTerraformDefinitionSubjects() []string {
	return []string{
		TerraformDefinitionCreateSubject,
		TerraformDefinitionUpdateSubject,
		TerraformDefinitionDeleteSubject,
	}
}

// GetTerraformInstanceSubjects returns the NATS subjects
// for terraform instances.
func GetTerraformInstanceSubjects() []string {
	return []string{
		TerraformInstanceCreateSubject,
		TerraformInstanceUpdateSubject,
		TerraformInstanceDeleteSubject,
	}
}

// GetTerraformSubjects returns the NATS subjects
// for all terraform objects.
func GetTerraformSubjects() []string {
	var terraformSubjects []string

	terraformSubjects = append(terraformSubjects, GetTerraformDefinitionSubjects()...)
	terraformSubjects = append(terraformSubjects, GetTerraformInstanceSubjects()...)

	return terraformSubjects
}

// NotificationPayload returns the notification payload that is delivered to the
// controller when a change is made.  It includes the object as presented by the
// client when the change was made.
func (td *TerraformDefinition) NotificationPayload(
	operation notifications.NotificationOperation,
	requeue bool,
	creationTime int64,
) (*[]byte, error) {
	notif := notifications.Notification{
		CreationTime: &creationTime,
		Object:       td,
		Operation:    operation,
	}

	payload, err := json.Marshal(notif)
	if err != nil {
		return &payload, fmt.Errorf("failed to marshal notification payload %+v: %w", td, err)
	}

	return &payload, nil
}

// DecodeNotifObject takes the threeport object in the form of a
// map[string]interface and returns the typed object by marshalling into JSON
// and then unmarshalling into the typed object.  We are not using the
// mapstructure library here as that requires custom decode hooks to manage
// fields with non-native go types.
func (td *TerraformDefinition) DecodeNotifObject(object interface{}) error {
	jsonObject, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("failed to marshal object map from consumed notification message: %w", err)
	}
	if err := json.Unmarshal(jsonObject, &td); err != nil {
		return fmt.Errorf("failed to unmarshal json object to typed object: %w", err)
	}
	return nil
}

// GetID returns the unique ID for the object.
func (td *TerraformDefinition) GetID() uint {
	return *td.ID
}

// String returns a string representation of the ojbect.
func (td TerraformDefinition) String() string {
	return fmt.Sprintf("v0.TerraformDefinition")
}

// NotificationPayload returns the notification payload that is delivered to the
// controller when a change is made.  It includes the object as presented by the
// client when the change was made.
func (ti *TerraformInstance) NotificationPayload(
	operation notifications.NotificationOperation,
	requeue bool,
	creationTime int64,
) (*[]byte, error) {
	notif := notifications.Notification{
		CreationTime: &creationTime,
		Object:       ti,
		Operation:    operation,
	}

	payload, err := json.Marshal(notif)
	if err != nil {
		return &payload, fmt.Errorf("failed to marshal notification payload %+v: %w", ti, err)
	}

	return &payload, nil
}

// DecodeNotifObject takes the threeport object in the form of a
// map[string]interface and returns the typed object by marshalling into JSON
// and then unmarshalling into the typed object.  We are not using the
// mapstructure library here as that requires custom decode hooks to manage
// fields with non-native go types.
func (ti *TerraformInstance) DecodeNotifObject(object interface{}) error {
	jsonObject, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("failed to marshal object map from consumed notification message: %w", err)
	}
	if err := json.Unmarshal(jsonObject, &ti); err != nil {
		return fmt.Errorf("failed to unmarshal json object to typed object: %w", err)
	}
	return nil
}

// GetID returns the unique ID for the object.
func (ti *TerraformInstance) GetID() uint {
	return *ti.ID
}

// String returns a string representation of the ojbect.
func (ti TerraformInstance) String() string {
	return fmt.Sprintf("v0.TerraformInstance")
}
