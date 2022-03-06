# Ubuntu Operator

<img src="images/ubuntunetes.png" width="300">

Control Ubuntu from Kubernetes.

![license](https://img.shields.io/github/license/cloud-native-skunkworks/ubuntu-operator)
![tags](https://img.shields.io/github/v/tag/cloud-native-skunkworks/ubuntu-operator)
![build](https://img.shields.io/github/workflow/status/cloud-native-skunkworks/ubuntu-operator/Docker%20Image%20CI)

```
apiVersion: ubuntu.machinery.io.canonical.com/v1alpha1
kind: UbuntuMachineConfiguration
metadata:
  name: ubuntumachineconfiguration-sample
spec:
  desiredModules:
  - name: "nvme_core"
    flags: ""
  - name: "rfcomm"
    flags: ""
  desiredPackages:
    apt:
    - name: "build-essentials"
    snap:
    - name: "microk8ss"
      confinement: "classic"
```


![modules](images/carbon.png)


Control your underlying Ubuntu distribution through Kubernetes....

![arch](images/arch.png)

## Roadmap

- [x] Kernel module support
- [ ] Package system support

## Installation

## Host-relay

`make install-relay`

## Operator 
```
make install # Uploads the CustomResourceDefinitions into your cluster
make deploy
```


## Development

### Operator 

After installing the CRD with `make install`
Run `go run main.go` to run the operator locally.


#### Notes

Regenerating the clientset was done from [this](https://www.fatalerrors.org/a/writing-crd-by-mixing-kubeuilder-and-code-generator.html) guide.
