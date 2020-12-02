package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/wdongyu/builder-manager/api/v1alpha1"
	"io"
	"io/ioutil"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListBuilds(w http.ResponseWriter, r *http.Request) {
	buildList := &v1alpha1.BuildList{}

	err := handler.client.List(context.TODO(), buildList, &client.ListOptions{})

	if err != nil {
		log.Warningf("failed to list builds: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(BuildList{Builds: buildList.Items}, w, http.StatusOK)
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
	build.Namespace = handler.resourcesNamespace

	err := handler.client.Delete(context.TODO(), build)
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON("", w, http.StatusOK)
	}
}
