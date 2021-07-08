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
add='10.2.0.9'
interval='10s'
scale=100
st='2018-01-01T00:00:00Z'
et='2018-01-02T00:00:00Z'

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
    echo
    echo "---------------Generating && Inserting Data-----------------"
    echo
    echo "Prepare data for Elastic...."
    bin/bulk_data_gen -seed 123 -format es-bulk7x -sampling-interval $interval -scale-var $scale -use-case devops -timestamp-start "$st" -timestamp-end "$et" > data/es.dat
    cat data/es.dat |bin/bulk_load_es --batch-size=$batchsize --workers=$workers --urls="http://$add:9200"
fi

echo
echo "------------------Querying Data-----------------"
echo

echo 
echo  "start query test, query max from 8 hosts group by 1hour, Elasticsearch"
echo
#Test case 1
#测试用例1，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') ;
# a,b,c,d,e,f,g,h are random 8 numbers.
bin/bulk_query_gen -seed 123 -format es-http -query-type 8-host-all -scale-var $scale -queries 1000 | bin/query_benchmarker_es  -urls="http://$add:9200"  -workers $workers -print-interval 0 > /dev/null 2>case1.txt
ESQS1=`grep wall case1.txt`
echo -e "${GREEN}Elasticsearch query test case 1 result:${NC}"
echo -e "${GREEN}$ESQS1${NC}"
TMP=`echo $ESQS1|awk '{print($6)}'`
ESQ1=`echo ${TMP%s*}`
#Test case 2
#测试用例2，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1小时为粒度，查询每1小时的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') interval(1h);
# a,b,c,d,e,f,g,h are random 8 numbers
bin/bulk_query_gen -seed 123 -format es-http -query-type 8-host-allbyhr -scale-var $scale -queries 1000 | bin/query_benchmarker_es  -urls="http://$add:9200"  -workers $workers -print-interval 0 > /dev/null 2>case2.txt
ESQS2=`grep wall case2.txt`
echo -e "${GREEN}Elasticsearch query test case 2 result:${NC}"
echo -e "${GREEN}$ESQS2${NC}"
TMP=`echo $ESQS2|awk '{print($6)}'`
ESQ2=`echo ${TMP%s*}`
#Test case 3
#测试用例3，测试用例3，随机查询12个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以10分钟为粒度，查询每10分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =12 hour
#ELASTICSEARCHQUERY=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_es  -urls="http://127.0.0.1:9200"  -workers $workers -print-interval 0|grep wall`
bin/bulk_query_gen -seed 123 -format es-http -query-type 8-host-12-hr -scale-var $scale -queries 1000 | bin/query_benchmarker_es  -urls="http://$add:9200"  -workers $workers -print-interval 0 > /dev/null 2>case3.txt
ESQS3=`grep wall case3.txt`
echo -e "${GREEN}Elasticsearch query test case 3 result:${NC}"
echo -e "${GREEN}$ESQS3${NC}"
TMP=`echo $ESQS3|awk '{print($6)}'`
ESQ3=`echo ${TMP%s*}`
#Test case 4
#测试用例4，随机查询1个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1分钟为粒度，查询每1分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =1 hours
#ELASTICSEARCHQUERY=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_es  -urls="http://127.0.0.1:9200"  -workers $workers -print-interval 0|grep wall`
bin/bulk_query_gen -seed 123 -format es-http -query-type 8-host-1-hr -scale-var $scale -queries 1000 | bin/query_benchmarker_es  -urls="http://$add:9200"  -workers $workers -print-interval 0 > /dev/null 2>case4.txt
ESQS4=`grep wall case4.txt`
echo -e "${GREEN}Elasticsearch query test case 4 result:${NC}"
echo -e "${GREEN}$ESQS4${NC}"
TMP=`echo $ESQS4|awk '{print($6)}'`
ESQ4=`echo ${TMP%s*}`


echo
echo
echo    "======================================================"
echo    "             tsdb performance comparision             "
echo    "======================================================"
echo    "                   Query test cases:                "
echo    " case 1: select the max(value) from all data    "
echo    " filtered out 8 hosts                                 "
echo    "       Query test case 1 takes:                      "
printf  "       Elasticsearch      |       %-4.2f Seconds    \n" $ESQ1
echo    "------------------------------------------------------"
echo    " case 2: select the max(value) from all data          "
echo    " filtered out 8 hosts with an interval of 1 hour     "
echo    " case 2 takes:                                       "
printf  "       Elasticsearch      |       %-4.2f Seconds    \n" $ESQ2
echo    "------------------------------------------------------"
echo    " case 3: select the max(value) from random 12 hours"
echo    " data filtered out 8 hosts with an interval of 10 min         "
echo    " filtered out 8 hosts interval(1h)                   "
echo    " case 3 takes:                                       "
printf  "       Elasticsearch      |       %-4.2f Seconds    \n" $ESQ3
echo    "------------------------------------------------------"
echo    " case 4: select the max(value) from random 1 hour data  "
echo    " data filtered out 8 hosts with an interval of 1 min         "
echo    " case 4 takes:                                        "
printf  "       Elasticsearch      |       %-4.2f Seconds    \n" $ESQ4
echo    "------------------------------------------------------"
echo

#bulk_query_gen/bulk_query_gen  -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_influxdb/query_benchmarker_influxdb  -urls="http://172.26.89.231:8086" 
#bulk_query_gen/bulk_query_gen  -format tdengine -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_tdengine/query_benchmarker_tdengine  -urls="http://172.26.89.231:6030" 
