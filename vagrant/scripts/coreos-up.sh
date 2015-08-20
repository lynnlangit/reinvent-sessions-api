#!/bin/bash

systemctl stop ap
systemctl stop db

systemctl start db
systemctl start ap
