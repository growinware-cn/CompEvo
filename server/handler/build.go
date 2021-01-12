package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/wdongyu/builder-manager/api/v1alpha1"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
)

func (handler *APIHandler) ListBuilds(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	buildList := &v1alpha1.BuildList{}

	err := handler.client.List(context.TODO(), buildList, GetListOptions(vars))
	if err != nil {
		log.Warningf("failed to list builds: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(BuildList{Builds: buildList.Items}, w, http.StatusOK)
	}
}

func (handler *APIHandler) GetBuild(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	var namespace, name string
	namespaceParam, ok := vars[NAMESPACE]
	if !ok || len(namespaceParam) == 0 {
		log.Warningf("Not specify namespace, use default namespace")
		namespace = DefaultNamespace
	} else {
		namespace = namespaceParam[0]
	}

	nameParam, ok := vars[NAME]
	if !ok || len(nameParam) == 0 {
		log.Errorf("Param name is not specify.")
		responseJSON(Message{"Param name is not specify."}, w, http.StatusInternalServerError)
		return
	} else {
		name = nameParam[0]
	}

	build := new(v1alpha1.Build)
	err := handler.client.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, build)
	if err != nil {
		log.Errorf("failed to get build %s/%s : %v", namespace, name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(build, w, http.StatusOK)
	}
}

func (handler *APIHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	build := new(v1alpha1.Build)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &build); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}

	if len(build.Namespace) == 0 {
		build.Namespace = handler.resourcesNamespace
	}

	if err = handler.client.Create(context.TODO(), build); err != nil {
		log.Warningf("Failed to create build %v: %v", build.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(build, w, http.StatusCreated)
	}
}

func (handler *APIHandler) StopBuild(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["build"]

	build := new(v1alpha1.Build)
	build.Name = name

	params := r.URL.Query()
	param, ok := params[NAMESPACE]
	if !ok || len(param) == 0 {
		build.Namespace = handler.resourcesNamespace
	} else {
		build.Namespace = param[0]
	}

	err := handler.client.Delete(context.TODO(), build)
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON("done", w, http.StatusOK)
	}
}
