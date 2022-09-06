# My Openfaas Functions

This repository contains my OpenFaaS functions and instructions on how to get
started yourself.

It details how to install your own OpenFaaS platform, how to deploy existing
OpenFaaS functions to it and how to create your own OpenFaas functions (and
deploy those).

The [OpenFaaS](https://www.openfaas.com/) platform lets you run Functions
as a Service, akin to AWS Lambdas or other Cloud providers' Serverless compute
offerings. You have to run the platform on your own infrastructure though, so
not _serverless_.

## Getting Started

### `faasd`

You first need to deploy the FaaS platform somewhere. You have several options,
either Kubernetes based or on bare-metal/VM (more [here](https://github.com/openfaas/faas-cli)).

We are going to use [faasd](https://github.com/openfaas/faasd):

> faasd is OpenFaaS reimagined, but without the cost and complexity of
Kubernetes. It runs on a single host with very modest requirements, making it
fast and easy to manage. Under the hood it uses containerd and Container
Networking Interface (CNI) along with the same core OpenFaaS components from
the main project.

### `multipass`

For managing the VM, lets use [multipass](https://multipass.run/):

> Get an instant Ubuntu VM with a single command. Multipass can launch and run
virtual machines and configure them with cloud-init like a public cloud.

Just follow [these steps](https://github.com/openfaas/faasd/blob/master/docs/MULTIPASS.md#lets-start-the-tutorial).
It will spin up a VM and install everything you need.

Take not of the IP address
You can now access the WebUI, which you can use to deploy and invoke existing
OpenFaas functions. You can find pre-packaged functions in the Function Store,
or build your own using any OpenFaaS template.

### `faas-cli`

You will also need [`faas-cli`](https://github.com/openfaas/faas-cli) on your
host machine:

> faas-cli is the official CLI for OpenFaaS

`faas-cli` will be available in your VM, but it's easier to also have it on
your host machine. This allows you to build and deploy FaaS functions from the
host machine, rather than having to ssh into the VM:

`brew install faas-cli`

You can also use `faas-cli` to deploy and invoke functions instead of using the
WebUI.

[Try to deploy a pre-packaged function.](https://github.com/openfaas/faasd/blob/master/docs/MULTIPASS.md#try-faasd-openfaas).

## Deploy a function

This repository comes with OpenFaaS functions that I wrote. Let's start by
building and deploying one of those so you get a feel for how the whole process
works.



## Create your first function

We will start with a simple "Hello World" example, that is a function which,
when called, responds with "Hello World" in the body. We will be using Go for this
example but you should be able to translate this in any [language supported by
OpenFaas](https://github.com/openfaas/templates/#classic-templates).

Use the `golang-middleware` template to bootstrap your project:

```shell
faas-cli new --lang golang-middleware hello-world
```

This will create a directory called `hello-world` and a `hello-world.yml` file.
The YAML file is your



### Bare Metal (Rpi)

```shell
git clone https://github.com/openfaas/faasd --depth=1
cd faasd

./hack/install.sh
```

## How it works

[Start here](https://www.openfaas.com/blog/template-store/#:~:text=%2D%2Dparallel%20flag.-,How%20it%20works,-If%20you%20run)

The development process of OpenFaaS functions goes something like this:

* pick the language you want to use for writing your function ([list of supported
languages](https://github.com/openfaas/templates/#classic-templates))
* Use the template for that language to bootstrap your OpenFaaS project/function
* Write your handler function
* Build and Publish your function
* Deploy your function to your OpenFaaS environment


