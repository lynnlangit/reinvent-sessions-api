#!/bin/bash

systemctl stop monit
systemctl stop ap
systemctl stop db

systemctl start db
systemctl start ap
systemctl start monit
