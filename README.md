# concourse-k8s-resource

[![Go Report Card](https://goreportcard.com/badge/github.com/kudoh/concourse-k8s-resource)](https://goreportcard.com/report/github.com/kudoh/concourse-k8s-resource)

Concourse CI custom resource for kubernetes.

This resource assumes that kubernetes deployment is executed by plain manifests(`kubectl apply -f`) or Kustomize(`kubectl apply -k`).
In addition to deploy(`kubectl apply`) operation, it also supports delete(`kubectl delete`) and undo(`kubectl rollout undo`) operations.

## Source Configuration

* **`api_server_url`** - Kubernetes API Server URL. (e.g. `https://172.16.10.11:6443`)
* `api_server_cert` - Kubernetes api server certificate.
* `client_cert` - Client certificate that authenticated to access cluster, if specified, `clusterKey` is required. 
* `client_key` - Client key corresponding to `client_cert`.
* `client_token` - Kubernetes ServiceAccount token. Required when accessing a cluster with ServiceAccount(neither `client_cert` nor `client_key`).
* `kubeconfig` - Contents of kubecofig(YAML). if specified, other auth settings(including `api_server_url`) are ignored.
* `skip_tls_verify` - true if you want to skip TLS verification.
* `namespace` - Name of target kubernetes namespace. if not specified, `default` is used.

## Behavior

### `check`

Triggers when the watched resources is updated.
The Revision consists of namespace, resource name, and revision(Deployment or StatefulSet resource). 
Multiple resources are combined with `+`.

Example version file is `dev:app1:100+dev:app2:50+dev:app3:10`.

### `in`

Retrieves the watched resource's version and writes to `version` file.

### `out`

Deploys the watched resources to kubernetes using a plain manifest or Kustomize overlays. After deployed, wait for the specified time to complete.

* **`paths`** - kubernetes manifests path. plain manifests path(e.g. `k8s/deployment.yaml`) or kustomize directory(e.g. `repo/overlays/prod`)
* `kustomize` - true if deploying by kustomize. Default to `false`.
* `status_check_timeout` - The time(seconds) to wait for deployment to complete. Default to 5 minutes.
* `command_timeout` - The time(seconds) to wait for kubectl apply or delete. Default to unlimited(0).
* `delete` - true if using `kubectl delete` operation. Default to `false`.
* `undo` - true if using `kubectl rollout undo` operation(target resources are `watchedResources`). Default to `false`.

## Example

See full configuration example [here](./test/test-pipeline.yaml).

### `resource_types` 

```yaml
resource_types:
- name: k8s
  type: docker-image
  source:
    repository: kudohn/concourse-k8s-resource
    tag: <version>
```

### `resources`

```yaml
resources:
- name: k8s
  type: k8s
  source:

    api_server_url: https://172.16.10.11:6443
    api_server_cert: |
      -----BEGIN CERTIFICATE-----
      ....
      -----END CERTIFICATE-----
    # use client certificate
    client_cert: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    client_key: |
      -----BEGIN PRIVATE KEY-----
      ...
      -----END PRIVATE KEY-----

    # or use service account token
    client_token: ....

    # or use kubeconfig
    kubeconfig: |
      apiVersion: v1
      kind: Config
      clusters:
        - cluster:
          ....
      contexts:
          ....
      current-context: ...
      users:
        - name: concourse
          user:
            ...

    skip_tls_verify: false
    namespace: dev
    # watched resources(deployment or statefulset is supported)
    watch_resources:
    - name: app1
      kind: Deployment
    - name: app2
      kind: Deployment
    - name: web
      kind: StatefulSet
```

### `put`

#### Deploy resources using plain k8s manifests

```yaml
jobs:
- name: deploy-app
  plan:
  - get: repo
  - put: k8s
    params:
      status_check_timeout: 60
      command_timeout: 30
      paths:
        - repo/test/plain/deploy1.yaml
        - repo/test/plain/deploy2.yaml
        - repo/test/plain/sts.yaml
```

#### Deploy resources using Kustomize manifests

```yaml
jobs:
- name: deploy-app
  plan:
  - get: repo
  - put: k8s
    params:
      kustomize: true
      status_check_timeout: 60
      command_timeout: 30
      paths:
      - repo/test/kustomize/overlays/prod
```

#### Undo Resources

```yaml
jobs:
- name: deploy-app
  plan:
  - get: repo
  - put: k8s
    params:
      undo: true
```

#### Delete Resources

```yaml
jobs:
- name: deploy-app
  plan:
  - get: repo
  - put: k8s
    params:
      delete: true
```

## License

[MIT License](./LICENSE)
