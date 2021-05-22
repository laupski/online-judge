all: build

build:
	docker build ./database -t online-judge-database
	docker build . -t online-judge-web
	docker build ./messaging -t online-judge-messaging

clean:
	-docker stop online-judge-database
	-docker stop pgadmin
	-docker stop online-judge-messaging

pgadmin:
	# Open up http://127.0.0.1:5433 to access pgAdmin
	docker run --rm --name pgadmin -p 5433:80 -e 'PGADMIN_DEFAULT_EMAIL=admin@admin.com' -e 'PGADMIN_DEFAULT_PASSWORD=admin' -d dpage/pgadmin4

postgres:
	docker build ./database -t online-judge-database
	docker run --rm -d --name online-judge-database -p 5432:5432 online-judge-database

rabbitmq:
	# Open up http://127.0.0.1:15672 to view RabbitMQ
	docker build ./messaging -t online-judge-messaging
	docker run --rm -d -p 15672:15672 -p 5672:5672 --name online-judge-rabbitmq online-judge-messaging

judge-rce:
	docker build . -t online-judge-web
	docker run --rm -d --name online-judge-rce online-judge-web online-judge judge local

debug:
	go build -gcflags="all=-N -l"