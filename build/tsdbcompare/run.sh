#!/bin/bash


# Color setting
RED='\033[0;31m'
GREEN='\033[1;32m'
GREEN_DARK='\033[0;32m'
GREEN_UNDERLINE='\033[4;32m'
NC='\033[0m'

#set -x
echo
echo "Prepare data for InfluxDB...."
#bin/bulk_data_gen -seed 123 -format influx-bulk -scale-var 100 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z" >data/influx.dat
bin/bulk_data_gen -seed 123 -format influx-bulk -sampling-interval 1s -scale-var 10 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z" >data/influx.dat

echo 
echo "Prepare data for TDengine...."
#bin/bulk_data_gen -seed 123 -format tdengine -tdschema-file config/TDengineSchema.toml -scale-var 100 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z"  > data/tdengine.dat
bin/bulk_data_gen -seed 123 -format tdengine -sampling-interval 1s -tdschema-file config/TDengineSchema.toml -scale-var 10 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z"  > data/tdengine.dat



docker network create --ip-range 172.15.1.255/24 --subnet 172.15.1.1/16 tsdbcomp >>/dev/null 2>&1


TDENGINE=`docker run -d --net tsdbcomp --ip 172.15.1.6 -p 6030:6030 -p 6020:6020 -p 6031:6031 -p 6032:6032 -p 6033:6033 -p 6034:6034 -p 6035:6035 -p 6036:6036 -p 6037:6037 -p 6038:6038 -p 6039:6039 tdengine/tdengine` 

sleep 5
echo
echo -e "Start test TDengine, result in ${GREEN}Green line${NC}"

TDENGINERES=`cat data/tdengine.dat |bin/bulk_load_tdengine --url 172.15.1.6:0 --batch-size 300   -do-load -report-tags n1 -workers 20 -fileout=false| grep loaded`
#TDENGINERES=`cat data/tdengine.dat |gunzip|bin/bulk_load_tdengine --url 172.15.1.6:0 --batch-size 300   -do-load -report-tags n1 -workers 10 -fileout=false| grep loaded`
echo
echo -e "${GREEN}TDengine writing result:${NC}"
echo -e "${GREEN}$TDENGINERES${NC}"
DATA=`echo $TDENGINERES|awk '{print($2)}'`
TMP=`echo $TDENGINERES|awk '{print($5)}'`
TDWTM=`echo ${TMP%s*}`


INFLUX=`docker run -d -p 8086:8086 --net tsdbcomp --ip 172.15.1.5 influxdb` >>/dev/null 2>&1
sleep 10
echo
echo -e "Start test InfluxDB, result in ${GREEN}Green line${NC}"


INFLUXRES=`cat data/influx.dat  |bin/bulk_load_influx --batch-size=5000 --workers=20 --urls="http://172.15.1.5:8086" | grep loaded`


echo
echo -e "${GREEN}InfluxDB writing result:${NC}"
echo -e "${GREEN}$INFLUXRES${NC}"

TMP=`echo $INFLUXRES|awk '{print($5)}'`
IFWTM=`echo ${TMP%s*}`

sleep 10
echo 
echo  "start query test, query max from 8 hosts group by 1hour, TDengine"
echo

#Test case 1
#测试用例1，查询一天的数据中的最大值，用8个hostname标签进行过滤，以1小时为颗粒查询
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') interval(1h);
# a,b,c,d,e,f,g,h are random 8 numbers.
TDQUERY=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-allbyhr -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="http://172.15.1.6:6020" -workers 50 -print-interval 0|grep wall`

#Test case 2
#测试用例2，查询12小时的数据中的最大值，用8个hostname标签进行过滤，以10分钟为颗粒查询
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') where time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =12 hours
#TDQUERY=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-12-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="http://172.15.1.6:6020" -workers 50 -print-interval 0|grep wall`

#Test case 3
#测试用例3，查询1个小时的数据，用1个hostname进行过滤，以1分钟为颗粒查询
#select max(usage_user) from cpu where(hostname='host_a') where time >x and time <y interval(10m);
# a are random number, y-x =1 hours
#TDQUERY=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 1-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="http://172.15.1.6:6020" -workers 50 -print-interval 0|grep wall`

echo
echo -e "${GREEN}TDengine query result:${NC}"
echo -e "${GREEN}$TDQUERY${NC}"
TMP=`echo $TDQUERY|awk '{print($4)}'`
TDQTM=`echo ${TMP%s*}`

sleep 10

echo 
echo  "start query test, query max from 8 hosts group by 1hour, Influxdb"
echo
#Test case 1
#测试用例1，查询一天的数据中的最大值，用8个hostname标签进行过滤，以1小时为颗粒查询
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') interval(1h);
# a,b,c,d,e,f,g,h are random 8 numbers.
INFLUXQUERY=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-allbyhr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://172.15.1.5:8086"  -workers 50 -print-interval 0|grep wall`

#Test case 2
#测试用例2，查询12小时的数据中的最大值，用8个hostname标签进行过滤，以10分钟为颗粒查询
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') where time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =12 hours
#INFLUXQUERY=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-12-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://172.15.1.5:8086"  -workers 50 -print-interval 0|grep wall`

#Test case 3
#测试用例3，查询1个小时的数据，用1个hostname进行过滤，以1分钟为颗粒查询
#select max(usage_user) from cpu where(hostname='host_a') where time >x and time <y interval(10m);
# a are random number, y-x =1 hours
#INFLUXQUERY=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://172.15.1.5:8086"  -workers 50 -print-interval 0|grep wall`

echo
echo -e "${GREEN}InfluxDB query result:${NC}"
echo -e "${GREEN}$INFLUXQUERY${NC}"

TMP=`echo $INFLUXQUERY|awk '{print($4)}'`
IFQTM=`echo ${TMP%s*}`

echo
echo
echo    "-----------------------------------------------"
echo    "         tsdb performance comparision          "
echo    "-----------------------------------------------"
echo -e "         writing $DATA records takes:          "
printf  "       InfluxDB       |    %-4.2f Seconds    \n" $IFWTM 
printf  "       TDengine       |    %-4.2f Seconds    \n" $TDWTM
echo    "-----------------------------------------------"
echo    "1000 queries: select max(usage_user) groupby 1h"
echo    "takes:                                         "
printf  "       InfluxDB       |    %-4.2f Seconds    \n" $IFQTM 
printf  "       TDengine       |    %-4.2f Seconds    \n" $TDQTM
echo    "-----------------------------------------------"
echo
docker stop $INFLUX >>/dev/null 2>&1
docker container rm -f $INFLUX >>/dev/null 2>&1
docker stop $TDENGINE >>/dev/null 2>&1
docker container rm -f $TDENGINE >>/dev/null 2>&1
docker network rm tsdbcomp >>/dev/null 2>&1
#bulk_query_gen/bulk_query_gen  -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_influxdb/query_benchmarker_influxdb  -urls="http://172.26.89.231:8086" 
#bulk_query_gen/bulk_query_gen  -format tdengine -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_tdengine/query_benchmarker_tdengine  -urls="http://172.26.89.231:6020" 
