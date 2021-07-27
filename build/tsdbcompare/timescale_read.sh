#!/bin/bash

helpdoc() {
cat <<EOF
简述：

        本shell脚本用于进行timescaleDB和TDengine的插入性能对比测试

用法:

        $0	-w <workers>	-a <address>	-g <genate data> \
        -t <timestamp-start> -e <timestamp-end> \
        -c <use-case>	-s <scale> -q <query_type> -n <queries>

参数:

        ===========================================================================================
        -w      workers,		uint,   设置查询工作客户端数量,默认16;
        -------------------------------------------------------------------------------------------
        -a      address,		string, 设置TDengine或timescaleDB服务器地址,默认10.2.0.5;
        -------------------------------------------------------------------------------------------
        -t      timestamp-start,string, 设置开始时间戳,默认2018-01-01T00:00:00Z';
        -------------------------------------------------------------------------------------------
       	-e      timestamp-end,	string, 设置结束时间戳,默认2018-01-02T00:00:00Z';
        ---------------------------------------   ----------------------------------------------------
        -g      gene,			uint,   设置是否生成数据(0:no,1:yes),默认1;
        -------------------------------------------------------------------------------------------
        -c      use-case,		string, 设置case类型(devops/iot/cpu-only),默认devops;
        -------------------------------------------------------------------------------------------
        -s      scale,			uint,   设置模拟的设备数量，默认100.
        -------------------------------------------------------------------------------------------
        -q      query_type,   string,   设置生成的查询类型(8-host-all/8-host-allbyhr/8-host-12-hr/ 8-host-1-hr)，默认8-host-all
        -------------------------------------------------------------------------------------------
        -n      queries,    uint,   设置生成的查询数量，默认100
        ===========================================================================================

EOF
}


# Color setting
RED='\033[0;31m'
GREEN='\033[1;32m'
GREEN_DARK='\033[0;32m'
GREEN_UNDERLINE='\033[4;32m'
NC='\033[0m'

workers=16
dbname="test"
interface='cgo'
gene=0
add=10.2.0.5
interval='10s'
scale=100
queries=1000
query_type="cpu-max-all-8"
st='2018-01-01T00:00:00Z'
et='2018-01-02T00:00:00Z'
usecase="devops"

echo -e "${GREEN}====================This test option: scale: ${scale}, query_type: ${query_type}, worker: ${workers} =====================${NC}"
echo -e "${GREEN}------------------------------------- timescaledb query test start -------------------------------------${NC}"

while getopts "hw:n:g:a:s:t:e:c:q:d:" opt
do
    case $opt in
      h)
        helpdoc
        exit 0
        ;;
      w)
        echo "workers:$OPTARG"
        workers=$OPTARG
        ;;
      n)
        echo "queries:$OPTARG"
        queries=$OPTARG
        ;;
      g)
        echo "whether generate data:$OPTARG"
        gene=$OPTARG
        ;;
      q)
        echo "query_type :$OPTARG"
        query_type=$OPTARG
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
      c)
        echo "use-case: $OPTARG"
        usecase=$OPTARG
        ;;
      d)
        echo "database name: $OPTARG"
        dbname=$OPTARG
        ;;
      ?)
        helpdoc
        exit 1
        ;;
    esac
done

echo "variables :"
echo "generate data: "
echo "  scale-var: $scale , query_type: $query_type , timestamp-start: $st , timestamp-stop: $et"
echo "  workers: $workers , use-case: $usecase , generate data: $gene , address: $add"


#################### 0 prepare data source ###############
datafile="data/timescaledb.gz"

if [[ $gene == 1 ]];then 
    if [ ! -d "data" ]; then
        mkdir data
    fi
    echo
    echo "---------------Generating && Inserting Data-----------------"
    echo
    echo "Prepare data for timescaledb..."
    ../../bin/timescale_generate_data   --seed=123 --format="timescaledb" --log-interval=$interval \
      --scale=$scale --use-case=$usecase \
      --timestamp-start=$st --timestamp-end=$et  | gzip  > ${datafile}

    cat ${datafile} | gunzip | ../../bin/timescale_load --workers=50 --db-name=$dbname --batch-size=5000 --host=$add  > tsres.log  2>&1
fi

echo
echo "------------------Querying Data-----------------"
echo

echo 
echo  "start query test, query max from 8 hosts group by 1hour, timescaledb"
echo
#Test case 1
#测试用例1，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') ;
# a,b,c,d,e,f,g,h are random 8 numbers.
../../bin/timescale_generate_queries --seed 123 --format="timescaledb" --query-type="cpu-max-all-8"  \
  --queries=1000 --use-case=$usecase --db-name=$dbname  --scale=100 | \
  gzip > data/query_cpu-max-all-8.gz
