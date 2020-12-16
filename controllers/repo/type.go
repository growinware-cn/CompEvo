package repo

const (
	DefaultCurlImage = "pstauffer/curl"

	LabelOwnerKey   = "owner"
	LabelProjectKey = "project"
	LabelServiceKey = "service"
	LabelJobName    = "job-name"

	DisableRepoFinalizer = "disable-repo"

	//DroneServer = "http://114.212.82.229:30170"
	DroneServer = "http://192.168.80.133:30170"

	DockerUsername = "username"
	DockerPassword = "password"
)

type SecretRequest struct {
	Name        string `json:"name"`
	Data        string `json:"data"`
	PullRequest bool   `json:"pull_request"`
}

type SecretResponse struct {
	Id     int32  `json:"id"`
	RepoId int32  `json:"repo_id"`
	Name   string `json:"name"`
}
