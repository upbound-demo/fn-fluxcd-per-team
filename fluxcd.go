package main

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"

	sourcev1 "github.com/fluxcd/source-controller/api/v1"
)

const (
	RepositoryIntervalDefault = 1 * time.Minute
	RepositoryBranchDefault   = "main"

	KustomizationIntervalDefault = 1 * time.Minute
)

type GitRepositoryOverride func(*sourcev1.GitRepository)

func GitRepository(name, namespace, repoUrl string, overrides ...GitRepositoryOverride) *sourcev1.GitRepository {
	gr := &sourcev1.GitRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sourcev1.GitRepositorySpec{
			URL: repoUrl,
			Interval: v1.Duration{
				Duration: RepositoryIntervalDefault,
			},
			Reference: &sourcev1.GitRepositoryRef{
				Branch: RepositoryBranchDefault,
			},
		},
	}
	gr.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   sourcev1.GroupVersion.Group,
		Version: "v1beta2",
		Kind:    sourcev1.GitRepositoryKind,
	})
	for _, override := range overrides {
		override(gr)
	}
	return gr
}

func WithPath(p string) KustomizationOverride {
	return func(k *kustomizev1.Kustomization) {
		k.Spec.Path = p
	}
}

type KustomizationOverride func(*kustomizev1.Kustomization)

func Kustomization(name string, gr *sourcev1.GitRepository, overrides ...KustomizationOverride) *kustomizev1.Kustomization {
	k := &kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: gr.GetNamespace(),
		},
		Spec: kustomizev1.KustomizationSpec{
			Interval: v1.Duration{
				Duration: KustomizationIntervalDefault,
			},
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind:      "GitRepository",
				Name:      gr.GetName(),
				Namespace: gr.GetNamespace(),
			},
		},
	}
	k.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   kustomizev1.GroupVersion.Group,
		Version: "v1beta2",
		Kind:    kustomizev1.KustomizationKind,
	})
	for _, override := range overrides {
		override(k)
	}
	return k
}
