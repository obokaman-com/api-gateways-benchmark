# Simple API Gateway benchmarking

Simple performance test using `hey` to launch 100K requests with 50, 100, 250, 500 & 1000 concurrent calls to a fake API backend endpoint (using `jaxgeller/lwan`)

All the API Gateways and even `hey` run using Docker, so the results may vary depending on the host machine and OS (remember that MacOS uses a virtual machine to run Docker).

Results are being saved in the `results` folder.

### Execute performance test on all Api Gateways

```shell
$ make all
```

### Execute performance tests individually

Write `make` + the name of the Gateway.

```shell
$ make [krakend|kong|tyk|nginx]
```
