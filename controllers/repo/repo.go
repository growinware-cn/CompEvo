package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	appsv1alpha1 "github.com/wdongyu/builder-manager/api/v1alpha1"
	"io/ioutil"
	"net/http"
)

func GetRepoInfo(repo *appsv1alpha1.Repo) (*appsv1alpha1.RepoResponse, error) {
	requestUrl := fmt.Sprintf("%s/api/repos/%s/%s", DroneServer, repo.Spec.Owner, repo.Spec.ServiceName)
	request, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repo.Spec.Token))

	body, err := sendRequest(request)
	if err != nil {
		return nil, err
	}
	repoResp := &appsv1alpha1.RepoResponse{}
	if err = json.Unmarshal(body, &repoResp); err != nil {
		log.Errorf("Fail to unmarshal response body : %v", err)
		return nil, err
	}

	return repoResp, nil
}

func EnableRepo(repo *appsv1alpha1.Repo) (*appsv1alpha1.RepoResponse, error) {
	requestUrl := fmt.Sprintf("%s/api/repos/%s/%s", DroneServer, repo.Spec.Owner, repo.Spec.ServiceName)
	request, err := http.NewRequest(http.MethodPost, requestUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repo.Spec.Token))

	body, err := sendRequest(request)
	if err != nil {
		return nil, err
	}
	repoResp := &appsv1alpha1.RepoResponse{}
	if err = json.Unmarshal(body, &repoResp); err != nil {
		log.Errorf("Fail to unmarshal response body : %v", err)
		return nil, err
	}

	return repoResp, nil
}

func DisableRepo(repo *appsv1alpha1.Repo) (*appsv1alpha1.RepoResponse, error) {
	requestUrl := fmt.Sprintf("%s/api/repos/%s/%s", DroneServer, repo.Spec.Owner, repo.Spec.ServiceName)
	request, err := http.NewRequest(http.MethodDelete, requestUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repo.Spec.Token))

	body, err := sendRequest(request)
	if err != nil {
		return nil, err
	}
	repoResp := &appsv1alpha1.RepoResponse{}
	if err = json.Unmarshal(body, &repoResp); err != nil {
		log.Errorf("Fail to unmarshal response body : %v", err)
		return nil, err
	}

	return repoResp, nil
}

func UpdateRepo(repo *appsv1alpha1.Repo) (*appsv1alpha1.RepoResponse, error) {
	requestUrl := fmt.Sprintf("%s/api/repos/%s/%s", DroneServer, repo.Spec.Owner, repo.Spec.ServiceName)
	requestBody := new(bytes.Buffer)
	repoRequest := &appsv1alpha1.RepoResponse{
		Visibility: repo.Spec.Setting.Visibility,
		ConfigPath: repo.Spec.Setting.ConfigPath,
		Trusted:    repo.Spec.Setting.Trusted,
		Protected:  repo.Spec.Setting.Protected,
	}
	err := json.NewEncoder(requestBody).Encode(repoRequest)
	if err != nil {
		log.Errorf("Fail to encode repo setting : %v", err)
		return nil, err
	}
	log.Printf("Encode json : %v", requestBody)
	request, err := http.NewRequest(http.MethodPatch, requestUrl, requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repo.Spec.Token))
	request.Header.Set("Content-Type", "application/json")

	body, err := sendRequest(request)
	if err != nil {
		return nil, err
	}
	repoResp := &appsv1alpha1.RepoResponse{}
	if err = json.Unmarshal(body, &repoResp); err != nil {
		log.Errorf("Fail to unmarshal response body : %v", err)
		return nil, err
	}

	return repoResp, nil
}

func sendRequest(request *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Errorf("Fail to get target repo : %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Fail to read response body : %v", err)
		return nil, err
	}

	return body, nil
}

func IsMatch(setting appsv1alpha1.RepoSetting, response appsv1alpha1.RepoResponse) bool {
	if setting.ConfigPath != response.ConfigPath ||
		*setting.Protected != *response.Protected ||
		*setting.Trusted != *response.Trusted ||
		setting.Visibility != response.Visibility {
		return false
	}

	return true
}

func (r *RepoReconciler) UpdateSecret(repo *appsv1alpha1.Repo) error {
	usernameRequestUrl := fmt.Sprintf("%s/api/repos/%s/%s/secrets/%s",
		DroneServer, repo.Spec.Owner, repo.Spec.ServiceName, DockerUsername)
	request, err := http.NewRequest(http.MethodGet, usernameRequestUrl, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repo.Spec.Token))

	body, err := sendRequest(request)
	if err != nil {
		return err
	}

	secretResponse := &SecretResponse{}
	if err = json.Unmarshal(body, &secretResponse); err != nil {
		return err
	}

	return nil
}
