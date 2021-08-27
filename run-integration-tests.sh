#!/bin/bash

docker build -t petstore .
docker-compose up --force-recreate -dV

echo "Waiting for the application to be healthy"
exitCode=1
while [ $exitCode -ne 0 ]
do
	sleep 1
	curl -f http://localhost:5000/health
	exitCode=$?
done

go test -count=1 ./integrationtests
