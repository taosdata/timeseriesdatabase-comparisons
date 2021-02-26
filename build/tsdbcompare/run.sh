#!/bin/bash


# Color setting
RED='\033[0;31m'
GREEN='\033[1;32m'
GREEN_DARK='\033[0;32m'
GREEN_UNDERLINE='\033[4;32m'
NC='\033[0m'
# docker rm -f `docker ps -a -q`
#set -x
echo
echo "---------------Generating Data-----------------"
echo
echo "Prepare data for InfluxDB...."
#bin/bulk_data_gen -seed 123 -format influx-bulk -scale-var 100 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z" >data/influx.dat
bin/bulk_data_gen -seed 123 -format influx-bulk -sampling-interval 1s -scale-var 10 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z" >data/influx.dat

echo 
echo "Prepare data for TDengine...."
#bin/bulk_data_gen -seed 123 -format tdengine -tdschema-file config/TDengineSchema.toml -scale-var 100 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z"  > data/tdengine.dat
bin/bulk_data_gen -seed 123 -format tdengine -sampling-interval 1s -tdschema-file config/TDengineSchema.toml -scale-var 10 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z"  > data/tdengine.dat



#docker network create --ip-range 172.15.1.255/24 --subnet 172.15.1.1/16 tsdbcomp >>/dev/null 2>&1


#TDENGINE=`docker run -d --net tsdbcomp --ip 127.0.0.1 -p 6030:6030 -p 6030:6030 -p 6031:6031 -p 6032:6032 -p 6033:6033 -p 6034:6034 -p 6035:6035 -p 6036:6036 -p 6037:6037 -p 6038:6038 -p 6039:6039 tdengine/tdengine:1.6.4.5` 
echo
echo "------------------Writing Data-----------------"
echo
sleep 5
echo
echo -e "Start test TDengine, result in ${GREEN}Green line${NC}"

taos -s 'drop database devops;'
TDENGINERES=`cat data/tdengine.dat |bin/bulk_load_tdengine --url 127.0.0.1:0 --batch-size 5000   -do-load -report-tags n1 -workers 20 -fileout=false| grep loaded`
#TDENGINERES=`cat data/tdengine.dat |gunzip|bin/bulk_load_tdengine --url 127.0.0.1:0 --batch-size 300   -do-load -report-tags n1 -workers 10 -fileout=false| grep loaded`
echo
echo -e "${GREEN}TDengine writing result:${NC}"
echo -e "${GREEN}$TDENGINERES${NC}"
DATA=`echo $TDENGINERES|awk '{print($2)}'`
TMP=`echo $TDENGINERES|awk '{print($5)}'`
TDWTM=`echo ${TMP%s*}`


#INFLUX=`docker run -d -p 8086:8086 --net tsdbcomp --ip 127.0.0.1 influxdb` >>/dev/null 2>&1
sleep 10
echo
echo -e "Start test InfluxDB, result in ${GREEN}Green line${NC}"

curl 'http://localhost:8086/query?q=drop%20database%20benchmark' -X POST

INFLUXRES=`cat data/influx.dat  |bin/bulk_load_influx --batch-size=5000 --workers=20 --urls="http://127.0.0.1:8086" | grep loaded`


echo
echo -e "${GREEN}InfluxDB writing result:${NC}"
echo -e "${GREEN}$INFLUXRES${NC}"

TMP=`echo $INFLUXRES|awk '{print($5)}'`
IFWTM=`echo ${TMP%s*}`

echo
echo "------------------Querying Data-----------------"
echo

sleep 10
echo 
echo  "start query test, query max from 8 hosts group by 1 hour, TDengine"
echo

#Test case 1
#测试用例1，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') ;
# a,b,c,d,e,f,g,h are random 8 numbers.
TDQS1=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-all -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="127.0.0.1:0" -workers 50 -threads 50 -print-interval 0 -http-client-type cgo|grep wall`
echo
echo -e "${GREEN}TDengine query test case 1 result:${NC}"
echo -e "${GREEN}$TDQS1${NC}"
TMP=`echo $TDQS1|awk '{print($4)}'`
TDQ1=`echo ${TMP%s*}`

#Test case 2
#测试用例2，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1小时为粒度，查询每1小时的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') interval(1h);
# a,b,c,d,e,f,g,h are random 8 numbers
TDQS2=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-allbyhr -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="127.0.0.1:0" -workers 50 -threads 50 -print-interval 0 -http-client-type cgo|grep wall`

echo
echo -e "${GREEN}TDengine query test case 2 result:${NC}"
echo -e "${GREEN}$TDQS2${NC}"
TMP=`echo $TDQS2|awk '{print($4)}'`
TDQ2=`echo ${TMP%s*}`

