package handler

import "github.com/wdongyu/builder-manager/api/v1alpha1"

const (
	APP              = "app"
	VERSION          = "version"
	NAMESPACE        = "namespace"
	DefaultNamespace = "default"
)

var (
	paramList = []string{APP, VERSION}
)

type BuildList struct {
	Builds []v1alpha1.Build `json:"builds"`
}

type RepoList struct {
	Repos []v1alpha1.Repo `json:"repos"`
}
