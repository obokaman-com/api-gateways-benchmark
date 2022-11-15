.PHONY: all krakend kong tyk nginx stop

all: krakend kong nginx tyk
# TOTAL_REQUESTS := 1000000
TOTAL_REQUESTS := 100000
TOKEN := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpbTIifQ.eyJpc3MiOiJhMzZjMzA0OWIzNjI0OWEzYzlmODg5MWNiMTI3MjQzYyIsImV4cCI6OTk0MjQzMDA1NCwibmJmIjoxNDQyNDI2NDU0LCJpYXQiOjE0NDI0MjY0NTQsInBvbCI6ImRlZmF1bHQiLCJzdWIiOiIxIiwianRpIjoiNWI0MzczZWU2M2M3ZDM1MGVmOGZkOWFmZmQ1MzgzNGUifQ.ew4Qe36jUYwKm1elij8opPabHtx75RDajBeHVRr-WUg"

test:
	cd ${PWD}/gateways/${name} ; \
# 	for c in 50 100 250 500 750 1000 1500 ; \
	for c in 50 100 250 500; \
	do \
		echo "Launching ${TOTAL_REQUESTS} requests to ${name} (${endpoint}) - Concurrency: $$c" ; \
		docker-compose up --build -d ; \
		sleep 3 ; \
		docker run --platform linux/amd64 --rm --net=host  williamyeh/hey -n ${TOTAL_REQUESTS} -c $$c${header} "http://localhost:8080/${endpoint}" > "${PWD}/results/test_results_${name}_${endpoint}_$${c}.txt" ; \
		docker-compose down --volumes --remove-orphans ; \
	done ; \
	cd ${PWD}

test-case:
	name=${name} endpoint=test header="" make -e test
	name=${name} endpoint=auth header=" -H \"Authorization: bearer ${TOKEN}\"" make -e test

krakend: stop
	name=krakend make -e test-case

kong: stop
	name=kong make -e test-case

tyk: stop
	name=tyk make -e test-case

nginx: stop
	name=nginx make -e test-case

stop:
	@cd ${PWD}/gateways/krakend ; docker-compose down --volumes --remove-orphans ; \
	cd ${PWD}/gateways/kong ; docker-compose down --volumes --remove-orphans ; \
	cd ${PWD}/gateways/tyk ; docker-compose down --volumes --remove-orphans ; \
	cd ${PWD}/gateways/nginx ; docker-compose down --volumes --remove-orphans ; \
	cd ${PWD}
