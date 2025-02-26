#! /bin/bash

curl -H 'Content-Type: application/json' debezium:8083/connectors --data '
{
  "name": "postgres-connector",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "plugin.name": "pgoutput",
    "database.hostname": "postgres",
    "database.port": "5432",
    "database.user": "postgres",
    "database.password": "password",
    "database.dbname" : "postgres",
    "database.server.name": "postgres",
    "table.include.list": "public.urls",
    "topic.prefix" : "dbz"
  }
}'