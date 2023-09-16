/*
Copyright 2023 ahwhya.

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

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Operator 的核心逻辑就是不断调谐资源对象的实际状态和期望状态(Spec)保持一致
// 大多数资源对象都有Spec和Status两个部分，但是也有部分资源对象不符合这种模式，比如 ConfigMap 之类的静态资源对象就不存在着 "期望的状态" 这一说法

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Deployment DeploymentTemplate `json:"deployment,omitempty"`
	Service    ServiceTemplate    `json:"service,omitempty"`
}

type DeploymentTemplate struct {
	appsv1.DeploymentSpec `json:",inline"`
}

type ServiceTemplate struct {
	corev1.ServiceSpec `json:",inline"`
}

// ApplicationStatus defines the observed state of Application
type ApplicationStatus struct {
	// 这里的 Status 也不是严格对应"实际状态"，而是观察并记录下来的当前对象最新"状态"
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Workflow appsv1.DeploymentStatus `json:"workflow"`
	Network  corev1.ServiceStatus    `json:"network"`
}

// 这个标记主要是被 controller-tools 识别，然后 controller-tools 的对象生成器就知道这个标记下面的对象代表一个 Kind，接着对象生成器会生成相应的 Kind 需要的代码，也就是实现 runtime.Object 接口
// 换言之，一个结构体要表示一个Kind，必须实现runtime.Object接口

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

//+kubebuilder:resource:path=applications,singular=application,scope=Namespaced,shortName=app

// Application is the Schema for the applications API
type Application struct {
	// Application 结构体是 Application 类型的"根类型"，和其他所有的 Kubernetes 资源类型一样包含 TypeMeta 和 ObjectMeta
	// TypeMeta 中存放的是当前资源的 Kind 和 APIVersion 信息
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta 中存放的是 Name、Namespace、Labels 和 Annotations 等信息
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// 通过 Items 存放一组 Application，用于 List 之类的批量操作
	Items []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
