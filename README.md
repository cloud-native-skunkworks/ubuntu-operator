# Ubuntu Operator


![modules](images/carbon.png)


```
apiVersion: ubuntu.machinery.io.canonical.com/v1alpha1
kind: UbuntuKernelModule
metadata:
  name: ubuntukernelmodule-sample
spec:
  desiredModules:
  - "nvme"
  - "rfcom"
  - "nvme_core"
```

Control your underlying Ubuntu distribution through Kubernetes....

![arch](images/arch.png)


## Installation

```
make install # Uploads the CustomResourceDefinitions into your cluster
make deploy
```
