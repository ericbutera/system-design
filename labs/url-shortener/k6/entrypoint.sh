#!/bin/sh

for file in $(ls /app/*.js); do
  echo "Running $file"
  k6 run --quiet $file
done