/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type RepoSetting struct {
	ConfigPath string `json:"configPath,omitempty"`

	Trusted *bool `json:"trusted,omitempty"`

	Protected *bool `json:"protected,omitempty"`

	Visibility string `json:"visibility,omitempty"`
}

// RepoSpec defines the desired state of Repo
type RepoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Token for administrative access
	Token string `json:"token"`

	// Name of project
	ProjectName string `json:"projectName"`

	// Name of service
	ServiceName string `json:"serviceName"`

	// Owner of service
	Owner string `json:"owner"`

	// Enable target repo or not
	Enable bool `json:"enable"`

	// Setting of repo
	Setting RepoSetting `json:"setting,omitempty"`
}

type RepoResponse struct {
	Id int32 `json:"id,omitempty"`

	UId string `json:"uid,omitempty"`

	UserId int32 `json:"user_id,omitempty"`

	Namespace string `json:"namespace,omitempty"`

	Name string `json:"name,omitempty"`

	Slug string `json:"slug,omitempty"`

	Scm string `json:"scm,omitempty"`

	GitHttpUrl string `json:"git_http_url,omitempty"`

	GitSshUrl string `json:"git_ssh_url,omitempty"`

	Link string `json:"link,omitempty"`

	DefaultBranch string `json:"default_branch,omitempty"`

	Private bool `json:"private,omitempty"`

	Visibility string `json:"visibility,omitempty"`

	Active bool `json:"active,omitempty"`

	ConfigPath string `json:"config_path,omitempty"`

	Trusted *bool `json:"trusted,omitempty"`

	Protected *bool `json:"protected,omitempty"`

	Counter int32 `json:"counter,omitempty"`
}

// RepoStatus defines the observed state of Repo
type RepoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The create time of the repo
	CreateTime *metav1.Time `json:"createTime,omitempty"`

	// The response of the repo
	Response RepoResponse `json:"response,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Repo is the Schema for the repoes API
type Repo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepoSpec   `json:"spec,omitempty"`
	Status RepoStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RepoList contains a list of Repo
type RepoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Repo{}, &RepoList{})
}
