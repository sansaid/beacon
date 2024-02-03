# beacon

> ⚠️ Beacon is not intended for use with production services or services that handle sensitive information! Its purpose is to serve running applications on a budget without having to pay a cloud PaaS an arm and a leg. This is only useful for proof of concepts or toy projects.

Beacon is a service that allows you to use your old machines as a platform. Think of it as a much simpler Docker Swarm or Kubernetes, but with the ability to expose your service to the internet.

With Beacon, you only need to push your images to a chosen container registry in order for your services to automatically be updated. You can run beacon in two modes:

## Solo mode

In **solo mode**, the beacon operates independently and does not delegate its management to the mothership. In other words, updates to the services that you're running must happen manually and by you.

To add services to be managed by beacon, you need to interact with it directly through the CLI on the device it's running on. If you have other beacons running, the other beacons will not be aware that you have changed anything on the current beacon.

### Setup

To operate in this mode, you will need to run one beacon per device. To do so, simply run the following command:

```sh
# Under construction - lytbeacon.com is not operational yet
curl -s -L https://lytbeacon.com/install.sh | bash -s --mode solo
```


## Fleet mode

> 🚧 **UNDER CONSTRUCTION** - this mode is not available yet. Please check again later!

In **fleet mode**, beacon operates as one of several beacons reporting to the _mothership_. This allows your mothership to manage which services are running where depending on the availability of resources reported by the beacons. You will need to connect your beacon to your account at `lytbeacon.com` in order to operate in this mode.

When running in this mode, you will manage your services entirely from the mothership, with a one time setup needed to register your device to the mothership.

### Setup

To operate in this mode, you will need to run one beacon per device. To do so, simply run the following command:

```sh
# Under construction - lytbeacon.com is not operational yet
curl -s -L https://lytbeacon.com/install.sh | bash -s --mode fleet
```

# Prerequisites

To run beacon, you need to have either `podman` or `docker` installed.

# Operating Principles

To use Beacon, there are a few operating principles that should be assumed

* Beacon is only useful with toy projects or proof of concepts that do not rely on high availability, scale or security
* For now, Beacon only works really well with worker type services (for example, Telegram bots, Discord bots, a report maker, etc.) - that is until we can make endpoint management available
* Beacon is very, very simple - its sole job is to make sure a container is running with the latest version currently available at an image repo
* Treat the image manifest of your Beacon service as the only input that defines your service; in other words:
  * Treat the `EXEC` field in your image manifest as the value you ultimately want to run in your environment (refer to the **Beacon is very, very simple** principle)
  * Use the `ENV` declarative in your image manifest to define non-sensitive environment variables you want to run your environment
  * For sensitive values, you should retrieve these values at runtime by calling a local or remote secret store (in future iterations of Beacon, we plan to provide this secret store by default)

# Concepts

## Beacon

`beacon` is a perpetual service running on your device that manages which services are running, which endpoints are exposed and which image is mapped to the service in the container registry. You can only have one `beacon` per device.

## Mothership

> 🚧 **UNDER CONSTRUCTION** - the mothership is still being built. Please check again later!

The `mothership` is the manager for all your beacons if you're running in fleet mode. It's only accessible through `lytbeacon.com`. You can do all the things you can do in solo mode with the mothership, except it's applied across all of your beacons. If a beacon is down, the mothership needs to decide how to rebuild the service.

## Service

A `service` is simply a service that you're running which is managed by a beacon. To run a service, it is expected that you either have `podman` or `docker` installed.

## Image

An `image` is just an OCI container image. With beacon, a service is directly mapped to an image repo in a container registry. Any new images regsitered in the repo will prompt a redeployment of the container using the new image.
