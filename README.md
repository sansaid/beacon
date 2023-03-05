# beacon

> ⚠️ Beacon is not intended for use with production services or services that handle sensitive information! Its purpose is to serve running applications on a budget without having to pay a cloud PaaS. This is only useful for proof of concepts or toy projects.

Beacon is a service that allows you to use your old machines as a platform. Think of it as a much simpler Docker Swarm and Kubernetes, but with the ability to expose your service to the internet.

With Beacon, you only need to push your images to a chosen container registry in order for your services to automatically be updated. You can run beacon in two modes:

## Fleet mode

In **fleet mode**, beacon operates as one of several beacons reporting to the _mothership_. This allows your mothership to manage which services are running where depending on the availability of resources reported by the beacons. You will need to connect your beacon to your account at `summonbeacon.com` in order to operate in this mode.

When running in this mode, you will manage your services entirely from the mothership, with a one time setup needed to register your device to the mothership.

### Setup

To operate in this mode, you will need to run one beacon per device. To do so, simply run the following command:

```sh
curl -s -L https://summonbeacon.com/install.sh | bash -s --mode fleet
```

## Solo mode

In **solo mode**, the beacon operates independently and does not delegate its management to the mothership.

To add services to be managed by beacon, you need to interact with it directly through the CLI on the device it's running on. If you have other beacons running, the other beacons will not be aware that you have changed anything on the current beacon.

### Setup

To operate in this mode, you will need to run one beacon per device. To do so, simply run the following command:

```sh
curl -s -L https://summonbeacon.com/install.sh | bash -s --mode solo
```

# Prerequisites

To run beacon, you need to have either `podman` or `docker` installed.

# Concepts

## Beacon

`beacon` is a perpetual service running on your device that manages which services are running on your device, which enpodints are exposed and which image is mapped to the service in the container registry. You can only have one `beacon` per device.

## Mothership

The `mothership` is the manager for all your beacons if you're running in fleet mode. It's only accessible through `summonbeacon.com`. You can do all the things you can do in solo mode with the mothership, except it's applied across all of your beacons. If a beacon is down, the mothership needs to decide how to rebuild the service 

## Service

A `service` is simply a service that you're running which is managed by a beacon. To run a service, it is expected that you either have `podman` or `docker` installed.

## Image

An `image` is just an OCI container image. With beacon, a service is directly mapped to an image repo in a container registry. Any new images regsitered in the repo will prompt a redeployment of the container using the new image.
