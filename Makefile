all: build

build:
	docker build ./database -t online-judge-database

clean:
	-docker stop online-judge-database
	-docker stop pgadmin

pgadmin:
	# Open up http://127.0.0.1:5433 to access pgAdmin
	docker run --rm --name pgadmin -p 5433:80 -e 'PGADMIN_DEFAULT_EMAIL=admin@admin.com' -e 'PGADMIN_DEFAULT_PASSWORD=admin' -d dpage/pgadmin4

postgres:
	docker build ./database -t online-judge-database
	docker run --rm -d --name online-judge-database -p 5432:5432 online-judge-database