#!/bin/bash

systemctl stop ap
systemctl stop db

docker build -f /home/core/share/docker/Dockerfile.base -t supinf/reinvent-sessions:base .
docker pull deangiberson/aws-dynamodb-local
