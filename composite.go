package main

import (
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composed"
	"k8s.io/apimachinery/pkg/util/json"

	"github.com/crossplane/crossplane/apis/apiextensions/fn/io/v1alpha1"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
)

func GetTeamEntries(cp *fieldpath.Paved) ([]TeamEntry, error) {
	teams := []TeamEntry{}
	if err := cp.GetValueInto("spec.teams", &teams); err != nil {
		return nil, errors.Wrap(err, "failed to get spec.teams from observed composite")
	}
	return teams, nil
}

func GetNameForProviderConfigs(observed []v1alpha1.ObservedResource) (string, error) {
	for _, r := range observed {
		comp := composed.New()
		if err := json.Unmarshal(r.Resource.Raw, &comp.Unstructured); err != nil {
			return "", errors.Wrap(err, "failed to unmarshal observed resource")
		}
		gvk := comp.GroupVersionKind()
		if gvk.Kind != "XEKS" || gvk.Group != "demo.upbound.io" {
			continue
		}
		name, err := fieldpath.Pave(comp.Object).GetString("spec.nameForProviderConfigs")
		if err != nil {
			return "", errors.Wrap(err, "failed to get spec.nameForProviderConfigs from XEKS")
		}
		return name, nil
	}
	return "", errors.New("failed to find XEKS reference")
}
