// generated by 'threeport-codegen api-model' - do not edit

package v0

import (
	"encoding/json"
	"fmt"
	notifications "github.com/threeport/threeport/pkg/notifications"
)

const (
	ObjectTypeEthereumNodeDefinition ObjectType = "EthereumNodeDefinition"
	ObjectTypeEthereumNodeInstance   ObjectType = "EthereumNodeInstance"

	EthereumNodeStreamName = "ethereumNodeStream"

	EthereumNodeDefinitionSubject       = "ethereumNodeDefinition.*"
	EthereumNodeDefinitionCreateSubject = "ethereumNodeDefinition.create"
	EthereumNodeDefinitionUpdateSubject = "ethereumNodeDefinition.update"
	EthereumNodeDefinitionDeleteSubject = "ethereumNodeDefinition.delete"

	EthereumNodeInstanceSubject       = "ethereumNodeInstance.*"
	EthereumNodeInstanceCreateSubject = "ethereumNodeInstance.create"
	EthereumNodeInstanceUpdateSubject = "ethereumNodeInstance.update"
	EthereumNodeInstanceDeleteSubject = "ethereumNodeInstance.delete"

	PathEthereumNodeDefinitions = "/v0/ethereum_node_definitions"
	PathEthereumNodeInstances   = "/v0/ethereum_node_instances"
)

// GetEthereumNodeDefinitionSubjects returns the NATS subjects
// for ethereum node definitions.
func GetEthereumNodeDefinitionSubjects() []string {
	return []string{
		EthereumNodeDefinitionCreateSubject,
		EthereumNodeDefinitionUpdateSubject,
		EthereumNodeDefinitionDeleteSubject,
	}
}

// GetEthereumNodeInstanceSubjects returns the NATS subjects
// for ethereum node instances.
func GetEthereumNodeInstanceSubjects() []string {
	return []string{
		EthereumNodeInstanceCreateSubject,
		EthereumNodeInstanceUpdateSubject,
		EthereumNodeInstanceDeleteSubject,
	}
}

// GetEthereumNodeSubjects returns the NATS subjects
// for all ethereum node objects.
func GetEthereumNodeSubjects() []string {
	var ethereumNodeSubjects []string

	ethereumNodeSubjects = append(ethereumNodeSubjects, GetEthereumNodeDefinitionSubjects()...)
	ethereumNodeSubjects = append(ethereumNodeSubjects, GetEthereumNodeInstanceSubjects()...)

	return ethereumNodeSubjects
}

// NotificationPayload returns the notification payload that is delivered to the
// controller when a change is made.  It includes the object as presented by the
// client when the change was made.
func (end *EthereumNodeDefinition) NotificationPayload(requeue bool, lastDelay int64) (*[]byte, error) {
	notif := notifications.Notification{
		LastRequeueDelay: &lastDelay,
		Object:           end,
		Requeue:          requeue,
	}

	payload, err := json.Marshal(notif)
	if err != nil {
		return &payload, fmt.Errorf("failed to marshal notification payload %+v: %w", end, err)
	}

	return &payload, nil
}

// GetID returns the unique ID for the object.
func (end *EthereumNodeDefinition) GetID() uint {
	return *end.ID
}

// NotificationPayload returns the notification payload that is delivered to the
// controller when a change is made.  It includes the object as presented by the
// client when the change was made.
func (eni *EthereumNodeInstance) NotificationPayload(requeue bool, lastDelay int64) (*[]byte, error) {
	notif := notifications.Notification{
		LastRequeueDelay: &lastDelay,
		Object:           eni,
		Requeue:          requeue,
	}

	payload, err := json.Marshal(notif)
	if err != nil {
		return &payload, fmt.Errorf("failed to marshal notification payload %+v: %w", eni, err)
	}

	return &payload, nil
}

// GetID returns the unique ID for the object.
func (eni *EthereumNodeInstance) GetID() uint {
	return *eni.ID
}
