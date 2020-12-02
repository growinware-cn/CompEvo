package internal

import (
	"context"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PodsViaLabels(c client.Client, namespace string, labels labels.Set) ([]v1.Pod, error) {
	pods := &v1.PodList{}
	if err := c.List(context.TODO(), pods, client.InNamespace(namespace),
		client.MatchingLabelsSelector{Selector: labels.AsSelector()}); err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func JobsViaLabels(c client.Client, namespace string, labels labels.Set) ([]batchv1.Job, error) {
	jobs := &batchv1.JobList{}
	if err := c.List(context.TODO(), jobs, client.InNamespace(namespace),
		client.MatchingLabelsSelector{Selector: labels.AsSelector()}); err != nil {
		return nil, err
	}
	return jobs.Items, nil
}
