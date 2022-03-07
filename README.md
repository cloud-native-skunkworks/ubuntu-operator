# Ubuntu Operator

<img src="images/ubuntunetes.png" width="300">

Control Ubuntu from Kubernetes.

![license](https://img.shields.io/github/license/cloud-native-skunkworks/ubuntu-operator)
![tags](https://img.shields.io/github/v/tag/cloud-native-skunkworks/ubuntu-operator)
![build](https://img.shields.io/github/workflow/status/cloud-native-skunkworks/ubuntu-operator/Docker%20Image%20CI)

![cs](images/code-example.png)


Control your underlying Ubuntu distribution through Kubernetes....

![arch](images/arch.png)

## Roadmap

- [x] Kernel module support
- [ ] Package system support

## Installation

Two step installation process.
1. Installing the host-relay on all hosts
2. Installing the Operator in cluster once.

### Host-relay

`make install-relay`

### Operator 
```
make install # Uploads the CustomResourceDefinitions into your cluster
make deploy
```


