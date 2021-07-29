#!/bin/bash


# Color setting
RED='\033[0;31m'
GREEN='\033[1;32m'
GREEN_DARK='\033[0;32m'
GREEN_UNDERLINE='\033[4;32m'
NC='\033[0m'

workers=16
interface='cgo'
gene=0
add='bschang1'
interval='10s'
scale=100
st='2018-01-01T00:00:00Z'
et='2018-01-02T00:00:00Z'

query=1000

while getopts "w:n:g:a:i:s:t:e:" opt
do
    case $opt in
        w)
        echo "workers:$OPTARG"
        workers=$OPTARG
        ;;
        n)
        echo "TD's interface:$OPTARG"
        interface=$OPTARG
        ;;
        g)
        echo "whether generate data:$OPTARG"
        gene=$OPTARG
        ;;
        i)
        echo "sampling interval:$OPTARG"
        interval=$OPTARG
        ;;
        a)
        echo "address:$OPTARG"
        add=$OPTARG
        ;;
        s)
        echo "scale-var:$OPTARG"
        scale=$OPTARG
        ;;
        t)
        echo "timestamp-start:$OPTARG"
        st=$OPTARG
        ;;
        e)
        echo "timestamp-end:$OPTARG"
        et=$OPTARG
        ;;
        ?)
        echo    "======================================================"
        echo    "w | query workers"
        echo    "------------------------------------------------------"
        echo    "n | TD's interface(cgo,fast)"
        echo    "------------------------------------------------------"
        echo    "i | sampling interval(default:10s)"
        echo    "------------------------------------------------------"
        echo    "a | address of TD & influx"
        echo    "------------------------------------------------------"
        echo    "s | scale-var(default:100)"
        echo    "------------------------------------------------------"
        echo    "t | timestamp-start(default:'2018-01-01T00:00:00Z')"
        echo    "------------------------------------------------------"
        echo    "e | timestamp-end(default:'2018-01-02T00:00:00Z')"
        echo    "------------------------------------------------------"
        echo    "g | genate data(0:no ,1:yes ,default:0)"
        echo    "======================================================"
        exit 1;;
    esac
done

echo "variables :"
echo "generate data: scale-var: $scale ,interval: $interval ,timestamp-start: $st ,timestamp-stop: $et"
echo "workers:$workers ,TD's interface: $interface ,generate data: $gene , address: $add"

