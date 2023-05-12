package main

import (
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"

	kubernetesv1alpha1 "github.com/crossplane-contrib/provider-kubernetes/apis/object/v1alpha1"
)

func WrapForKubernetes(in runtime.RawExtension, providerConfigName string) (runtime.RawExtension, error) {
	o := kubernetesv1alpha1.Object{
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
	raw, err := json.Marshal(o)
	if err != nil {
		return runtime.RawExtension{}, errors.Wrap(err, "failed to marshal object")
	}
	return runtime.RawExtension{Raw: raw}, nil
}
