FROM postgres:11.5-alpine

ENV POSTGRES_USER docker
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB links

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
CMD ["postgres"]

# docker build --tag=linksdb:latest .
# docker run -d --rm -p 5432:5432 linksdb
