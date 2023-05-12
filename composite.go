package main

import (
	"strings"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composite"
)

func GetTeamEntries(cp *fieldpath.Paved) ([]TeamEntry, error) {
	teams := []TeamEntry{}
	if err := cp.GetValueInto("spec.teams", &teams); err != nil {
		return nil, errors.Wrap(err, "failed to get spec.teams from observed composite")
	}
	return teams, nil
}

func GetXEKSName(cp *composite.Unstructured) (string, error) {
	for _, ref := range cp.GetResourceReferences() {
		if ref.Kind == "XEKS" && strings.HasPrefix(ref.APIVersion, "demo.upbound.io/") {
			return ref.Name, nil
		}
	}
	return "", errors.New("failed to find XEKS reference")
}
