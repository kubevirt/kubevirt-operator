# kubevirt-operator
Operator that manages KubeVirt

### Build the Operator COntainer
```bash
operator-sdk build docker.io/rthallisey/kubevirt-operator
```

### Launch the Operator
```bash
# Deploy the app-operator
$ kubectl create -f deploy/rbac.yaml
$ kubectl create -f deploy/crd.yaml
$ kubectl create -f deploy/operator.yaml

# Create the operator resource (kubectl get apps)
$ kubectl create -f deploy/cr.yaml
```
