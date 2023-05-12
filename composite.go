package main

import (
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composed"
)

func GetTeamEntries(cp *fieldpath.Paved) ([]TeamEntry, error) {
	teams := []TeamEntry{}
	if err := cp.GetValueInto("spec.teams", &teams); err != nil {
		return nil, errors.Wrap(err, "failed to get spec.teams from observed composite")
	}
	return teams, nil
}

func GetNameForProviderConfigs(observed map[string]*composed.Unstructured) (string, error) {
	for _, r := range observed {
		gvk := r.GroupVersionKind()
		if gvk.Kind != "XEKS" || gvk.Group != "demo.upbound.io" {
			continue
		}
		name, err := fieldpath.Pave(r.Object).GetString("spec.nameForProviderConfigs")
		if err != nil {
			return "", errors.Wrap(err, "failed to get spec.nameForProviderConfigs from XEKS")
		}
		return name, nil
	}
	return "", nil
}
