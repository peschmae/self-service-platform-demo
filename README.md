# Project self-service-platform

Very small application that allows to create a namespace with some default netpols, and customizable netpols.

I've written this for a demonstration of a very small self-service kiosk in context of platform engineering.

If you run the app locally, it will use the currently active kubectl context to connect to a cluster, but
if the app is run within a cluster, it will automatically use the serviceAccount of the pod.

Since the app creates namespaces and netpols it requires quite a lot of permissions, and the example
deployment config in the `k8s` directory, grants cluster admin!

## Getting Started

The project was created using the [go-blueprint](https://github.com/Melkeydev/go-blueprint)

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
