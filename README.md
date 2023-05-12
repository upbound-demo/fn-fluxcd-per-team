# FluxCD Resources For Every Team

This Crossplane composition function relies on the CompositeResourceDefinition
of `XKubernetesCluster` [here](https://github.com/upbound-demo/internal-cloud-platform/tree/main/platform/apis/production/kubernetescluster).

A [`GitRepository`](https://fluxcd.io/flux/components/source/gitrepositories/)
for every team and a
[`Kustomization`](https://fluxcd.io/flux/components/kustomize/kustomization/)
for every team's environment is wrapped in an
[`Object`](https://marketplace.upbound.io/providers/crossplane-contrib/provider-kubernetes/v0.8.0/resources/kubernetes.crossplane.io/Object/v1alpha1)
resource so that it's applied to a remote cluster.

It expects a `spec.teams` field with a list of team entries whose struct has to
comply with the `Team` struct defined [here](./teams.go) in this repository and
the composition needs to have an object of type `XEKS` in `demo.upbound.io` group
whose name will be used as the `providerConfigRef.name` of all generated `Object`
resources.

## Developing

Run the test.
```bash
cat test-input.yaml | go run . > /tmp/result.yaml
diff /tmp/result.yaml test-output.yaml
```

Build and push the image.
```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --tag ghcr.io/upbound-demo/xfn-fluxcd-per-team:v0.1.0 \
  --push .
```