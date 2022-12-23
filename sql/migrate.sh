#!/bin/sh
DIR=$1

if [ -z "$1" ]; then
    DIR=sql/migrations
fi

for FILE in $DIR/tables/*.sql
do
    psql ${POSTGRES_URL} < $FILE
done


for FILE in $DIR/columns/*.sql
do
    psql ${POSTGRES_URL} < $FILE
done