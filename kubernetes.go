package main

import (
	kubernetesv1alpha1 "github.com/crossplane-contrib/provider-kubernetes/apis/object/v1alpha1"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func WrapForKubernetes(providerConfigName string, in runtime.RawExtension) *kubernetesv1alpha1.Object {
	o := &kubernetesv1alpha1.Object{
		Spec: kubernetesv1alpha1.ObjectSpec{
			ForProvider: kubernetesv1alpha1.ObjectParameters{
				Manifest: in,
			},
			ResourceSpec: v1.ResourceSpec{
				ProviderConfigReference: &v1.Reference{
					Name: providerConfigName,
				},
			},
		},
	}
	o.SetGroupVersionKind(kubernetesv1alpha1.ObjectGroupVersionKind)
	return o
}