cat data/query_cpu-max-all-8.gz | gunzip | ../../bin/timescale_run_queries  --workers=$workers  --db-name=$dbname  --hosts=$add  > /dev/null 2>case1.log
TSQS1=`awk '/all queries/{getline; print}' case1.log`
echo -e "${GREEN}timescaledb query test case 1 result:${NC}"
echo -e "${GREEN}$ESQS1${NC}"
TMP=`echo $ESQS1|awk '{print $6}'`
TSQ1=`echo ${TMP%s*}`
#Test case 2
#测试用例2，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1小时为粒度，查询每1小时的最大值。
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') interval(1h);
# a,b,c,d,e,f,g,h are random 8 numbers
../../bin/timescale_generate_queries --seed 123 --format="timescaledb" --query-type="cpu-max-all-8-by-1hr"  \
  --queries=1000 --use-case=$usecase --db-name=$dbname  --scale=100 | \
  gzip > data/query_cpu-max-all-8.gz
cat data/query_cpu-max-all-8.gz | gunzip | ../../bin/timescale_run_queries  --workers=$workers  --db-name=$dbname  --hosts=$add  > /dev/null 2>case1.log
TSQS1=`awk '/all queries/{getline; print}' case1.log`
echo -e "${GREEN}timescaledb query test case 1 result:${NC}"
echo -e "${GREEN}$ESQS1${NC}"
TMP=`echo $ESQS1|awk '{print $6}'`
TSQ1=`echo ${TMP%s*}`
#Test case 3
#测试用例3，测试用例3，随机查询12个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以10分钟为粒度，查询每10分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =12 hour
#ELASTICSEARCHQUERY=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_es  -urls="http://127.0.0.1:9200"  -workers $workers -print-interval 0|grep wall`
../../bin/timescale_generate_queries --seed 123 --format="timescaledb" --query-type="cpu-max-all-8-10-12hr"  \
  --queries=1000 --use-case=$usecase --db-name=$dbname  --scale=100 | \
  gzip > data/query_cpu-max-all-8.gz
cat data/query_cpu-max-all-8.gz | gunzip | ../../bin/timescale_run_queries  --workers=$workers  --db-name=$dbname  --hosts=$add  > /dev/null 2>case1.log
TSQS1=`awk '/all queries/{getline; print}' case1.log`
echo -e "${GREEN}timescaledb query test case 1 result:${NC}"
echo -e "${GREEN}$ESQS1${NC}"
TMP=`echo $ESQS1|awk '{print $6}'`
TSQ1=`echo ${TMP%s*}`
#Test case 4
#测试用例4，随机查询1个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1分钟为粒度，查询每1分钟的最大值
#select max(usage_user) from cpu where(hostname='host_a' and hostname='host_b'and hostname='host_c'and hostname='host_d'and hostname='host_e'and hostname='host_f' and hostname='host_g'and hostname='host_h') and time >x and time <y interval(10m);
# a,b,c,d,e,f,g,h are random 8 numbers, y-x =1 hours
#ELASTICSEARCHQUERY=`bin/bulk_query_gen  -seed 123 -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_es  -urls="http://127.0.0.1:9200"  -workers $workers -print-interval 0|grep wall`
../../bin/timescale_generate_queries --seed 123 --format="timescaledb" --query-type="cpu-max-all-8-1-1hr"  \
  --queries=1000 --use-case=$usecase --db-name=$dbname  --scale=100 | \
  gzip > data/query_cpu-max-all-8.gz
cat data/query_cpu-max-all-8.gz | gunzip | ../../bin/timescale_run_queries  --workers=$workers  --db-name=$dbname  --hosts=$add  > /dev/null 2>case1.log
TSQS1=`awk '/all queries/{getline; print}' case1.log`
echo -e "${GREEN}timescaledb query test case 1 result:${NC}"
echo -e "${GREEN}$ESQS1${NC}"
TMP=`echo $ESQS1|awk '{print $6}'`
TSQ1=`echo ${TMP%s*}`


echo
echo
echo    "======================================================"
echo    "             tsdb performance comparision             "
echo    "======================================================"
echo    "             Query test cases with wokers:$workers                "
echo    " case 1: select the max(value) from all data    "
echo    " filtered out 8 hosts                                 "
echo    "       Query test case 1 takes:                      "
printf  "       timescaledb      |       %-4.2f Seconds    \n" $TSQ1
echo    "------------------------------------------------------"
echo    " case 2: select the max(value) from all data          "
echo    " filtered out 8 hosts with an interval of 1 hour     "
echo    " case 2 takes:                                       "
printf  "       timescaledb      |       %-4.2f Seconds    \n" $TSQ2
echo    "------------------------------------------------------"
echo    " case 3: select the max(value) from random 12 hours"
echo    " data filtered out 8 hosts with an interval of 10 min         "
echo    " filtered out 8 hosts interval(1h)                   "
echo    " case 3 takes:                                       "
printf  "       timescaledb      |       %-4.2f Seconds    \n" $TSQ3
echo    "------------------------------------------------------"
echo    " case 4: select the max(value) from random 1 hour data  "
echo    " data filtered out 8 hosts with an interval of 1 min         "
echo    " case 4 takes:                                        "
printf  "       timescaledb      |       %-4.2f Seconds    \n" $TSQ4
echo    "------------------------------------------------------"
echo

#bulk_query_gen/bulk_query_gen  -format influx-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_influxdb/query_benchmarker_influxdb  -urls="http://172.26.89.231:8086" 
#bulk_query_gen/bulk_query_gen  -format tdengine -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_tdengine/query_benchmarker_tdengine  -urls="http://172.26.89.231:6030" 