if [[ $gene == 1 ]];then 
    if [ ! -d "data" ]; then
        mkdir data
    fi

    ssh root@$add << eeooff
    systemctl stop taosd
    rm -rf /data/lib/taos/*
    sudo service cassandra start
    sleep 30
    echo 'drop keyspace if exists measurements;' | cqlsh $add
    echo 1 > /proc/sys/vm/drop_caches
    sudo service cassandra stop
    systemctl start taosd
    rm -rf /data/cassandra/data/measurements
    sudo service cassandra start
    sleep 10
    exit
eeooff


    echo
    echo "---------------Generating && Inserting Data-----------------"
    echo
    echo "Prepare data for Cassandra...."
    bin/bulk_data_gen -seed 123 -format cassandra -sampling-interval $interval -scale-var $scale -use-case devops -timestamp-start "$st" -timestamp-end "$et" >data/cassandra.dat
    cat data/cassandra.dat  |bin/bulk_load_cassandra --batch-size=2000 --workers=16 --url $add | grep loaded

    echo 
    echo "Prepare data for TDengine...."
    bin/bulk_data_gen -seed 123 -format tdengine -sampling-interval $interval -tdschema-file config/TDengineSchema.toml -scale-var $scale -use-case devops -timestamp-start "$st" -timestamp-end "$et"  > data/tdengine.dat
    cat data/tdengine.dat |bin/bulk_load_tdengine --url $add --batch-size 2000  -do-load -report-tags n1 -workers 50 -fileout=false -http-api='false'  | grep loaded
fi

echo
echo "------------------Querying Data-----------------"
echo

ssh root@$add << eeooff
systemctl stop taosd
service cassandra stop
echo 1 > /proc/sys/vm/drop_caches
systemctl start taosd
exit
eeooff

sleep 30

scp monitor.sh $add:/root
echo 
echo  "start query test, query max from 8 hosts group by 1 hour, TDengine"
echo

#Test case 1
#测试用例1，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') ;
# a,b,c,d,e,f,g,h are random 8 numbers.
ssh root@bschang1 -f /root/monitor.sh taosd >/dev/null 2>1
nohup ./monitor.sh taos >/dev/null 2>1 &
TDQS1=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-all -scale-var $scale -queries $query | bin/query_benchmarker_tdengine  -urls=$add -workers $workers -threads $workers -print-interval 0 -http-client-type $interface | grep wall`
echo
echo -e "${GREEN}TDengine query test case 1 result:${NC}"
echo -e "${GREEN}$TDQS1${NC}"
TMP=`echo $TDQS1|awk '{print($4)}'`
TDQ1=`echo ${TMP%s*}`
ssh root@bschang1 -f pkill -9 monitor.sh >/dev/null 2>1
nohup pkill -9 monitor.sh 2>1 &

#Test case 2
#测试用例2，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1小时为粒度，查询每1小时的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') interval(1h);
# a,b,c,d,e,f,g,h are random 8 numbers
ssh root@bschang1 -f /root/monitor.sh taosd >/dev/null 2>1
nohup ./monitor.sh taos >/dev/null 2>1 &
TDQS2=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-allbyhr -scale-var $scale -queries $query | bin/query_benchmarker_tdengine  -urls=$add -workers $workers -threads $workers -print-interval 0 -http-client-type $interface | grep wall`
echo
echo -e "${GREEN}TDengine query test case 2 result:${NC}"
echo -e "${GREEN}$TDQS2${NC}"
TMP=`echo $TDQS2|awk '{print($4)}'`
TDQ2=`echo ${TMP%s*}`
ssh root@bschang1 -f pkill -9 monitor.sh >/dev/null 2>1
nohup pkill -9 monitor.sh 2>1 &

#Test case 3
#测试用例3，测试用例3，随机查询12个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以10分钟为粒度，查询每10分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =12 hour
ssh root@bschang1 -f /root/monitor.sh taosd >/dev/null 2>1
nohup ./monitor.sh taos >/dev/null 2>1 &
TDQS3=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-12-hr -scale-var $scale -queries $query | bin/query_benchmarker_tdengine  -urls=$add -workers $workers -threads $workers -print-interval 0 -http-client-type $interface | grep wall`
echo
echo -e "${GREEN}TDengine query test case 3 result:${NC}"
echo -e "${GREEN}$TDQS3${NC}"
TMP=`echo $TDQS3|awk '{print($4)}'`
TDQ3=`echo ${TMP%s*}`
ssh root@bschang1 -f pkill -9 monitor.sh >/dev/null 2>1
nohup pkill -9 monitor.sh 2>1 &

#Test case 4
#测试用例4，随机查询1个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1分钟为粒度，查询每1分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =1 hours
ssh root@bschang1 -f /root/monitor.sh taosd >/dev/null 2>1
nohup ./monitor.sh taos >/dev/null 2>1 &
TDQS4=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-1-hr -scale-var $scale -queries $query | bin/query_benchmarker_tdengine  -urls=$add -workers $workers -threads $workers -print-interval 0 -http-client-type $interface | grep wall`
echo
echo -e "${GREEN}TDengine query test case 4 result:${NC}"
echo -e "${GREEN}$TDQS4${NC}"
TMP=`echo $TDQS4|awk '{print($4)}'`
TDQ4=`echo ${TMP%s*}`
ssh root@bschang1 -f pkill -9 monitor.sh >/dev/null 2>1
nohup pkill -9 monitor.sh 2>1 &


#Test case 5
#测试用例4，随机查询1个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1分钟为粒度，查询每1分钟的最大值
#为了与cassandra的查询方式保持一致，本测试将不会使用interval。 查询的方式将改为运行多个select语句，并依靠 ts> and ts< 模拟interval
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y;
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =1 hours
ssh root@bschang1 -f /root/monitor.sh taosd >/dev/null 2>1
nohup ./monitor.sh taos >/dev/null 2>1 &
TDQS5=`bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-1-hr-no-interval -scale-var $scale -queries $query | bin/query_benchmarker_tdengine  -urls=$add -workers $workers -threads $workers -print-interval 0 -http-client-type $interface | grep wall`
echo
echo -e "${GREEN}TDengine query test case 5 result:${NC}"
echo -e "${GREEN}$TDQS5${NC}"
TMP=`echo $TDQS5|awk '{print($4)}'`
TDQ5=`echo ${TMP%s*}`
ssh root@bschang1 -f pkill -9 monitor.sh >/dev/null 2>1
nohup pkill -9 monitor.sh 2>1 &

sleep 10

ssh root@$add << eeooff
systemctl stop taosd 
echo 1 > /proc/sys/vm/drop_caches
sudo service cassandra start
sleep 40
exit
eeooff

echo 
echo  "start query test, query max from 8 hosts group by 1hour, cassandra"
echo
#cassandra Test case 1
#测试用例1，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据的最大值。
#因为cassandra运行中每个worker最后的进程不会工作并马上结束，所以总查询数位query+worker
#SELECT max(usage_user) FROM measurements.cpu WHERE hostname in( 'host_a' ,  'host_b' ,  'host_c' ,  'host_d' ,  'host_e' ,  'host_f' ,  'host_g' ,  'host_h') ;
# a,b,c,d,e,f,g,h are random 8 numbers.
ssh root@bschang1 -f /root/monitor.sh java >/dev/null 2>1
nohup ./monitor.sh cqlsh >/dev/null 2>1 &
sum=$((query + workers))
IFQS1=`bin/bulk_query_gen  -seed 123 -format cassandra -query-type 8-host-all -scale-var $scale -queries $((query + workers)) | bin/query_benchmarker_cassandra  -url=$add  -workers $workers -print-interval 0 -aggregation-plan server| grep wall`
#`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-all -scale-var $scale -queries 5 | bin/query_benchmarker_influxdb  -urls="http://$add:8086"  -workers $workers -print-interval 0|grep wall`
echo -e "${GREEN}cassandra query test case 1 result:${NC}"
echo -e "${GREEN}$IFQS1${NC}"
TMP=`echo $IFQS1|awk '{print($4)}'`
IFQ1=`echo ${TMP%s*}`
ssh root@bschang1 -f pkill -9 monitor.sh >/dev/null 2>1
nohup pkill -9 monitor.sh 2>1 &

#Test case 2
#测试用例2，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1小时为粒度，查询每1小时的最大值。
#将运行多个用来模拟interval（1h） 因为cassandra运行中每个worker最后的进程不会工作并马上结束，所以总查询数位query+worker
#SELECT max(usage_user) FROM measurements.cpu WHERE hostname in( 'host_a' ,  'host_b' ,  'host_c' ,  'host_d' ,  'host_e' ,  'host_f' ,  'host_g' ,  'host_h') and time >= 'HOUR_START' and time < 'HOUR_END';
# a,b,c,d,e,f,g,h are random 8 numbers HOUR_END - HOUR_START = 1hour
ssh root@bschang1 -f /root/monitor.sh java >/dev/null 2>1
nohup ./monitor.sh cqlsh >/dev/null 2>1 &
IFQS2=`bin/bulk_query_gen  -seed 123 -format cassandra -query-type 8-host-allbyhr -scale-var $scale -queries $((query + workers)) | bin/query_benchmarker_cassandra  -url=$add  -workers $workers -print-interval 0 -aggregation-plan server| grep wall`
echo -e "${GREEN}cassandra query test case 2 result:${NC}"
echo -e "${GREEN}$IFQS2${NC}"
TMP=`echo $IFQS2|awk '{print($4)}'`
IFQ2=`echo ${TMP%s*}`
ssh root@bschang1 -f pkill -9 monitor.sh >/dev/null 2>1
nohup pkill -9 monitor.sh 2>1 &

#Test case 3
#测试用例3，测试用例3，随机查询12个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以10分钟为粒度，查询每10分钟的最大值
#将运行多个用来模拟interval（10m） 因为cassandra运行中每个worker最后的进程不会工作并马上结束，所以总查询数位query+worker
#SELECT max(usage_user) FROM measurements.cpu WHERE hostname in( 'host_a' ,  'host_b' ,  'host_c' ,  'host_d' ,  'host_e' ,  'host_f' ,  'host_g' ,  'host_h') and time >= 'HOUR_START' and time < 'HOUR_END';
# a,b,c,d,e,f,g,h are random 8 numbers, HOUR_END - HOUR_START = 10m
ssh root@bschang1 -f /root/monitor.sh java >/dev/null 2>1
nohup ./monitor.sh cqlsh >/dev/null 2>1 &
IFQS3=`bin/bulk_query_gen  -seed 123 -format cassandra -query-type 8-host-12-hr -scale-var $scale -queries $((query + workers)) | bin/query_benchmarker_cassandra  -url=$add  -workers $workers -print-interval 0 -aggregation-plan server| grep wall`
echo -e "${GREEN}cassandra query test case 3 result:${NC}"
echo -e "${GREEN}$IFQS3${NC}"
TMP=`echo $IFQS3|awk '{print($4)}'`
IFQ3=`echo ${TMP%s*}`
ssh root@bschang1 -f pkill -9 monitor.sh >/dev/null 2>1
nohup pkill -9 monitor.sh 2>1 &

#Test case 4
#测试用例4，随机查询1个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1分钟为粒度，查询每1分钟的最大值
#将运行多个用来模拟interval（10m） 因为cassandra运行中每个worker最后的进程不会工作并马上结束，所以总查询数位query+worker
#SELECT max(usage_user) FROM measurements.cpu WHERE hostname in( 'host_a' ,  'host_b' ,  'host_c' ,  'host_d' ,  'host_e' ,  'host_f' ,  'host_g' ,  'host_h') and time >= 'HOUR_START' and time < 'HOUR_END';
# a,b,c,d,e,f,g,h are random 8 numbers, HOUR_END - HOUR_START = 1m
ssh root@bschang1 -f /root/monitor.sh java >/dev/null 2>1
nohup ./monitor.sh cqlsh >/dev/null 2>1 &
IFQS4=`bin/bulk_query_gen  -seed 123 -format cassandra -query-type 8-host-1-hr -scale-var $scale -queries $((query + workers)) | bin/query_benchmarker_cassandra  -url=$add  -workers $workers -print-interval 0 -aggregation-plan server| grep wall`
echo -e "${GREEN}cassandra query test case 4 result:${NC}"
echo -e "${GREEN}$IFQS4${NC}"
TMP=`echo $IFQS4|awk '{print($4)}'`
IFQ4=`echo ${TMP%s*}`
ssh root@bschang1 -f pkill -9 monitor.sh >/dev/null 2>1
nohup pkill -9 monitor.sh 2>1 &

ssh root@$add << eeooff
service cassandra stop
exit
eeooff

echo
echo
echo    "======================================================"
echo    "             tsdb performance comparision             "
echo    "======================================================"
echo    "                   Query test cases:                "
echo    " case 1: select the max(value) from all data    "
echo    " filtered out 8 hosts                                 "
echo    "       Query test case 1 takes:                      "
printf  "       cassandra          |       %-4.2f Seconds    \n" $IFQ1 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDQ1
echo    "------------------------------------------------------"
echo    " case 2: select the max(value) from all data          "
echo    " filtered out 8 hosts with an interval of 1 hour     "
echo    " case 2 takes:                                       "
printf  "       cassandra          |       %-4.2f Seconds    \n" $IFQ2 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDQ2
echo    "------------------------------------------------------"
echo    " case 3: select the max(value) from random 12 hours"
echo    " data filtered out 8 hosts with an interval of 10 min         "
echo    " filtered out 8 hosts interval(1h)                   "
echo    " case 3 takes:                                       "
printf  "       cassandra          |       %-4.2f Seconds    \n" $IFQ3 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDQ3
echo    "------------------------------------------------------"
echo    " case 4: select the max(value) from random 1 hour data  "
echo    " data filtered out 8 hosts with an interval of 1 min         "
echo    " case 4 takes:                                        "
printf  "       cassandra          |       %-4.2f Seconds    \n" $IFQ4 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDQ4
printf  "   TDengine no interval   |       %-4.2f Seconds    \n" $TDQ5
echo    "------------------------------------------------------"
echo
docker stop $INFLUX >>/dev/null 2>&1
docker container rm -f $INFLUX >>/dev/null 2>&1
docker stop $TDENGINE >>/dev/null 2>&1
docker container rm -f $TDENGINE >>/dev/null 2>&1
docker network rm tsdbcomp >>/dev/null 2>&1
#bulk_query_gen/bulk_query_gen  -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 5 | query_benchmarker_influxdb/query_benchmarker_influxdb  -urls="http://172.26.89.231:8086" 
#bulk_query_gen/bulk_query_gen  -format tdengine -query-type 1-host-1-hr -scale-var 10 -queries 5 | query_benchmarker_tdengine/query_benchmarker_tdengine  -urls="http://172.26.89.231:6030" 

