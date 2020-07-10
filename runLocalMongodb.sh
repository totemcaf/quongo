#!/usr/bin/env bash

# cat pw.txt |  docker login --username carlosfau  --password-stdin

docker run -d --name mongo \
    -e MONGO_INITDB_ROOT_USERNAME=quongo \
    -e MONGO_INITDB_ROOT_PASSWORD=quongosecret \
    -p 27017:27017 \
    mongo:3.6
