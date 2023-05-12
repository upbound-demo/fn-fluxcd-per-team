package main

import (
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane/apis/apiextensions/fn/io/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
)

const (
	EnvironmentDev     = "development"
	EnvironmentStaging = "staging"
	EnvironmentProd    = "production"
)

var (
	EnvironmentsDefault = []Environment{
		{
			Name: EnvironmentDev,
			Path: EnvironmentDev,
		},
		{
			Name: EnvironmentStaging,
			Path: EnvironmentStaging,
		},
		{
			Name: EnvironmentProd,
			Path: EnvironmentProd,
		},
	}
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

func GetResources(namespace, providerConfigName string, teams []TeamEntry) ([]v1alpha1.DesiredResource, error) {
	resources := []v1alpha1.DesiredResource{}
	for _, team := range teams {
		gr := GitRepository(team.Name, namespace, team.Repository)
		grRaw, err := json.Marshal(gr)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal GitRepository %s", gr.GetName())
		}
		grORaw, err := WrapForKubernetes(runtime.RawExtension{Raw: grRaw}, providerConfigName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to wrap GitRepository %s into Object", gr.GetName())
		}
		resources = append(resources, v1alpha1.DesiredResource{
			Name:     "gitrepository-" + team.Name,
			Resource: grORaw,
		})
		for _, env := range team.Environments {
			k := Kustomization(team.Name+"-"+env.Name, gr, WithPath(env.Path))
			kRaw, err := json.Marshal(k)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to marshal Kustomization %s", k.GetName())
			}
			kORaw, err := WrapForKubernetes(runtime.RawExtension{Raw: kRaw}, providerConfigName)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to wrap GitRepository %s into Object", gr.GetName())
			}
			resources = append(resources, v1alpha1.DesiredResource{
				Name:     "kustomization-" + team.Name + "-" + env.Name,
				Resource: kORaw,
			})
		}
	}
	return resources, nil
}
