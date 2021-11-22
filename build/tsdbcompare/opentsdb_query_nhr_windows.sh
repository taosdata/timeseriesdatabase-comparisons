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
add='127.0.0.1'
interval='10s'
scale=100
st='2018-01-01T00:00:00Z'
et='2018-01-02T00:00:01Z'

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
        echo    "a | address of TD & opentsdb"
        echo    "------------------------------------------------------"
        echo    "s | scale-var(default:100)"
        echo    "------------------------------------------------------"
        echo    "t | timestamp-start(default:$st)"
        echo    "------------------------------------------------------"
        echo    "e | timestamp-end(default:$et)"
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

    echo 
    echo "Prepare data for TDengine...."
    bin/bulk_data_gen -seed 123 -format tdengine -sampling-interval $interval -tdschema-file config/TDengineSchema.toml -scale-var $scale -use-case devops -timestamp-start "$st" -timestamp-end "$et"  > data/tdengine.dat
    cat data/tdengine.dat |bin/bulk_load_tdengine --url $add --batch-size 5000  -do-load -report-tags n1 -workers $workers -fileout=false -http-api=$interface
fi

echo
echo "------------------Querying Data-----------------"
echo


echo 
echo  "start query test, query max from 8 hosts group by 1 hour, TDengine"
echo


q1h=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 1h tags result:${NC}"
echo -e "${GREEN}$q1h${NC}"
TMP=`echo $q1h|awk '{print($4)}'`
qo1h=`echo ${TMP%s*}`

q2h=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-2h" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 2h tags result:${NC}"
echo -e "${GREEN}$q2h${NC}"
TMP=`echo $q2h|awk '{print($4)}'`
qo2h=`echo ${TMP%s*}`

q4h=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-4h" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 4h tags result:${NC}"
echo -e "${GREEN}$q4h${NC}"
TMP=`echo $q4h|awk '{print($4)}'`
qo4h=`echo ${TMP%s*}`

q8h=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-8h" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 8h tags result:${NC}"
echo -e "${GREEN}$q8h${NC}"
TMP=`echo $q8h|awk '{print($4)}'`
qo8h=`echo ${TMP%s*}`


q12h=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-12h" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 12h tags result:${NC}"
echo -e "${GREEN}$q12h${NC}"
TMP=`echo $q12h|awk '{print($4)}'`
qo12h=`echo ${TMP%s*}`

qmixfunc=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-mixfunc" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 16 tags result:${NC}"
echo -e "${GREEN}$qmixfunc${NC}"
TMP=`echo $qmixfunc|awk '{print($4)}'`
qomixfunc=`echo ${TMP%s*}`


echo
echo
echo    "======================================================"
echo    "             tsdb performance comparision             "
echo    "======================================================"
echo    "                   Query test cases:                "
printf  "       OpenTSDB1h         |       %-4.2f Seconds    \n" $qo1h 
printf  "       OpenTSDB2h         |       %-4.2f Seconds    \n" $qo2h
printf  "       OpenTSDB4h         |       %-4.2f Seconds    \n" $qo4h
printf  "       OpenTSDB8h         |       %-4.2f Seconds    \n" $qo8h
printf  "       OpenTSDB12h        |       %-4.2f Seconds    \n" $qo12h 
printf  "      OpenTSDBmixfunc     |       %-4.2f Seconds    \n" $qomixfunc 
echo    "------------------------------------------------------"
echo
# docker stop $INFLUX >>/dev/null 2>&1
# docker container rm -f $INFLUX >>/dev/null 2>&1
# docker stop $TDENGINE >>/dev/null 2>&1
# docker container rm -f $TDENGINE >>/dev/null 2>&1
# docker network rm tsdbcomp >>/dev/null 2>&1
# bulk_query_gen/bulk_query_gen  -format opentsdb-http -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_opentsdbdb/query_benchmarker_opentsdbdb  -urls="http://172.26.89.231:8086" 
#bulk_query_gen/bulk_query_gen  -format tdengine -query-type 1-host-1-hr -scale-var 10 -queries 1000 | query_benchmarker_tdengine/query_benchmarker_tdengine  -urls="http://172.26.89.231:6030" 
