package main

import (
	"fmt"
	"io"
	"os"

	"github.com/crossplane/crossplane-runtime/pkg/errors"

	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composite"
	"github.com/crossplane/crossplane/apis/apiextensions/fn/io/v1alpha1"
	"sigs.k8s.io/yaml"
)

const (
	NamespaceForFlux = "flux-system"
)

func main() {
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, errors.Wrap(err, "failed to read stdin"))
		os.Exit(1)
	}
	result, err := Run(in)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	_, _ = fmt.Fprint(os.Stdout, result)
}

func Run(in []byte) (string, error) {
	obj := &v1alpha1.FunctionIO{}
	if err := yaml.Unmarshal(in, obj); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal stdin")
	}
	xkubernetesCluster := composite.New()
	if err := yaml.Unmarshal(obj.Observed.Composite.Resource.Raw, &xkubernetesCluster.Unstructured); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal observed composite")
	}
	teams, err := GetTeamEntries(fieldpath.Pave(xkubernetesCluster.Object))
	if err != nil {
		return "", errors.Wrap(err, "failed to get teams from observed composite")
	}
	// NOTE(muvaf): This assumes that the ProviderConfig of the cluster for
	// provider-kubernetes has the same name as the XEKS resource.
	providerConfigName, err := GetXEKSName(xkubernetesCluster)
	if err != nil {
		return "", errors.Wrap(err, "failed to get XEKS name from observed composite")
	}
	resources, err := GetResources(NamespaceForFlux, providerConfigName, teams)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate resources")
	}

	obj.Desired.Resources = resources
	result, err := yaml.Marshal(obj)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal resulting functionio")
	}
	return string(result), nil
}
