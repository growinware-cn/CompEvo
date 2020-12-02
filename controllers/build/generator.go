package build

import (
	"fmt"
	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newJob(build *appsv1alpha1.Build) *batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NameforJob(build),
			Namespace: build.Namespace,
			Labels:    LabelsForJob(build),
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            build.Spec.ServiceName,
							Image:           DefaultImage,
							ImagePullPolicy: v1.PullIfNotPresent,
							Args:            ArgsForCreateBuild(build),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	}
	return job
}

func newDisabledJob(build *appsv1alpha1.Build) *batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NameforDisableJob(build),
			Namespace: build.Namespace,
			Labels:    LabelsForDisableJob(build),
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            build.Spec.ServiceName,
							Image:           DefaultCurlImage,
							ImagePullPolicy: v1.PullIfNotPresent,
							Args:            StringsForStopBuild(build),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	}
	return job
}

func ArgsForCreateBuild(build *appsv1alpha1.Build) []string {
	var strings []string

	strings = append(strings, fmt.Sprintf("--TOKEN=%s", build.Spec.Token))

	postUrl := fmt.Sprintf("--POST_URL=%s/api/repos/%s/%s/builds?%s_%s=%s",
		DroneServer, build.Spec.Owner, build.Spec.ServiceName, AliyunLogsPrefix, UIDForManagedService(build), STDOUT)
	if build.Spec.Branch != "" {
		postUrl += fmt.Sprintf("&branch=%s", build.Spec.Branch)
	}
	if build.Spec.Commit != "" {
		postUrl += fmt.Sprintf("&commit=%s", build.Spec.Commit)
	}

	strings = append(strings, postUrl)

	strings = append(strings, fmt.Sprintf("--GET_URL=%s/api/repos/%s/%s/builds/",
		DroneServer, build.Spec.Owner, build.Spec.ServiceName))

	return strings
}

func StringsForCreateBuild(build *appsv1alpha1.Build) []string {
	strings := []string{"curl", "-s", "-X", "POST", "-H"}

	strings = append(strings, fmt.Sprintf("Authorization: Bearer %s", build.Spec.Token))

	if build.Spec.Branch == "" {
		// Creates a build using the latest commit to the default branch (master)
		strings = append(strings, fmt.Sprintf("%s/api/repos/%s/%s/builds?%s_%s=%s",
			DroneServer, build.Spec.Owner, build.Spec.ServiceName, AliyunLogsPrefix, UIDForManagedService(build), STDOUT))
	} else if build.Spec.Commit == "" {
		// Creates a build using the latest commit to the named branch
		strings = append(strings, fmt.Sprintf("%s/api/repos/%s/%s/builds?branch=%s&%s_%s=%s",
			DroneServer, build.Spec.Owner, build.Spec.ServiceName, build.Spec.Branch, AliyunLogsPrefix, UIDForManagedService(build), STDOUT))
	} else {
		// Creates a build using the named branch and commit sha
		strings = append(strings, fmt.Sprintf("%s/api/repos/%s/%s/builds?branch=%s&commit=%s&%s_%s=%s",
			DroneServer, build.Spec.Owner, build.Spec.ServiceName, build.Spec.Branch, build.Spec.Commit, AliyunLogsPrefix, UIDForManagedService(build), STDOUT))
	}

	return strings
}

func StringsForStopBuild(build *appsv1alpha1.Build) []string {
	strings := []string{"curl", "-s", "-X", "DELETE", "-H"}

	strings = append(strings, fmt.Sprintf("Authorization: Bearer %s", build.Spec.Token))

	strings = append(strings, fmt.Sprintf("%s/api/repos/%s/%s/builds/%d",
		DroneServer, build.Spec.Owner, build.Spec.ServiceName, build.Status.Response.Number))

	return strings
}
