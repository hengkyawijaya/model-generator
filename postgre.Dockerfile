FROM postgres:10.6-alpine

COPY ./migration/*.sql /docker-entrypoint-initdb.d/

