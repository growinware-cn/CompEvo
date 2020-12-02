package syncer

import "k8s.io/apimachinery/pkg/runtime"

type MutateFn func(existing runtime.Object) error

// OperationResult is the action result of a CreateOrUpdate call
type OperationResult string

const ( // They should complete the sentence "Deployment default/foo has been ..."
	// OperationResultNone means that the resource has not been changed
	OperationResultNone OperationResult = "unchanged"
	// OperationResultCreated means that a new resource is created
	OperationResultCreated OperationResult = "created"
	// OperationResultUpdated means that an existing resource is updated
	OperationResultUpdated OperationResult = "updated"
)

const (
	eventNormal  = "Normal"
	eventWarning = "Warning"
)
