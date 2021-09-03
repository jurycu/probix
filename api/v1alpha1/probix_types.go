/*
Copyright 2021 jurycu.

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

type Result string

const (
	SUCCESS Result = "SUCCESS"
	PENDING Result = "PENDING"
	FAILED  Result = "FAILED"
)

// ProbixSpec defines the desired state of Probix
type ProbixSpec struct {
	//TODO add relabels and labels
	// The job name assigned to scraped metrics by default.

	Targets       []ProbixTarget `json:"targets,omitempty"`
	Interval      string         `json:"interval,omitempty"`
	ScrapeTimeout string         `json:"scrapeTimeout,omitempty"`
}

type ProbixTarget struct {
	//完整的target路径
	MetricsName string `json:"metricsName,omitempty"`
	MetricsHelp string `json:"metricsHelp,omitempty"`

	Target string `json:"target,omitempty"`
	//请求方法,默认为GET
	Method string `json:"method,omitempty"`
	//当请求参数为POST时，可以传入body参数，GET请求只支持path传参
	Body string `json:"body,omitempty"`
}

// ProbixStatus defines the observed state of Probix
type ProbixStatus struct {
	//数据拉取状态
	Status Result `json:"status,omitempty"`
	//备注信息
	Message string `json:"message,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Probix is the Schema for the probixes API
type Probix struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProbixSpec   `json:"spec,omitempty"`
	Status ProbixStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProbixList contains a list of Probix
type ProbixList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Probix `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Probix{}, &ProbixList{})
}
