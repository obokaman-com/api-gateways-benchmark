.PHONY: all krakend tyk nginx

all: krakend tyk nginx

krakend:
	cd ${PWD}/gateways/krakend \
		&& docker-compose up --build -d && docker run --rm --net=host  williamyeh/hey -n 100000 -c 1000 http://localhost:8080/test/ > ${PWD}/results/test_results_krakend.txt \
		&& docker-compose down --volumes

tyk:
	cd ${PWD}/gateways/tyk \
		&& docker-compose up --build -d && docker run --rm --net=host  williamyeh/hey -n 100000 -c 1000 http://localhost:8080/test/ > ${PWD}/results/test_results_tyk.txt \
		&& docker-compose down --volumes

nginx:
	cd ${PWD}/gateways/nginx \
		&& docker-compose up --build -d && docker run --rm --net=host  williamyeh/hey -n 100000 -c 1000 http://localhost:8080/test/ > ${PWD}/results/test_results_nginx.txt \
		&& docker-compose down --volumes
