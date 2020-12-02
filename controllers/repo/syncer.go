package repo

import (
	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
	"github.com/wdongyu/builder-manager/syncer"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewJobSyncer(repo *appsv1alpha1.Repo, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	job := newJob(repo)
	return syncer.NewObjectSyncer("job", repo, job, c, scheme, func(existing runtime.Object) error {
		out := existing.(*batchv1.Job)
		if !reflect.DeepEqual(out.Spec, job.Spec) {
			out.Spec = job.Spec
		}
		return nil
	})
}

func NewDisableJobSyncer(repo *appsv1alpha1.Repo, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	job := newDisableJob(repo)
	return syncer.NewObjectSyncer("disableJob", repo, job, c, scheme, func(existing runtime.Object) error {
		out := existing.(*batchv1.Job)
		if !reflect.DeepEqual(out.Spec, job.Spec) {
			out.Spec = job.Spec
		}
		return nil
	})
}
