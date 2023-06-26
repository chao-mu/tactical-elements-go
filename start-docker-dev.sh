#!/usr/bin/env bash

source .env
docker stop te-db
docker rm te-db
docker run -p 127.0.0.1:5432:5432 --name te-db  -e POSTGRES_PASSWORD=$TE_PASS -e POSTGRES_USER=$TE_USER -d postgres
