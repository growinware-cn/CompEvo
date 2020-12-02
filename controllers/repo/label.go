package repo

import (
	"fmt"
	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

func NameforJob(repo *appsv1alpha1.Repo) string {
	return fmt.Sprintf("%s-%s-%s", strings.ToLower(repo.Kind), repo.Spec.ProjectName, repo.Spec.ServiceName)
}

func NameforDisableJob(repo *appsv1alpha1.Repo) string {
	return fmt.Sprintf("%s.%s-%s-%s", strings.ToLower(repo.Kind), strings.ToLower(DisableRepoFinalizer), repo.Spec.ProjectName, repo.Spec.ServiceName)
}

func LabelsForJob(repo *appsv1alpha1.Repo) labels.Set {
	l := labels.Set{}

	l[LabelOwnerKey] = repo.Kind
	l[LabelProjectKey] = repo.Spec.ProjectName
	l[LabelServiceKey] = repo.Spec.ServiceName
	l[DisableRepoFinalizer] = "false"

	return l
}

func LabelsForDisableJob(repo *appsv1alpha1.Repo) labels.Set {
	l := LabelsForJob(repo)
	l[DisableRepoFinalizer] = "true"

	return l
}

func LabelsForPod(repo *appsv1alpha1.Repo) labels.Set {
	l := labels.Set{}
	l[LabelJobName] = NameforJob(repo)

	return l
}

func LabelsForDisablePod(repo *appsv1alpha1.Repo) labels.Set {
	l := labels.Set{}
	l[LabelJobName] = NameforDisableJob(repo)

	return l
}
