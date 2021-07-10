#!/bin/bash

echo "-----Auto deployment start!-----"

docker login --username=usr --password=passwd registry-vpc.cn-beijing.aliyuncs.com

echo "-----docker pull-----"

docker pull registry-vpc.cn-beijing.aliyuncs.com/usr/nsaop-backend:latest

echo "-----docker run-----"

docker kill nsaop-backend
docker rm nsaop-backend

docker run --name nsaop-backend -d -p 10080:10080 \
-v /etc/localtime:/etc/localtime \
-v /home/config/config.json:/app/config/config.json \
registry-vpc.cn-beijing.aliyuncs.com/usr/nsaop-backend:latest \
./nsaop-backend --email=true

python3 del_old_back_img.py

echo "----- Auto deployment end! -----"