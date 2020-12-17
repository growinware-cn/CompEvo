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

package repo

import (
	"context"
	"github.com/go-logr/logr"
	log "github.com/sirupsen/logrus"
	"github.com/wdongyu/builder-manager/syncer"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
)

// RepoReconciler reconciles a Repo object
type RepoReconciler struct {
	client.Client
	ApiReader  client.Reader
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	KubeClient kubernetes.Interface
}

// +kubebuilder:rbac:groups=apps.builder.io,resources=repoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.builder.io,resources=repoes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch;
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/log,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch;delete

func (r *RepoReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("repo", req.NamespacedName)

	// your logic here
	log.Printf("Reconciling repo %s/%s\n", req.Namespace, req.Name)

	repo := &appsv1alpha1.Repo{}
	err := r.ApiReader.Get(context.TODO(), req.NamespacedName, repo)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	response := &appsv1alpha1.RepoResponse{}

	if repo.ObjectMeta.DeletionTimestamp.IsZero() {
		if len(repo.ObjectMeta.Finalizers) == 0 {
			repo.ObjectMeta.Finalizers = append(repo.ObjectMeta.Finalizers, DisableRepoFinalizer)
			err := r.Update(context.TODO(), repo)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		// Get repo info
		repoInfo, err := GetRepoInfo(repo)
		if err != nil {
			return ctrl.Result{}, err
		}
		response = repoInfo.DeepCopy()
		log.Printf("Reconciling current repo info : %v", repoInfo)

		if repo.Spec.Enable && !repoInfo.Active {
			log.Printf("Enable the target repo")
			response, err = EnableRepo(repo)
			if err != nil {
				return ctrl.Result{}, err
			}
		} else if !repo.Spec.Enable && repoInfo.Active {
			log.Printf("Disable the target repo")
			response, err = DisableRepo(repo)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		if !IsMatch(repo.Spec.Setting, *repoInfo) {
			log.Printf("Update the target repo")
			response, err = UpdateRepo(repo)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		log.Printf("Finalize : Disable the target repo")
		response, err = DisableRepo(repo)
		if err != nil {
			return ctrl.Result{}, err
		}
		repo.ObjectMeta.Finalizers = []string{}
		if err := r.Update(context.TODO(), repo); err != nil {
			log.Errorf("Fail to update repo's finalizers : %v", err)
			return ctrl.Result{}, err
		}
	}

	//if err = r.UpdateSecret(repo); err != nil {
	//	return ctrl.Result{}, err
	//}

	return ctrl.Result{}, r.updateStatus(repo, response)
}

func (r *RepoReconciler) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.Recorder); err != nil {
			return err
		}
	}
	return nil
}

func (r *RepoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Repo{}).Owns(&batchv1.Job{}).
		Complete(r)
}
