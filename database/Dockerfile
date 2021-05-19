FROM postgres
ENV POSTGRES_DB online-judge
ENV POSTGRES_PASSWORD postgres
COPY startup.sql /docker-entrypoint-initdb.d/
