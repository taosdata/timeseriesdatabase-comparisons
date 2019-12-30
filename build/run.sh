#!/bin/bash

# Color setting
RED='\033[0;31m'
GREEN='\033[1;32m'
GREEN_DARK='\033[0;32m'
GREEN_UNDERLINE='\033[4;32m'
NC='\033[0m'

#set -x

echo "Prepare data for InfluxDB...."
bin/bulk_data_gen -format influx-bulk -scale-var 100 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z"  >data/influx.dat

echo 
echo "Prepare data for TDengine...."
bin/bulk_data_gen -format tdengine -tdschema-file config/TDengineSchema.toml -scale-var 100 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z"  > data/tdengine.dat

docker pull influxdb >>/dev/null 2>&1
docker pull tdengine/tdengine >>/dev/null 2>&1

docker run -d influxdb >>/dev/null 2>&1
docker run -d tdengine/tdengine >>/dev/null 2>&1

docker network create --ip-range 172.15.1.255/24 --subnet 172.15.1.1/16 tsdbcomp >>/dev/null 2>&1

INFLUX=`docker run -d -p 8086:8086 --net tsdbcomp --ip 172.15.1.5 influxdb` >>/dev/null 2>&1

sleep 5
echo
echo -e "Start test InfluxDB, result in ${GREEN}Green line${NC}"

INFLUXRES=`cat data/influx.dat |bin/bulk_load_influx --batch-size=5000 --workers=10 --urls="http://172.15.1.5:8086" | grep loaded`
echo
echo -e "${GREEN}InfluxDB writing result:${NC}"
echo -e "${GREEN}$INFLUXRES${NC}"
docker stop $INFLUX >>/dev/null 2>&1

TDENGINE=`docker run -d --net tsdbcomp --ip 172.15.1.6 -p 6030:6030  tdengine/tdengine`

sleep 5
echo
echo -e "Start test TDengine, result in ${GREEN}Green line${NC}"
#cat data/tdengine.dat | gunzip|./bulk_load_tdengine  --batch-size=5000 --do-load true --report-tags=n1 --workers=10 --url="172.15.1.6:0" | grep loaded

TDENGINERES=`cat data/tdengine.dat |bin/bulk_load_tdengine --url 172.15.1.6:0 --batch-size 300   -do-load -report-tags n1 -workers 10 -fileout=false| grep loaded`
echo
echo -e "${GREEN}TDengine writing result:${NC}"
echo -e "${GREEN}$TDENGINERES${NC}"
docker stop $TDENGINE >>/dev/null 2>&1

docker network rm tsdbcomp >>/dev/null 2>&1
