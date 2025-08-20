package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplicationSpec struct {
	// +optional
	// +kubebuilder:validation:MaxLength:=512
	Description string `json:"description,omitempty"`
}

type ApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

type Application struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Application
	// +required
	Spec ApplicationSpec `json:"spec"`

	// status defines the observed state of Application
	// +optional
	Status ApplicationStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
