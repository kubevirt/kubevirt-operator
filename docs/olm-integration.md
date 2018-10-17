# OLM integration

The OLM integration guide will outline the basic concepts needed in order to
integrate with OLM.

### Glossary
| Resource                 | Short Name | Owner   | Description                                                                                |
|--------------------------|------------|---------|--------------------------------------------------------------------------------------------|
| ClusterServiceVersion | CSV        | OLM     | application metadata: name, version, icon, required resources, installation, etc...        |
| InstallPlan           | IP         | Catalog | calculated list of resources to be created in order to automatically install/upgrade a CSV |
| CatalogSource         | CS         | Catalog | a repository of CSVs, CRDs, and packages that define an application                        |
| Subscription          | Sub        | Catalog | used to keep CSVs up to date by tracking a channel in a package                            |

### Getting Started
The operator-sdk has a useful command to generate a CSV and CRDs.  This command
is subject to change in the future, but for now it will get you off the ground.

```bash
$ operator-sdk generate olm-catalog ...
```

#### CatalogSource
```yaml
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: kubevirt-operator
  namespace: operator-lifecycle-manager
spec:
  sourceType: internal
  configMap: kubevirt-operator
  displayName: KubeVirt Operator
  publisher: Red Hat
```

#### ConfigMap
*This config map is referenced by the CatalogSource (see above).*

##### Outline
```yaml
Header
data:
  customResourceDefinition: |-
    - list
  clusterServiceVersion: |-
    - list
  packages: |-
    - list
```

##### Header
```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: kubevirt-operator
  namespace: operator-lifecycle-manager
data:
```

##### Custom Resource Definition (CRD)
CRD yaml [fields](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/Documentation/design/building-your-csv.md#crd-templates).

```yaml
**Header**
...
  customResourceDefinitions: |-
    - apiVersion: apiextensions.k8s.io/v1beta1
      kind: CustomResourceDefinition
      metadata:
        name: virtualization.virt.kubevirt.io
      spec:
        version: v1alpha1
        group: virt.kubevirt.io
        names:
          kind: Virt
          listKind: virtList
          plural: virtualization
          singular: virt
        scope: Namespaced

...
**ClusterServiceVersion (CSV)**
**Pacakges**
```

##### Custom Service Version (CSV)
The CSV section is split into three parts:
  - [ClusterServiceVersion Metadata](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/Documentation/design/building-your-csv.md#operator-install)
  - [Operator Metadata](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/Documentation/design/building-your-csv.md#operator-metadata)
  - [Operator Install](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/Documentation/design/building-your-csv.md#operator-install)

```yaml
**Header**
**ClusterResourceDefinition**
...
  clusterServiceVersions: |-


      # ----- CSV Metadata ----- #


    - apiVersion: operators.coreos.com/v1alpha1
      kind: ClusterServiceVersion
      metadata:
        name: kubevirtoperator.v0.8.0
        annotations:
          tectonic-visibility: ocs
      spec:


      # ----- Operator Metadata ----- #


        maturity: stable
        version: 0.8.0
        displayName: KubeVirt
        description: |-
           KubeVirt enables the migration of existing virtualized workloads directly into the development workflows supported by Kubernetes.
           This provides a path to more rapid application modernization by:
              - Supporting development of new microservice applications in containers that interact with existing virtualized applications.
              - Combining existing virtualized workloads with new container workloads on the same platform, thereby making it easier to decompose monolithic virtualized workloads into containers over time.

        keywords: ['kubevirt', 'virtualization', 'virt']
        maintainers:
        - name: RedHat
          email: rhallise@redhat.com
        provider:
          name: RedHat
        labels:
          operated-by: kubevirt
        selector:
          matchLabels:
            operated-by: kubevirt
        links:
        - name: Learn more about the project
          url: http://kubevirt.io/

        icon: # convert image to base64
        - base64data: ...
          mediatype: image/png


       # ----- Operator Install ----- #


        install:
          strategy: deployment
          spec:
            clusterPermissions: # Cluster level permissions (clusterRole and clusterRoleBinding)
            permissions: # User level permissions
            - serviceAccountName: kubevirt-operator

              ** COPY `deploy/rbac.yaml HERE` **

            deployments:

              ** COPY `deploy/operator.yaml` HERE **
...
**Packages**
```

##### Packages
```yaml
**Header**
**ClusterResourceDefinition**
**ClusterServiceVersion**
...
  packages: |-
    - packageName: kubevirt
      channels:
      - name: alpha
        currentCSV: kubevirtoperator.v0.9.0-alpha.0
```
