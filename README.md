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

KrakenD
```shell
$ curl -i http://localhost:8080/test
$ curl -i -H "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpbTIifQ.eyJpc3MiOiJhMzZjMzA0OWIzNjI0OWEzYzlmODg5MWNiMTI3MjQzYyIsImV4cCI6OTk0MjQzMDA1NCwibmJmIjoxNDQyNDI2NDU0LCJpYXQiOjE0NDI0MjY0NTQsInBvbCI6ImRlZmF1bHQiLCJzdWIiOiIxIiwianRpIjoiNWI0MzczZWU2M2M3ZDM1MGVmOGZkOWFmZmQ1MzgzNGUifQ.ew4Qe36jUYwKm1elij8opPabHtx75RDajBeHVRr-WUg" http://localhost:8080/auth
$ curl -i -H "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpbTIifQ.eyJpc3MiOiJhMzZjMzA0OWIzNjI0OWEzYzlmODg5MWNiMTI3MjQzYyIsImV4cCI6OTk0MjQzMDA1NCwibmJmIjoxNDQyNDI2NDU0LCJpYXQiOjE0NDI0MjY0NTQsInBvbCI6ImRlZmF1bHQiLCJzdWIiOiIxIiwianRpIjoiNWI0MzczZWU2M2M3ZDM1MGVmOGZkOWFmZmQ1MzgzNGUifQ.ew4Qe36jUYwKm1elij8opPabHtx75RDajBeHVRr-WUg" http://localhost:8080/backend
```

Others
```shell
$ curl -i http://localhost:8080/test
$ curl -i -H "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpbTIifQ.eyJpc3MiOiJhMzZjMzA0OWIzNjI0OWEzYzlmODg5MWNiMTI3MjQzYyIsImV4cCI6OTk0MjQzMDA1NCwibmJmIjoxNDQyNDI2NDU0LCJpYXQiOjE0NDI0MjY0NTQsInBvbCI6ImRlZmF1bHQiLCJzdWIiOiIxIiwianRpIjoiNWI0MzczZWU2M2M3ZDM1MGVmOGZkOWFmZmQ1MzgzNGUifQ.ew4Qe36jUYwKm1elij8opPabHtx75RDajBeHVRr-WUg" http://localhost:8080/auth
$ curl -i -H "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpbTIifQ.eyJpc3MiOiJhMzZjMzA0OWIzNjI0OWEzYzlmODg5MWNiMTI3MjQzYyIsImV4cCI6OTk0MjQzMDA1NCwibmJmIjoxNDQyNDI2NDU0LCJpYXQiOjE0NDI0MjY0NTQsInBvbCI6ImRlZmF1bHQiLCJzdWIiOiIxIiwianRpIjoiNWI0MzczZWU2M2M3ZDM1MGVmOGZkOWFmZmQ1MzgzNGUifQ.ew4Qe36jUYwKm1elij8opPabHtx75RDajBeHVRr-WUg" http://localhost:8080/backend1
$ curl -i -H "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpbTIifQ.eyJpc3MiOiJhMzZjMzA0OWIzNjI0OWEzYzlmODg5MWNiMTI3MjQzYyIsImV4cCI6OTk0MjQzMDA1NCwibmJmIjoxNDQyNDI2NDU0LCJpYXQiOjE0NDI0MjY0NTQsInBvbCI6ImRlZmF1bHQiLCJzdWIiOiIxIiwianRpIjoiNWI0MzczZWU2M2M3ZDM1MGVmOGZkOWFmZmQ1MzgzNGUifQ.ew4Qe36jUYwKm1elij8opPabHtx75RDajBeHVRr-WUg" http://localhost:8080/backend2
$ curl -i -H "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpbTIifQ.eyJpc3MiOiJhMzZjMzA0OWIzNjI0OWEzYzlmODg5MWNiMTI3MjQzYyIsImV4cCI6OTk0MjQzMDA1NCwibmJmIjoxNDQyNDI2NDU0LCJpYXQiOjE0NDI0MjY0NTQsInBvbCI6ImRlZmF1bHQiLCJzdWIiOiIxIiwianRpIjoiNWI0MzczZWU2M2M3ZDM1MGVmOGZkOWFmZmQ1MzgzNGUifQ.ew4Qe36jUYwKm1elij8opPabHtx75RDajBeHVRr-WUg" http://localhost:8080/backend3
```