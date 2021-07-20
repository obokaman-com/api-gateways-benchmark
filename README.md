# Simple API Gateway benchmarking

Simple performace test using `hey` to launch 100K requests with 1K concurrency to a fake API backend endpoint (using `jaxgeller/lwan`)

### Execute performance test on all Api Gateways

```shell
$ make all
```

### Execute performance tests individually

```shell
$ make [krakend|kong|tyk|nginx]
```
