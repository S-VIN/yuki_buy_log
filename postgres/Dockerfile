FROM postgres:15-alpine
COPY schema.sql /docker-entrypoint-initdb.d/
COPY testdata.sql /docker-entrypoint-initdb.d/
