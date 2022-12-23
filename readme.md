# Shorty

Shorty is a url shortening service. 

# Docker Compose

```
$ ./docker-compose up -d
```

Once the stack is up and running you can seed the database with unique keys. This must be done exactly once.

```
$ ./docker-compose exec server sh
# psql $POSTGRES_URL -c '\copy urls (key) from /sql/seed.csv csv;'
```
