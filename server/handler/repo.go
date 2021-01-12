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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListRepos(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	repoList := &v1alpha1.RepoList{}

	err := handler.client.List(context.TODO(), repoList, GetListOptions(vars))

	if err != nil {
		log.Warningf("failed to list builds: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(RepoList{Repos: repoList.Items}, w, http.StatusOK)
	}
}

func (handler *APIHandler) GetRepo(w http.ResponseWriter, r *http.Request) {
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

	repo := new(v1alpha1.Repo)
	err := handler.client.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, repo)
	if err != nil {
		log.Errorf("failed to get repo %s/%s : %v", namespace, name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(repo, w, http.StatusOK)
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

	params := r.URL.Query()
	param, ok := params[NAMESPACE]
	if !ok || len(param) == 0 {
		repo.Namespace = handler.resourcesNamespace
	} else {
		repo.Namespace = param[0]
	}

	err := handler.client.Delete(context.TODO(), repo)
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON("done", w, http.StatusOK)
	}
}

func (handler *APIHandler) UpdateRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["repo"]

	repo := new(v1alpha1.Repo)
	repo.Name = name

	params := r.URL.Query()
	param, ok := params[NAMESPACE]
	if !ok || len(param) == 0 {
		repo.Namespace = handler.resourcesNamespace
	} else {
		repo.Namespace = param[0]
	}

	patchRepo := new(v1alpha1.RepoSetting)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &patchRepo); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}

	patchData := []PatchInterfaceValue{}
	if patchRepo.ConfigPath != "" {
		patchData = append(patchData, PatchInterfaceValue{
			Op:    "replace",
			Path:  "/spec/setting/configPath",
			Value: patchRepo.ConfigPath,
		})
	}
	if patchRepo.Trusted != nil {
		patchData = append(patchData, PatchInterfaceValue{
			Op:    "replace",
			Path:  "/spec/setting/trusted",
			Value: patchRepo.Trusted,
		})
	}
	if patchRepo.Protected != nil {
		patchData = append(patchData, PatchInterfaceValue{
			Op:    "replace",
			Path:  "/spec/setting/protected",
			Value: patchRepo.Protected,
		})
	}
	if patchRepo.Visibility != "" {
		patchData = append(patchData, PatchInterfaceValue{
			Op:    "replace",
			Path:  "/spec/setting/visibility",
			Value: patchRepo.Visibility,
		})
	}

	payLoad, _ := json.Marshal(patchData)
	rawPatch := client.RawPatch(types.JSONPatchType, payLoad)
	if err = handler.client.Patch(context.TODO(), repo, rawPatch); err != nil {
		log.Errorf("Fail to patch repo %s/%s : %v", repo.Namespace, repo.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON("done", w, http.StatusOK)
	}
}
