package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:method=GetApplication,verb=get
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Application struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`
	// +required
	Spec ApplicationSpec `json:"spec"`
	// +optional
	Status ApplicationStatus `json:"status,omitempty,omitzero"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +optional
	Items []Application `json:"items"`
}

type ApplicationSpec struct {
	// +optional
	Description string `json:"description,omitempty"`
}

type ApplicationStatus struct {
}
