/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package build

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/wdongyu/builder-manager/syncer"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
)

// buildReconciler reconciles a build object
type BuildReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	KubeClient kubernetes.Interface
}

// +kubebuilder:rbac:groups=apps.builder.io,resources=builds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.builder.io,resources=builds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch;
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/log,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch;delete

func (r *BuildReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("build", req.NamespacedName)

	// your logic here
	log.Printf("Reconciling build %s/%s\n", req.Namespace, req.Name)

	build := &appsv1alpha1.Build{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, build)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	if build.ObjectMeta.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
	}

	if build.Status.CreateTime == nil {
		now := metav1.Now()
		build.Status.CreateTime = &now
	}

	syncers := []syncer.Interface{
		NewJobSyncer(build, r.Client, r.Scheme),
	}

	//if build.ObjectMeta.DeletionTimestamp.IsZero() {
	//	if len(build.ObjectMeta.Finalizers) == 0 {
	//		build.ObjectMeta.Finalizers = append(build.ObjectMeta.Finalizers, StopBuildFinalizer)
	//		err := r.Update(context.TODO(), build)
	//		if err != nil {
	//			return ctrl.Result{}, err
	//		}
	//	}
	//} else {
	//	// Create job to stop the build
	//	syncers = append(syncers, NewDisableJobSyncer(build, r.Client, r.Scheme))
	//}

	if err := r.sync(syncers); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, r.updateStatus(build)
}

func (r *BuildReconciler) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.Recorder); err != nil {
			return err
		}
	}
	return nil
}

func (r *BuildReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Build{}).Owns(&batchv1.Job{}).
		Complete(r)
}
