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

// buildSpec defines the desired state of build
type BuildSpec struct {
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

	// Branch
	Branch string `json:"branch,omitempty"`

	// SHA commit info
	Commit string `json:"commit,omitempty"`

	// Number of the target build
	// BuildNo int32 `json:"buildNo,omitempty"`

	// Number of the target step in certain build
	// StepNo int32 `json:"StepNo,omitempty"`
}

type BuildPhase string

const (
	PhaseSuccess = "Succeeded"

	PhaseRunning = "Running"
)

type BuildResponse struct {
	Id int32 `json:"id,omitempty"`

	RepoId int32 `json:"repo_id,omitempty"`

	Number int32 `json:"number,omitempty"`

	Message string `json:"message,omitempty"`

	Before string `json:"before,omitempty"`

	After string `json:"after,omitempty"`
}

// buildStatus defines the observed state of build
type BuildStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The create time of the build
	CreateTime *metav1.Time `json:"createTime,omitempty"`

	// The phase of the build
	RequestPhase BuildPhase `json:"requestPhase,omitempty"`

	// The response of the build request
	Response BuildResponse `json:"response,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// build is the Schema for the builds API
type Build struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BuildSpec   `json:"spec,omitempty"`
	Status BuildStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// buildList contains a list of build
type BuildList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Build `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Build{}, &BuildList{})
}
