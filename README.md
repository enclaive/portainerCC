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

PortainerCC is based on [Portainer.io Community Edition](https://github.com/portainer/portainer) and extends Portainer with confidential computing capabilities. PortainerCC allows to deploy confidential gramine containers and attest them via Intel-SGX Remote Attestation.

## Features

In its current state, PortainerCC offers these features:

- Creating and storing Intel SGX Signing Keys
- Building and deploying a Remote Attestation System based on Edgeless Systems Marblerun
- Deploying a MariaDB instance running with GramineOS that gets remote attested and receives login credentials via Secret Provisioning

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
