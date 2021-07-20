.PHONY: all krakend kong tyk nginx stop

all: krakend kong tyk nginx

krakend:
	cd ${PWD}/gateways/krakend ; \
	docker-compose up --build -d ; \
	for c in 50 100 250 500 1000 ; \
	do \
		echo "Launching 100K requests to KrakenD - Concurrency: $$c" ; \
		docker run --rm --net=host  williamyeh/hey -n 100000 -c $$c http://localhost:8080/test/ > ${PWD}/results/test_results_krakend_$${c}.txt ; \
	done ; \
	docker-compose down --volumes ; \
	cd ${PWD}

kong:
	cd ${PWD}/gateways/kong ; \
	docker-compose up --build -d ; \
	for c in 50 100 250 500 1000 ; \
	do \
		echo "Launching 100K requests to Kong - Concurrency: $$c" ; \
		docker run --rm --net=host  williamyeh/hey -n 100000 -c $$c http://localhost:8080/test/answer.json > ${PWD}/results/test_results_kong_$${c}.txt ; \
	done ; \
	docker-compose down --volumes ; \
	cd ${PWD}

tyk:
	cd ${PWD}/gateways/tyk ; \
	docker-compose up --build -d ; \
	for c in 50 100 250 500 1000 ; \
	do \
		echo "Launching 100K requests to Tyk - Concurrency: $$c" ; \
		docker run --rm --net=host  williamyeh/hey -n 100000 -c $$c http://localhost:8080/test/ > ${PWD}/results/test_results_tyk_$${c}.txt ; \
	done ; \
	docker-compose down --volumes ; \
	cd ${PWD}

nginx:
	cd ${PWD}/gateways/nginx ; \
	docker-compose up --build -d ; \
	for c in 50 100 250 500 1000 ; \
	do \
		echo "Launching 100K requests to Nginx - Concurrency: $$c" ; \
		docker run --rm --net=host  williamyeh/hey -n 100000 -c $$c http://localhost:8080/test/ > ${PWD}/results/test_results_nginx_$${c}.txt ; \
	done ; \
	docker-compose down --volumes ; \
	cd ${PWD}

stop:
	cd ${PWD}/gateways/krakend ; docker-compose down --volumes ; \
	cd ${PWD}/gateways/kong ; docker-compose down --volumes ; \
	cd ${PWD}/gateways/tyk ; docker-compose down --volumes ; \
	cd ${PWD}/gateways/nginx ; docker-compose down --volumes ; \
	cd ${PWD}
