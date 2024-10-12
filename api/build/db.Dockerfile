FROM postgres:16.3

COPY ./internal/database/ /docker-entrypoint-initdb.d/
