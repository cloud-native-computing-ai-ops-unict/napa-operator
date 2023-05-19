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

// Objective represents single threshold for SLO, for internal usage.
type Objective struct {
	DisplayName     string `json:"displayName,omitempty"`
	Op              string `json:"op,omitempty" example:"lte"`
	Value           string `json:"value,omitempty" validate:"numeric,omitempty"`
	Target          string `json:"target" validate:"required,numeric,gte=0,lt=1" example:"0.9"`
	TimeSliceTarget string `json:"timeSliceTarget,omitempty" validate:"gte=0,lte=1,omitempty" example:"0.9"`
	TimeSliceWindow string `json:"timeSliceWindow,omitempty" example:"5m"`
}

// Calendar struct represents calendar time window.
type Calendar struct {
	StartTime string `json:"startTime" validate:"required,dateWithTime" example:"2020-01-21 12:30:00"`
	TimeZone  string `json:"timeZone" validate:"required,timeZone" example:"America/New_York"`
}

// TimeWindow represents content of time window.
type TimeWindow struct {
	Duration  string    `json:"duration" validate:"required,validDuration" example:"1h"`
	IsRolling bool      `json:"isRolling" example:"true"`
	Calendar  *Calendar `json:"calendar,omitempty" validate:"required_if=IsRolling false"`
}

// SLOSpec defines the desired state of SLO
type SLOSpec struct {
	Description     string       `json:"description,omitempty" validate:"max=1050,omitempty"`
	Service         string       `json:"service" validate:"required" example:"webapp-service"`
	Indicator       *SLIInline   `json:"indicator,omitempty" validate:"required_without=IndicatorRef"`
	IndicatorRef    *string      `json:"indicatorRef,omitempty"`
	BudgetingMethod string       `json:"budgetingMethod" validate:"required,oneof=Occurrences Timeslices" example:"Occurrences"` //nolint:lll
	TimeWindow      []TimeWindow `json:"timeWindow" validate:"required,len=1,dive"`
	Objectives      []Objective  `json:"objectives" validate:"required,dive"`
	// We don't make clear in the spec if this is a ref or inline.
	// We will make it a ref for now.
	// https://github.com/OpenSLO/OpenSLO/issues/133
	AlertPolicies []string `json:"alertPolicies" validate:"dive"`
}

// SLOStatus defines the observed state of SLO
type SLOStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SLO is the Schema for the sloes API
type SLO struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SLOSpec   `json:"spec,omitempty"`
	Status SLOStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SLOList contains a list of SLO
type SLOList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SLO `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SLO{}, &SLOList{})
}
