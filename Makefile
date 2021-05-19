all: build

build:
	make clean
	make network
	make postgres

clean:
	-docker stop online-judge-database
	-docker stop pgadmin
	-docker network rm online-judge-network

pgadmin:
	# Make sure the network online-judge-network is created, will cause an error if not created, use `make network`
	# Open up http://127.0.0.1:5433 to access pgAdmin
	docker run --rm --name pgadmin --network online-judge-network -p 5433:80 -e 'PGADMIN_DEFAULT_EMAIL=admin@admin.com' -e 'PGADMIN_DEFAULT_PASSWORD=admin' -d dpage/pgadmin4

postgres:
	docker build ./database -t online-judge-database
	docker run --rm -d --name online-judge-database --network online-judge-network -p 5432:5432 online-judge-database

network:
	docker network create online-judge-network --driver bridge