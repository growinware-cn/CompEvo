package syncer

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"reflect"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ObjectSyncer is a syncer.Interface for syncing kubernetes.Objects
type ObjectSyncer struct {
	Owner          runtime.Object
	Obj            runtime.Object
	Name           string
	Client         client.Client
	Scheme         *runtime.Scheme
	SyncFn         MutateFn
	previousObject runtime.Object
}

// GetObject returns the ObjectSyncer subject
func (s *ObjectSyncer) GetObject() interface{} { return s.Obj }

// GetOwner returns the ObjectSyncer owner
func (s *ObjectSyncer) GetOwner() runtime.Object { return s.Owner }

// NewObjectSyncer creates a new kubernetes.Object syncer for a given object
// with an owner and persists data using controller-runtime's CreateOrUpdate.
// The name is used for logging and event emitting purposes and should be an
// valid go identifier in upper camel case. (eg. MysqlStatefulSet)
func NewObjectSyncer(name string, owner, obj runtime.Object, c client.Client, scheme *runtime.Scheme, syncFn MutateFn) Interface {
	return &ObjectSyncer{
		Owner:  owner,
		Obj:    obj,
		Name:   name,
		SyncFn: syncFn,
		Client: c,
		Scheme: scheme,
	}
}

// Sync does the actual syncing and implements the syncer.Interface Sync method
func (s *ObjectSyncer) Sync(ctx context.Context) (SyncResult, error) {
	result := SyncResult{}

	key, err := getNameAndNamespace(s.Obj)
	if err != nil {
		return result, err
	}

	result.Operation, err = CreateOrUpdate(ctx, s.Client, s.Obj, s.mutateFn())

	if err != nil {
		result.SetEventData(eventWarning, basicEventReason(s.Name, err),
			fmt.Sprintf("%T %s failed syncing: %s", s.Obj, key, err))
		log.Errorf("%s: key %s, kind: %s, err: %v", string(result.Operation), key, fmt.Sprintf("%T", s.Obj), err)
	} else {
		result.SetEventData(eventNormal, basicEventReason(s.Name, err),
			fmt.Sprintf("%T %s %s successfully", s.Obj, key, result.Operation))
	}

	return result, err
}

// CreateOrUpdate creates or updates the given object obj in the Kubernetes
// cluster. The object's desired state should be reconciled with the existing
// state using the passed in ReconcileFn. obj must be a struct pointer so that
// obj can be updated with the content returned by the Server.
//
// It returns the executed operation and an error.
func CreateOrUpdate(ctx context.Context, c client.Client, obj runtime.Object, f MutateFn) (OperationResult, error) {
	// op is the operation we are going to attempt
	op := OperationResultNone

	// get the existing object meta
	metaObj, ok := obj.(metav1.Object)
	if !ok {
		return OperationResultNone, fmt.Errorf("%T does not implement metav1.Object interface", obj)
	}

	// retrieve the existing object
	key := client.ObjectKey{
		Name:      metaObj.GetName(),
		Namespace: metaObj.GetNamespace(),
	}
	err := c.Get(ctx, key, obj)

	// reconcile the existing object
	existing := obj.DeepCopyObject()
	existingObjMeta := existing.(metav1.Object)
	existingObjMeta.SetName(metaObj.GetName())
	existingObjMeta.SetNamespace(metaObj.GetNamespace())

	if e := f(obj); e != nil {
		return OperationResultNone, e
	}

	if metaObj.GetName() != existingObjMeta.GetName() {
		return OperationResultNone, fmt.Errorf("ReconcileFn cannot mutate objects name")
	}

	if metaObj.GetNamespace() != existingObjMeta.GetNamespace() {
		return OperationResultNone, fmt.Errorf("ReconcileFn cannot mutate objects namespace")
	}

	if errors.IsNotFound(err) {
		err = c.Create(ctx, obj)
		op = OperationResultCreated
	} else if err == nil {
		if reflect.DeepEqual(existing, obj) {
			return OperationResultNone, nil
		}
		err = c.Update(ctx, obj)
		op = OperationResultUpdated
	} else {
		return OperationResultNone, err
	}

	if err != nil {
		op = OperationResultNone
	}
	return op, err
}

// Given an ObjectSyncer, returns a controllerutil.MutateFn which also sets the
// owner reference if the subject has one
func (s *ObjectSyncer) mutateFn() MutateFn {
	return func(existing runtime.Object) error {
		s.previousObject = existing.DeepCopyObject()
		err := s.SyncFn(existing)
		if err != nil {
			return err
		}
		if s.Owner != nil {
			existingMeta, ok := existing.(metav1.Object)
			if !ok {
				return fmt.Errorf("%T is not a metav1.Object", existing)
			}
			ownerMeta, ok := s.Owner.(metav1.Object)
			if !ok {
				return fmt.Errorf("%T is not a metav1.Object", s.Owner)
			}
			err := controllerutil.SetControllerReference(ownerMeta, existingMeta, s.Scheme)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
