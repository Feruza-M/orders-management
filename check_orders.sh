#!/bin/bash

containers=(
orders-app
orders-db
orders-nginx
)

echo "Checking containers..."

for c in "${containers[@]}"
do
 if docker ps --format '{{.Names}}' | grep -q $c
 then
   echo "$c running OK"
 else
   echo "$c NOT running"
   exit 1
 fi
done

echo "Checking app..."

curl -f http://localhost:8082/api/health && \
echo "Application healthy" || \
echo "Application failed"
