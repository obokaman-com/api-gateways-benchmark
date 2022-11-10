.PHONY: all krakend kong tyk nginx stop

all: krakend kong tyk nginx
TOTAL_REQUESTS := 100000

krakend:
	@cd ${PWD}/gateways/krakend ; \
	docker-compose up --build -d ; \
	for c in 50 100 250 500 1000 1500 2000 ; \
	do \
		echo "Launching ${TOTAL_REQUESTS} requests to KrakenD - Concurrency: $$c" ; \
		sleep 3 ; \
		docker run --platform linux/amd64 --rm --net=host  williamyeh/hey -n ${TOTAL_REQUESTS} -c $$c http://localhost:8080/test > ${PWD}/results/test_results_krakend_$${c}.txt ; \
	done ; \
	docker-compose down --volumes ; \
	cd ${PWD}

kong:
	@cd ${PWD}/gateways/kong ; \
	docker-compose up --build -d ; \
	for c in 50 100 250 500 1000 1500 2000 ; \
	do \
		echo "Launching ${TOTAL_REQUESTS} requests to Kong - Concurrency: $$c" ; \
		sleep 3 ; \
		docker run --platform linux/amd64 --rm --net=host  williamyeh/hey -n ${TOTAL_REQUESTS} -c $$c http://localhost:8080/test > ${PWD}/results/test_results_kong_$${c}.txt ; \
	done ; \
	docker-compose down --volumes ; \
	cd ${PWD}

tyk:
	@cd ${PWD}/gateways/tyk ; \
	docker-compose up --build -d ; \
	for c in 50 100 250 500 1000 1500 2000 ; \
	do \
		echo "Launching ${TOTAL_REQUESTS} requests to Tyk - Concurrency: $$c" ; \
		sleep 3 ; \
		docker run --platform linux/amd64 --rm --net=host  williamyeh/hey -n ${TOTAL_REQUESTS} -c $$c http://localhost:8080/test > ${PWD}/results/test_results_tyk_$${c}.txt ; \
	done ; \
	docker-compose down --volumes ; \
	cd ${PWD}

nginx:
	@cd ${PWD}/gateways/nginx ; \
	docker-compose up --build -d ; \
	for c in 50 100 250 500 1000 1500 2000 ; \
	do \
		echo "Launching ${TOTAL_REQUESTS} requests to Nginx - Concurrency: $$c" ; \
		sleep 3 ; \
		docker run --platform linux/amd64 --rm --net=host  williamyeh/hey -n ${TOTAL_REQUESTS} -c $$c http://localhost:8080/test > ${PWD}/results/test_results_nginx_$${c}.txt ; \
	done ; \
	docker-compose down --volumes ; \
	cd ${PWD}

stop:
	@cd ${PWD}/gateways/krakend ; docker-compose down --volumes ; \
	cd ${PWD}/gateways/kong ; docker-compose down --volumes ; \
	cd ${PWD}/gateways/tyk ; docker-compose down --volumes ; \
	cd ${PWD}/gateways/nginx ; docker-compose down --volumes ; \
	cd ${PWD}
