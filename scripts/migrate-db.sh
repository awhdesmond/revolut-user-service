#!/bin/sh

for db in postgres postgres-test
do
    flyway \
        -url="jdbc:postgresql://127.0.0.1:5432/${db}" \
        -user="postgres" \
        -password="postgres" \
        -locations="filesystem:db/migrations/" \
        -validateMigrationNaming=true \
        migrate
done
