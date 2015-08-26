#!/bin/bash

systemctl stop ap
systemctl stop db

docker pull supinf/reinvent-sessions-api:base
docker pull deangiberson/aws-dynamodb-local
