#!/bin/bash
docker-compose -f docker-compose-zxy.yml down
docker ps -a|grep ago | awk '{print $1}' |xargs docker rm
docker rmi deploy_ui
docker-compose -f docker-compose-zxy.yml up -d