#Test case 3
#测试用例3，测试用例3，随机查询12个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以10分钟为粒度，查询每10分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =12 hour
TDQS3=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-12-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="127.0.0.1:0" -workers 50 -threads 50 -print-interval 0 -http-client-type cgo|grep wall`
echo
echo -e "${GREEN}TDengine query test case 3 result:${NC}"
echo -e "${GREEN}$TDQS3${NC}"
TMP=`echo $TDQS3|awk '{print($4)}'`
TDQ3=`echo ${TMP%s*}`

#Test case 4
#测试用例4，随机查询1个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1分钟为粒度，查询每1分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =1 hours
TDQS4=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="127.0.0.1:0" -workers 50 -threads 50 -print-interval 0 -http-client-type cgo|grep wall`
echo
echo -e "${GREEN}TDengine query test case 4 result:${NC}"
echo -e "${GREEN}$TDQS4${NC}"
TMP=`echo $TDQS4|awk '{print($4)}'`
TDQ4=`echo ${TMP%s*}`

sleep 10

echo 
echo  "start query test, query max from 8 hosts group by 1hour, Influxdb"
echo
#Test case 1
#测试用例1，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') ;
# a,b,c,d,e,f,g,h are random 8 numbers.
IFQS1=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-all -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://127.0.0.1:8086"  -workers 50 -print-interval 0|grep wall`
echo -e "${GREEN}InfluxDB query test case 1 result:${NC}"
echo -e "${GREEN}$IFQS1${NC}"
TMP=`echo $IFQS1|awk '{print($4)}'`
IFQ1=`echo ${TMP%s*}`
#Test case 2
#测试用例2，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1小时为粒度，查询每1小时的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') interval(1h);
# a,b,c,d,e,f,g,h are random 8 numbers
IFQS2=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-allbyhr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://127.0.0.1:8086"  -workers 50 -print-interval 0|grep wall`
echo -e "${GREEN}InfluxDB query test case 2 result:${NC}"
echo -e "${GREEN}$IFQS2${NC}"
TMP=`echo $IFQS2|awk '{print($4)}'`
IFQ2=`echo ${TMP%s*}`
#Test case 3
#测试用例3，测试用例3，随机查询12个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以10分钟为粒度，查询每10分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =12 hour
#INFLUXQUERY=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://127.0.0.1:8086"  -workers 50 -print-interval 0|grep wall`
IFQS3=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-12-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://127.0.0.1:8086"  -workers 50 -print-interval 0|grep wall`
echo -e "${GREEN}InfluxDB query test case 3 result:${NC}"
echo -e "${GREEN}$IFQS3${NC}"
TMP=`echo $IFQS3|awk '{print($4)}'`
IFQ3=`echo ${TMP%s*}`
#Test case 4
#测试用例4，随机查询1个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1分钟为粒度，查询每1分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =1 hours
#INFLUXQUERY=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://127.0.0.1:8086"  -workers 50 -print-interval 0|grep wall`
IFQS4=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://127.0.0.1:8086"  -workers 50 -print-interval 0|grep wall`
echo -e "${GREEN}InfluxDB query test case 4 result:${NC}"
echo -e "${GREEN}$IFQS4${NC}"
TMP=`echo $IFQS4|awk '{print($4)}'`
IFQ4=`echo ${TMP%s*}`


echo
echo
echo    "======================================================"
echo    "             tsdb performance comparision             "
echo    "======================================================"
echo -e "       Writing $DATA records test takes:          "
printf  "       InfluxDB           |       %-4.2f Seconds    \n" $IFWTM 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDWTM
echo    "------------------------------------------------------"
echo    "                   Query test cases:                "
echo    " case 1: select the max(value) from all data    "
echo    " filtered out 8 hosts                                 "
echo    "       Query test case 1 takes:                      "
printf  "       InfluxDB           |       %-4.2f Seconds    \n" $IFQ1 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDQ1
echo    "------------------------------------------------------"
echo    " case 2: select the max(value) from all data          "
echo    " filtered out 8 hosts with an interval of 1 hour     "
echo    " case 2 takes:                                       "
printf  "       InfluxDB           |       %-4.2f Seconds    \n" $IFQ2 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDQ2
echo    "------------------------------------------------------"
echo    " case 3: select the max(value) from random 12 hours"
echo    " data filtered out 8 hosts with an interval of 10 min         "
echo    " filtered out 8 hosts interval(1h)                   "
echo    " case 3 takes:                                       "
printf  "       InfluxDB           |       %-4.2f Seconds    \n" $IFQ3 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDQ3
echo    "------------------------------------------------------"
echo    " case 4: select the max(value) from random 1 hour data  "
echo    " data filtered out 8 hosts with an interval of 1 min         "
echo    " case 4 takes:                                        "
printf  "       InfluxDB           |       %-4.2f Seconds    \n" $IFQ4 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDQ4
echo    "------------------------------------------------------"
echo
docker stop $INFLUX >>/dev/null 2>&1
docker container rm -f $INFLUX >>/dev/null 2>&1
docker stop $TDENGINE >>/dev/null 2>&1
docker container rm -f $TDENGINE >>/dev/null 2>&1
docker network rm tsdbcomp >>/dev/null 2>&1
#bulk_query_gen/bulk_query_gen  -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_influxdb/query_benchmarker_influxdb  -urls="http://172.26.89.231:8086" 
#bulk_query_gen/bulk_query_gen  -format tdengine -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_tdengine/query_benchmarker_tdengine  -urls="http://172.26.89.231:6030" 
