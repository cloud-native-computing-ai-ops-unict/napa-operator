/*
Copyright 2023.

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

type SLIInline struct {
	Metadata metav1.ObjectMeta `json:"metadata" validate:"required"`
	Spec     SLISpec           `json:"spec" validate:"required"`
}

// RatioMetric represents the ratio metric.
type RatioMetric struct {
	Counter bool                `json:"counter" example:"true"`
	Good    *MetricSourceHolder `json:"good,omitempty" validate:"required_without=Bad"`
	Bad     *MetricSourceHolder `json:"bad,omitempty" validate:"required_without=Good"`
	Total   MetricSourceHolder  `json:"total" validate:"required"`
}

// MetricSourceHolder represents the metric source holder.
type MetricSourceHolder struct {
	MetricSource MetricSource `json:"metricSource" validate:"required"`
}

// MetricSource represents the metric source.
type MetricSource struct {
	MetricSourceRef  string            `json:"metricSourceRef,omitempty" validate:"required_without=MetricSourceSpec"`
	Type             string            `json:"type,omitempty" validate:"required_without=MetricSourceRef"`
	MetricSourceSpec map[string]string `json:"spec" validate:"required_without=MetricSourceRef"`
}

// SLISpec defines the desired state of SLI
type SLISpec struct {
	ThresholdMetric *MetricSourceHolder `json:"thresholdMetric,omitempty" validate:"required_without=RatioMetric"`
	RatioMetric     *RatioMetric        `json:"ratioMetric,omitempty" validate:"required_without=ThresholdMetric"`
}

// SLIStatus defines the observed state of SLI
type SLIStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SLI is the Schema for the sli API
type SLI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SLISpec   `json:"spec,omitempty"`
	Status SLIStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SLIList contains a list of SLI
type SLIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SLI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SLI{}, &SLIList{})
}
