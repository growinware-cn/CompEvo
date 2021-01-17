package build

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
	"github.com/wdongyu/builder-manager/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"strings"
)

func (r *BuildReconciler) updateStatus(build *appsv1alpha1.Build) error {
	if build.Status.CreateTime == nil {
		now := metav1.Now()
		build.Status.CreateTime = &now
	}

	jobs, err := internal.JobsViaLabels(r.ApiReader, build.Namespace, LabelsForJob(build))
	if err != nil {
		log.Errorf("Fail to list jobs with labels : %v", LabelsForJob(build))
		return err
	}

	if len(jobs) == 0 {
		log.Info("Job has not been created")
		return nil
	}

	if jobs[0].Status.Succeeded != 1 {
		log.Infof("Job %s/%s is running", jobs[0].Namespace, jobs[0].Name)
		build.Status.RequestPhase = appsv1alpha1.PhaseRunning
	} else {
		log.Infof("Job %s/%s is completed", jobs[0].Namespace, jobs[0].Name)

		pods, err := internal.PodsViaLabels(r.ApiReader, build.Namespace, LabelsForPod(build))
		if err != nil {
			log.Errorf("Fail to get pod for labels %v : %v", LabelsForPod(build), err)
			return err
		}

		podLogs, err := internal.GetLog(r.KubeClient, build.Namespace, pods[0].Name)
		if err != nil {
			log.Errorf("Fail to get pod log : %v", err)
			return err
		}
		left := strings.Index(podLogs, "{")
		right := strings.LastIndex(podLogs, "}")
		if left == -1 || right == -1 {
			build.Status.RequestPhase = appsv1alpha1.PhaseFailed
		} else if err = json.Unmarshal([]byte(podLogs[left:right+1]), &build.Status.Response); err != nil {
			log.Errorf("Fail to decode pod log to build response : %v", err)
			build.Status.RequestPhase = appsv1alpha1.PhaseFailed
		} else {
			build.Status.RequestPhase = appsv1alpha1.PhaseSuccess
		}
	}

	if err = r.syncStatus(build); err != nil {
		return err
	}

	//if !build.ObjectMeta.DeletionTimestamp.IsZero() {
	//	disableJobs, err := internal.JobsViaLabels(r.Client, build.Namespace, LabelsForDisableJob(build))
	//	if err != nil {
	//		log.Errorf("Fail to list jobs with labels : %v", LabelsForDisableJob(build))
	//		return err
	//	}
	//
	//	if len(disableJobs) != 0 && disableJobs[0].Status.Succeeded == 1 {
	//		pods, err := internal.PodsViaLabels(r.Client, build.Namespace, LabelsForDisablePod(build))
	//		if err != nil {
	//			log.Errorf("Fail to get pod for labels %v : %v", LabelsForDisablePod(build), err)
	//			return err
	//		}
	//
	//		podLogs, err := internal.GetLog(r.KubeClient, build.Namespace, pods[0].Name)
	//		if err != nil {
	//			log.Errorf("Fail to get pod log : %v", err)
	//			return err
	//		}
	//		log.Info(podLogs)
	//
	//		build.ObjectMeta.Finalizers = []string{}
	//		if err := r.Update(context.TODO(), build); err != nil {
	//			log.Errorf("Fail to update build's finalizers : %v", err)
	//			return err
	//		}
	//	} else if len(disableJobs) != 0 {
	//		log.Infof("DisableJob %s/%s is running", disableJobs[0].Namespace, disableJobs[0].Name)
	//	}
	//}

	return nil
}

func (r *BuildReconciler) syncStatus(build *appsv1alpha1.Build) error {
	old := &appsv1alpha1.Build{}
	err := r.ApiReader.Get(context.TODO(), types.NamespacedName{
		Name:      build.Name,
		Namespace: build.Namespace,
	}, old)
	if err != nil {
		log.Errorf("Fail to update the status of build %s/%s : %v\n", build.Namespace, build.Name, err)
		return err
	}

	if !reflect.DeepEqual(old.Status, build.Status) {
		//return r.Client.Update(context.TODO(), bm)
		return r.Status().Update(context.TODO(), build)
	}

	return nil
}
