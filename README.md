# Portainer.cc - Building and Deploying Runtime Encrypted Workloads leveraging Confidential Compute

![](https://github.com/enclaive/portainerCC/blob/develop/wip-screens.gif)

## Table of Contents

- [About The Project](#about-the-project)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Install PortainerCC](#install-portainercc)
  - [Remote Attestation and Secret Provisioning](#remote-attestation-and-secret-provisioning)
- [How-tos](#how-tos)
- [Licence](#licence)

## About The Project

In view of the ever increasing shift of applications to the cloud, new mechanisms need to be developed to protect the workload. In contrast to on-prem, physical resources are no more isolated in the cloud. Rather virtual machines, kubernetes clusters and serverless functions, share physical resources. Moreover, the resources are maintained by a third party known as the cloud provider who has root access to the resources. For decades it is well known that the application isolation provided by hypervisors and operating systems is weak. A vast amount of exploits have been demonstrated how to escapte the present security and trust model.

Confidential Computing, for short CC, is a new, promising technology addressing the problem. CC makes it for the very first time practically possible to encrypt data during runtime in such a way that only the CPU has access to it. This makes it possible to protect application code and data in the light of vertical and horizontal exploits.

Portainer.cc is a project extending the promiment community tool [Portainer.io](https://github.com/portainer/portainer) with confidential computing capabilities. to make it easy to run application-containers confidentially in the cloud. PortainerCC builds upon [Gramine OS](https://github.com/gramineproject/gramine) and [Marblerun](https://github.com/edgelesssys/marblerun) to run and remotely attest containerized Gramine-applications.

## Features (v.0.1.0-beta)

Portainer.cc offers these features:

- Build and deploy any application in an Intel SGX enclave supporting Gramine libOS [Gramine](https://github.com/gramineproject/gramine)
- Key managmement for container authentication and file/volume encryption
- Authenticated container provisioning of secrets, environment variables, files and keys supporting [Marblerun](https://github.com/edgelesssys/marblerun)
- Example template to build, deploy and securely provision MariaDB


## Getting Started

### Prerequisites

For Portainer.cc to work, you need to make sure that all environments you want to use are Intel SGX compatible and can use Intel SGX Datacenter Attestation Primitives for Remote Attestation and meet these requirements:

- [Intel SGX and DCAP](https://download.01.org/intel-sgx/latest/dcap-latest/linux/docs/Intel_SGX_SW_Installation_Guide_for_Linux.pdf) are installed

- A [Provisioning Certificate Caching Service](https://docs.edgeless.systems/ego/reference/attest#set-up-the-pccs) is up and running

### Install Portainer.cc

To install Portainer.cc, run the following command:

```
docker run -d -p 8000:8000 -p 9000:9000 -p 9443:9443 \
-v /var/run/docker.sock:/var/run/docker.sock:z \
-v /var/run/docker.sock:/var/run/alternative.sock:z \
-v /tmp:/tmp \
-v pccdata:/data \
--name portainerCC \
marcely0/pcc
```

The Portainer.cc Image comes with some predefined confidential templates.
You can mount your own templates with the following parameter when you start your Container:
```
-v ./temps.json:/confidential-templates.json 
```

### How-Tos
You can check out some of the How-Tos:

[Step by Step guide to set up PortainerCC with an PortainerCC Agent](https://github.com/enclaive/portainerCC/wiki/PortainerCC-Full-Guide)

[Create a confidential Application for Portainer.cc](https://github.com/enclaive/portainerCC/wiki/PortainerCC-Create-Template)


### Remote Attestation and Secret Provisioning

[Deprecated - Step by Step guide to run MariaDB in PortainerCC](https://github.com/enclaive/portainerCC/wiki/PortainerCC-MariaDB-Guide)

## Licence

Distributed under the zlib licence. See [LICENCE](./License) for reference.
