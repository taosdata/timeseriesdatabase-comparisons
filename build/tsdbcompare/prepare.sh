#!/bin/bash
set -x

#安装docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh


#安装tdengine客户端
tar -zxf TDengine-client-1.6.4.5.tar.gz
cd TDengine-client-1.6.4.5
./install_client.sh
cd ..

#拉取influxdb和tdengine的镜像
docker pull influxdb 
docker pull tdengine/tdengine:v1.6.4.5.c 
