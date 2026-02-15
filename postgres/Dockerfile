FROM postgres:16-alpine
COPY schema.sql /docker-entrypoint-initdb.d/
COPY testdata.sql /docker-entrypoint-initdb.d/
