package main

import (
	kubernetesv1alpha1 "github.com/crossplane-contrib/provider-kubernetes/apis/object/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
)

type TeamEntry struct {
	// Name of the team.
	Name string `json:"name"`

	// Repository is the URL of the git repository that contains the
	// manifests for this team.
	Repository string `json:"repository"`

	// Environments is a list of environments that should be deployed for this
	// team. If not specified, three default environments will be deployed.
	// development -> ./development
	// staging -> ./staging
	// production -> ./production
	Environments []Environment `json:"environments,omitempty"`
}

type Environment struct {
	// Name of the environment.
	Name string `json:"name"`

	// Path of the folder in the repository containing the manifests for this
	// environment.
	Path string `json:"path"`
}

type DesiredResource struct {
	Name     string
	Resource *kubernetesv1alpha1.Object
}

func GetResources(fluxNamespace, providerConfigName string, teams []TeamEntry) ([]DesiredResource, error) {
	resources := []DesiredResource{}
	for _, team := range teams {
		gr := GitRepository(team.Name, fluxNamespace, team.Repository)
		grRaw, err := json.Marshal(gr)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal GitRepository %s", gr.GetName())
		}
		resources = append(resources, DesiredResource{
			Name:     "gitrepository-" + team.Name,
			Resource: WrapForKubernetes(providerConfigName, runtime.RawExtension{Raw: grRaw}),
		})
		for _, env := range team.Environments {
			k := Kustomization(team.Name+"-"+env.Name, gr, WithPath(env.Path))
			kRaw, err := json.Marshal(k)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to marshal Kustomization %s", k.GetName())
			}
			resources = append(resources, DesiredResource{
				Name:     "kustomization-" + team.Name + "-" + env.Name,
				Resource: WrapForKubernetes(providerConfigName, runtime.RawExtension{Raw: kRaw}),
			})
		}
	}
	return resources, nil
}
