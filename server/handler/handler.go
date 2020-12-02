package handler

import (
	"encoding/json"
	"github.com/wdongyu/builder-manager/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type APIHandler struct {
	resourcesNamespace string
	client             client.Client
}

func NewAPIHandler() (*APIHandler, error) {
	kubeConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := setupClient(kubeConfig)
	if err != nil {
		return nil, err
	}

	apiHandler := &APIHandler{
		client: clientSet,
	}

	apiHandler.resourcesNamespace = os.Getenv("RESOURCES_NAMESPACE")
	if len(apiHandler.resourcesNamespace) == 0 {
		apiHandler.resourcesNamespace = DefaultNamespace
	}

	return apiHandler, nil
}

func setupClient(config *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	for _, addToSchemeFunc := range []func(s *runtime.Scheme) error{
		v1alpha1.AddToScheme,
		v1.AddToScheme,
		batchv1.AddToScheme,
	} {
		if err := addToSchemeFunc(scheme); err != nil {
			return nil, err
		}
	}

	clientSet, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

type Message struct {
	Message string `json:"message"`
}

func responseJSON(body interface{}, w http.ResponseWriter, statusCode int) {
	jsonResponse, err := json.Marshal(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
