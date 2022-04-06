# Just a docker image for workflows so that postgres table is already created
FROM postgres:latest
COPY CreateTable.sql /docker-entrypoint-initdb.d/