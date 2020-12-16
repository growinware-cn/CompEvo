package handler

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetListOptions(vars url.Values) *client.ListOptions {
	paramMap := make(map[string]string)
	for _, paramName := range paramList {
		p, ok := vars[paramName]
		if !ok || len(p) == 0 {

		} else {
			paramMap[paramName] = p[0]
		}
	}

	var targetNamespace string
	namespace, ok := vars[NAMESPACE]
	if !ok || len(namespace) == 0 {
		log.Warningf("Not specify namespace, use default namespace")
		targetNamespace = DefaultNamespace
	} else {
		targetNamespace = namespace[0]
	}

	selector := &metav1.LabelSelector{
		MatchLabels: paramMap,
	}
	labelSelector, _ := metav1.LabelSelectorAsSelector(selector)

	listOptions := &client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     targetNamespace,
	}

	log.Printf("List options result : %v", listOptions)
	return listOptions
}
