#!/bin/bash

systemctl stop monit
systemctl stop ap
systemctl stop db

docker build -f /home/core/share/docker/Dockerfile.base -t supinf/reinvent-sessions-api:base .
docker pull deangiberson/aws-dynamodb-local
docker pull pottava/docker-webui:latest
