package repo

import (
	"fmt"
	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newJob(repo *appsv1alpha1.Repo) *batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NameforJob(repo),
			Namespace: repo.Namespace,
			Labels:    LabelsForJob(repo),
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            repo.Spec.ServiceName,
							Image:           DefaultCurlImage,
							ImagePullPolicy: v1.PullIfNotPresent,
							Args:            StringsForEnableRepo(repo),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	}
	return job
}

func newDisableJob(repo *appsv1alpha1.Repo) *batchv1.Job {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NameforDisableJob(repo),
			Namespace: repo.Namespace,
			Labels:    LabelsForDisableJob(repo),
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            repo.Spec.ServiceName,
							Image:           DefaultCurlImage,
							ImagePullPolicy: v1.PullIfNotPresent,
							Args:            StringsForDisableRepo(repo),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	}
	return job
}

func StringsForEnableRepo(repo *appsv1alpha1.Repo) []string {
	strings := []string{"curl", "-s", "-X", "POST", "-H"}

	strings = append(strings, fmt.Sprintf("Authorization: Bearer %s", repo.Spec.Token))

	strings = append(strings, fmt.Sprintf("%s/api/repos/%s/%s",
		DroneServer, repo.Spec.Owner, repo.Spec.ServiceName))

	return strings
}

func StringsForDisableRepo(repo *appsv1alpha1.Repo) []string {
	strings := []string{"curl", "-s", "-X", "DELETE", "-H"}

	strings = append(strings, fmt.Sprintf("Authorization: Bearer %s", repo.Spec.Token))

	strings = append(strings, fmt.Sprintf("%s/api/repos/%s/%s",
		DroneServer, repo.Spec.Owner, repo.Spec.ServiceName))

	return strings
}
