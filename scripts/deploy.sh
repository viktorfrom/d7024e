#!/bin/bash

# Run this from root NOT inside /scripts
docker stack rm kadlab
docker build cmd/client/. -t apiclient
docker build . -t kadlab
docker stack deploy --compose-file docker-compose.yml kadlab