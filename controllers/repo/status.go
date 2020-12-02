package repo

import (
	"context"
	log "github.com/sirupsen/logrus"
	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
)

func (r *RepoReconciler) updateStatus(repo *appsv1alpha1.Repo, response *appsv1alpha1.RepoResponse) error {
	if repo.Status.CreateTime == nil {
		now := metav1.Now()
		repo.Status.CreateTime = &now
	}

	repo.Status.Response = *response

	return r.syncStatus(repo)
}

func (r *RepoReconciler) syncStatus(repo *appsv1alpha1.Repo) error {
	old := &appsv1alpha1.Repo{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      repo.Name,
		Namespace: repo.Namespace,
	}, old)
	if err != nil {
		log.Errorf("Fail to update the status of repo %s/%s : %v\n", repo.Namespace, repo.Name, err)
		return err
	}

	if !reflect.DeepEqual(old.Status, repo.Status) {
		return r.Status().Update(context.TODO(), repo)
	}

	return nil
}
