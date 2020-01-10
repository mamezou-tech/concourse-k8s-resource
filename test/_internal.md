# Testing

## Concourse on local k8s

```sh
kubectl create ns concourse
helm repo add concourse https://concourse-charts.storage.googleapis.com/
helm upgrade --install concourse concourse/concourse --namespace concourse \
  --set persistence.worker.storageClass=openebs-cstor-sparse \
  --set postgresql.persistence.storageClass=openebs-cstor-sparse

export POD_NAME=$(kubectl get pods --namespace concourse -l "app=concourse-web" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward --namespace concourse $POD_NAME 8080:8080

fly login -t ci -c http://127.0.0.1:8080 -u test -p test -n main
fly -t ci set-pipeline -c test/test-pipeline.yaml -p test
fly -t ci up -p test
fly -t ci uj -j test/test
```

## Local Testing

```sh
# check
cat test/json/check_request_clientcert.json | go run cmd/check/main.go
cat test/json/check_request_kubeconfig.json | go run cmd/check/main.go
cat test/json/check_request_sa.json | go run cmd/check/main.go
# in
cat test/json/check_request_clientcert.json | go run cmd/in/main.go
cat test/json/check_request_kubeconfig.json | go run cmd/in/main.go
cat test/json/check_request_sa.json | go run cmd/in/main.go
# out
## deploy
cat test/json/out_request_apply_plain.json | go run cmd/out/main.go
cat test/json/out_request_apply_kustomize.json | go run cmd/out/main.go
## undo
cat test/json/out_request_undo_plain.json | go run cmd/out/main.go
cat test/json/out_request_undo_kustomize.json | go run cmd/out/main.go
## delete
cat test/json/out_request_delete_plain.json | go run cmd/out/main.go
cat test/json/out_request_delete_kustomize.json | go run cmd/out/main.go
```