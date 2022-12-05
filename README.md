# PortainerCC

## Table of Contents

- [About The Project](#about-the-project)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Install PortainerCC](#install-portainercc)
  - [Remote Attestation and Secret Provisioning](#remote-attestation-and-secret-provisioning)
- [Licence](#licence)

## About The Project

Especially in view of the ever increasing shift of applications to the cloud, the question is becoming more and more important whether the cloud environment used, over which the end user has only limited control, can be trusted. Confidential computing is one approach to solving this problem. Confidential computing makes it possible to encrypt data during processing in such a way that only the CPU has access to it. This makes it possible to protect data processed in the cloud against access by the cloud provider or other users of the cloud.

PortainerCC is based on [Portainer.io Community Edition](https://github.com/portainer/portainer) and extends Portainer with confidential computing capabilities to make it easy to run application-containers confidentially in the cloud. PortainerCC builds upon [Gramine OS](https://github.com/gramineproject/gramine) and [Marblerun](https://github.com/edgelesssys/marblerun) to run and remotely attest Gramine-applications.

## Features

In its current state, PortainerCC offers these features:

- Creating and storing Intel SGX Signing Keys
- Building and deploying a Remote Attestation System based on [Edgeless Systems Marblerun](https://github.com/edgelesssys/marblerun)
- Deploying a MariaDB instance running on [Gramine](https://github.com/gramineproject/gramine) that gets remote attested and receives login credentials via Secret Provisioning

## Getting Started

### Prerequisites

For PortainerCC to work, you need to make sure that all environments you want to use are Intel SGX compatible and can use Intel SGX Datacenter Attestation Primitives for Remote Attestation and meet these requirements:

- [Intel SGX and DCAP](https://download.01.org/intel-sgx/latest/dcap-latest/linux/docs/Intel_SGX_SW_Installation_Guide_for_Linux.pdf) are installed

- A [Provisioning Certificate Caching Service](https://docs.edgeless.systems/ego/reference/attest#set-up-the-pccs) is up and running

### Install PortainerCC

To install PortainerCC, run the following command:

```
docker run -d -p 8000:8000 -p 9443:9443 --name portainercc --restart=always -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data sgxdcaprastuff/portainercc
```

### Remote Attestation and Secret Provisioning

WIKILINK

## Licence

Distributed under the zlib licence. See [LICENCE](./License) for reference.
