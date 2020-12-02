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

func (handler *APIHandler) ListRepos(w http.ResponseWriter, r *http.Request) {
	repoList := &v1alpha1.RepoList{}

	err := handler.client.List(context.TODO(), repoList, &client.ListOptions{})

	if err != nil {
		log.Warningf("failed to list builds: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(RepoList{Repos: repoList.Items}, w, http.StatusOK)
	}
}

func (handler *APIHandler) EnableRepo(w http.ResponseWriter, r *http.Request) {
	repo := new(v1alpha1.Repo)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &repo); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}

	if len(repo.Namespace) == 0 {
		repo.Namespace = handler.resourcesNamespace
	}

	log.Printf("Receive repo : %v", repo)
	if err = handler.client.Create(context.TODO(), repo); err != nil {
		log.Warningf("Failed to enable repo %v: %v", repo.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(repo, w, http.StatusCreated)
	}
}

func (handler *APIHandler) DisableRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["repo"]

	repo := new(v1alpha1.Repo)
	repo.Name = name
	repo.Namespace = handler.resourcesNamespace

	err := handler.client.Delete(context.TODO(), repo)
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON("", w, http.StatusOK)
	}
}
