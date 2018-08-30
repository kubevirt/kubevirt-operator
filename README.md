# kubevirt-operator
Operator that manages KubeVirt

### Build the Operator Container
```bash
wget -o kubevirt.yaml https://github.com/kubevirt/kubevirt/releases/download/v0.6.4/kubevirt.yaml
operator-sdk build docker.io/rthallisey/kubevirt-operator
```

### Launch the Operator
```bash
# Give the kubevirt operator user 'default' cluster-admin privilages
$ oc adm policy add-cluster-role-to-user cluster-admin -z default

# Deploy the app-operator
$ kubectl create -f deploy/rbac.yaml
$ kubectl create -f deploy/crd.yaml
$ kubectl create -f deploy/operator.yaml

# Create the operator resource (kubectl get apps)
$ kubectl create -f deploy/cr.yaml
```
