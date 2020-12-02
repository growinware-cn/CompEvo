package build

import (
	"fmt"
	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"strings"
)

func NameforJob(build *appsv1alpha1.Build) string {
	return fmt.Sprintf("%s-%d", build.Spec.ServiceName, build.Status.CreateTime.Unix())
}

func NameforDisableJob(build *appsv1alpha1.Build) string {
	return fmt.Sprintf("%s-%s-%d", strings.ToLower(StopBuildFinalizer), build.Spec.ServiceName, build.Status.CreateTime.Unix())
}

func UIDForManagedService(build *appsv1alpha1.Build) string {
	//return build.Spec.ServiceName
	return fmt.Sprintf("%s-%s", build.Spec.ProjectName, build.Spec.ServiceName)
}

func LabelsForJob(build *appsv1alpha1.Build) labels.Set {
	l := labels.Set{}

	l[LabelOwnerKey] = build.Kind
	l[LabelProjectKey] = build.Spec.ProjectName
	l[LabelServiceKey] = build.Spec.ServiceName
	l[LabelCreateTime] = fmt.Sprintf("%d", build.Status.CreateTime.Unix())
	l[StopBuildFinalizer] = "false"

	return l
}

func LabelsForDisableJob(build *appsv1alpha1.Build) labels.Set {
	l := LabelsForJob(build)
	l[StopBuildFinalizer] = "true"

	return l
}

func LabelsForPod(build *appsv1alpha1.Build) labels.Set {
	l := labels.Set{}
	l[LabelJobName] = NameforJob(build)

	return l
}

func LabelsForDisablePod(build *appsv1alpha1.Build) labels.Set {
	l := labels.Set{}
	l[LabelJobName] = NameforDisableJob(build)

	return l
}
