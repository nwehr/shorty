FROM golang:1.18
WORKDIR /app
COPY . .
RUN make server

FROM alpine
WORKDIR /app
RUN apk add --no-cache ca-certificates \
    && update-ca-certificates \
    && apk add postgresql-client
COPY --from=0 /app/shorty-server /app/sql/migrate.sh /app/
COPY --from=0 /app/templates /app/templates
COPY --from=0 /app/sql/migrations /sql/migrations
COPY --from=0 /app/sql/seed/seed.csv /sql/
EXPOSE 8080 8080
CMD ["/bin/sh", "-c", "/app/migrate.sh /sql/migrations; /app/shorty-server"]