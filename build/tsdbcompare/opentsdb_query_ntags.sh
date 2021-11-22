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

# 
# 
# q512=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-512" -query-type="8-host-all" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 512 tags result:${NC}"
# echo -e "${GREEN}$q512${NC}"
# TMP=`echo $q512|awk '{print($4)}'`
# o512=`echo ${TMP%s*}`
# 
# q256=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-256" -query-type="8-host-all" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 256 tags result:${NC}"
# echo -e "${GREEN}$q256${NC}"
# TMP=`echo $q256|awk '{print($4)}'`
# qo256=`echo ${TMP%s*}`
# 
# q128=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-128" -query-type="8-host-all" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 128 tags result:${NC}"
# echo -e "${GREEN}$q128${NC}"
# TMP=`echo $q128|awk '{print($4)}'`
# qo128=`echo ${TMP%s*}`
# 
# q64=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-64" -query-type="8-host-all" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 64 tags result:${NC}"
# echo -e "${GREEN}$q64${NC}"
# TMP=`echo $q64|awk '{print($4)}'`
# qo64=`echo ${TMP%s*}`
# 
# 
# q32=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-32" -query-type="8-host-all" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 32 tags result:${NC}"
# echo -e "${GREEN}$q32${NC}"
# TMP=`echo $q32|awk '{print($4)}'`
# qo32=`echo ${TMP%s*}`
# 
# q16=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-16" -query-type="8-host-all" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 16 tags result:${NC}"
# echo -e "${GREEN}$q16${NC}"
# TMP=`echo $q16|awk '{print($4)}'`
# qo16=`echo ${TMP%s*}`
# 
# q8=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb" -query-type="8-host-all" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 8 tags result:${NC}"
# echo -e "${GREEN}$q8${NC}"
# TMP=`echo $q8|awk '{print($4)}'`
# qo8=`echo ${TMP%s*}`
# 
# q1=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-1" -query-type="8-host-all" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 1 tags result:${NC}"
# echo -e "${GREEN}$q1${NC}"
# TMP=`echo $q1|awk '{print($4)}'`
# qo1=`echo ${TMP%s*}`
# 
# 
# echo
# echo
# echo    "======================================================"
# echo    "             tsdb performance comparision             "
# echo    "======================================================"
# echo    "          Query test cases Nhost-all(1):                "
# printf  "       OpenTSDB512         |       %-4.2f Seconds    \n" $qo512 
# printf  "       OpenTSDB256         |       %-4.2f Seconds    \n" $qo256
# printf  "       OpenTSDB128         |       %-4.2f Seconds    \n" $qo128
# printf  "       OpenTSDB64          |       %-4.2f Seconds    \n" $qo64
# printf  "       OpenTSDB32          |       %-4.2f Seconds    \n" $qo32 
# printf  "       OpenTSDB16          |       %-4.2f Seconds    \n" $qo16 
# printf  "       OpenTSDB8           |       %-4.2f Seconds    \n" $qo8 
# printf  "       OpenTSDB1           |       %-4.2f Seconds    \n" $qo1 
# echo    "------------------------------------------------------"
# echo
# 
# 
# q512=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-512" -query-type="8-host-allbyhr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 512 tags result:${NC}"
# echo -e "${GREEN}$q512${NC}"
# TMP=`echo $q512|awk '{print($4)}'`
# o512=`echo ${TMP%s*}`
# 
# q256=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-256" -query-type="8-host-allbyhr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 256 tags result:${NC}"
# echo -e "${GREEN}$q256${NC}"
# TMP=`echo $q256|awk '{print($4)}'`
# qo256=`echo ${TMP%s*}`
# 
# q128=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-128" -query-type="8-host-allbyhr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 128 tags result:${NC}"
# echo -e "${GREEN}$q128${NC}"
# TMP=`echo $q128|awk '{print($4)}'`
# qo128=`echo ${TMP%s*}`
# 
# q64=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-64" -query-type="8-host-allbyhr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 64 tags result:${NC}"
# echo -e "${GREEN}$q64${NC}"
# TMP=`echo $q64|awk '{print($4)}'`
# qo64=`echo ${TMP%s*}`
# 
# 
# q32=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-32" -query-type="8-host-allbyhr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 32 tags result:${NC}"
# echo -e "${GREEN}$q32${NC}"
# TMP=`echo $q32|awk '{print($4)}'`
# qo32=`echo ${TMP%s*}`
# 
# q16=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-16" -query-type="8-host-allbyhr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 16 tags result:${NC}"
# echo -e "${GREEN}$q16${NC}"
# TMP=`echo $q16|awk '{print($4)}'`
# qo16=`echo ${TMP%s*}`
# 
# q8=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb" -query-type="8-host-allbyhr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 8 tags result:${NC}"
# echo -e "${GREEN}$q8${NC}"
# TMP=`echo $q8|awk '{print($4)}'`
# qo8=`echo ${TMP%s*}`
# 
# q1=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-1" -query-type="8-host-allbyhr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 1 tags result:${NC}"
# echo -e "${GREEN}$q1${NC}"
# TMP=`echo $q1|awk '{print($4)}'`
# qo1=`echo ${TMP%s*}`
# 
# 
# echo
# echo
# echo    "======================================================"
# echo    "             tsdb performance comparision             "
# echo    "======================================================"
# echo    "          Query test cases Nhost-by-1Hr(2):                "
# printf  "       OpenTSDB512         |       %-4.2f Seconds    \n" $qo512 
# printf  "       OpenTSDB256         |       %-4.2f Seconds    \n" $qo256
# printf  "       OpenTSDB128         |       %-4.2f Seconds    \n" $qo128
# printf  "       OpenTSDB64          |       %-4.2f Seconds    \n" $qo64
# printf  "       OpenTSDB32          |       %-4.2f Seconds    \n" $qo32 
# printf  "       OpenTSDB16          |       %-4.2f Seconds    \n" $qo16 
# printf  "       OpenTSDB8           |       %-4.2f Seconds    \n" $qo8 
# printf  "       OpenTSDB1           |       %-4.2f Seconds    \n" $qo1 
# echo    "------------------------------------------------------"
# echo
# 
q512=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-512" -query-type="8-host-12-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 512 tags result:${NC}"
echo -e "${GREEN}$q512${NC}"
TMP=`echo $q512|awk '{print($4)}'`
o512=`echo ${TMP%s*}`

q256=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-256" -query-type="8-host-12-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 256 tags result:${NC}"
echo -e "${GREEN}$q256${NC}"
TMP=`echo $q256|awk '{print($4)}'`
qo256=`echo ${TMP%s*}`

q128=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-128" -query-type="8-host-12-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 128 tags result:${NC}"
echo -e "${GREEN}$q128${NC}"
TMP=`echo $q128|awk '{print($4)}'`
qo128=`echo ${TMP%s*}`

q64=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-64" -query-type="8-host-12-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 64 tags result:${NC}"
echo -e "${GREEN}$q64${NC}"
TMP=`echo $q64|awk '{print($4)}'`
qo64=`echo ${TMP%s*}`


q32=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-32" -query-type="8-host-12-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 32 tags result:${NC}"
echo -e "${GREEN}$q32${NC}"
TMP=`echo $q32|awk '{print($4)}'`
qo32=`echo ${TMP%s*}`

q16=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-16" -query-type="8-host-12-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 16 tags result:${NC}"
echo -e "${GREEN}$q16${NC}"
TMP=`echo $q16|awk '{print($4)}'`
qo16=`echo ${TMP%s*}`

q8=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb" -query-type="8-host-12-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 8 tags result:${NC}"
echo -e "${GREEN}$q8${NC}"
TMP=`echo $q8|awk '{print($4)}'`
qo8=`echo ${TMP%s*}`

q1=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-1" -query-type="8-host-12-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
echo -e "${GREEN}Opentsdb query test 1 tags result:${NC}"
echo -e "${GREEN}$q1${NC}"
TMP=`echo $q1|awk '{print($4)}'`
qo1=`echo ${TMP%s*}`


echo
echo
echo    "======================================================"
echo    "             tsdb performance comparision             "
echo    "======================================================"
echo    "          Query test cases Nhost-by-10min(3):                "
printf  "       OpenTSDB512         |       %-4.2f Seconds    \n" $qo512 
printf  "       OpenTSDB256         |       %-4.2f Seconds    \n" $qo256
printf  "       OpenTSDB128         |       %-4.2f Seconds    \n" $qo128
printf  "       OpenTSDB64          |       %-4.2f Seconds    \n" $qo64
printf  "       OpenTSDB32          |       %-4.2f Seconds    \n" $qo32 
printf  "       OpenTSDB16          |       %-4.2f Seconds    \n" $qo16 
printf  "       OpenTSDB8           |       %-4.2f Seconds    \n" $qo8 
printf  "       OpenTSDB1           |       %-4.2f Seconds    \n" $qo1 
echo    "------------------------------------------------------"
echo

# q512=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-512" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 512 tags result:${NC}"
# echo -e "${GREEN}$q512${NC}"
# TMP=`echo $q512|awk '{print($4)}'`
# o512=`echo ${TMP%s*}`
# 
# q256=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-256" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 256 tags result:${NC}"
# echo -e "${GREEN}$q256${NC}"
# TMP=`echo $q256|awk '{print($4)}'`
# qo256=`echo ${TMP%s*}`
# 
# q128=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-128" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 128 tags result:${NC}"
# echo -e "${GREEN}$q128${NC}"
# TMP=`echo $q128|awk '{print($4)}'`
# qo128=`echo ${TMP%s*}`
# 
# q64=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-64" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 64 tags result:${NC}"
# echo -e "${GREEN}$q64${NC}"
# TMP=`echo $q64|awk '{print($4)}'`
# qo64=`echo ${TMP%s*}`
# 
# q32=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-32" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 32 tags result:${NC}"
# echo -e "${GREEN}$q32${NC}"
# TMP=`echo $q32|awk '{print($4)}'`
# qo32=`echo ${TMP%s*}`
# 
# q16=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-16" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 16 tags result:${NC}"
# echo -e "${GREEN}$q16${NC}"
# TMP=`echo $q16|awk '{print($4)}'`
# qo16=`echo ${TMP%s*}`
# 
# q8=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 8 tags result:${NC}"
# echo -e "${GREEN}$q8${NC}"
# TMP=`echo $q8|awk '{print($4)}'`
# qo8=`echo ${TMP%s*}`
# 
# q1=`bin/bulk_query_gen -timestamp-end="2018-01-02T08:00:01Z" -seed=123 -format="opentsdb-1" -query-type="8-host-1-hr" -timestamp-start="2017-12-31T23:59:59Z"  -scale-var=1000 -queries=1000 |bin/query_benchmarker_opentsdb -urls="http://192.168.1.84:4242" -workers 16 -print-interval 0 |grep wall`
# echo -e "${GREEN}Opentsdb query test 1 tags result:${NC}"
# echo -e "${GREEN}$q1${NC}"
# TMP=`echo $q1|awk '{print($4)}'`
# qo1=`echo ${TMP%s*}`
# 
# 
# echo
# echo
# echo    "======================================================"
# echo    "             tsdb performance comparision             "
# echo    "======================================================"
# echo    "          Query test cases Nhost-by-1min(4):                "
# printf  "       OpenTSDB512         |       %-4.2f Seconds    \n" $qo512 
# printf  "       OpenTSDB256         |       %-4.2f Seconds    \n" $qo256
# printf  "       OpenTSDB128         |       %-4.2f Seconds    \n" $qo128
# printf  "       OpenTSDB64          |       %-4.2f Seconds    \n" $qo64
# printf  "       OpenTSDB32          |       %-4.2f Seconds    \n" $qo32 
# printf  "       OpenTSDB16          |       %-4.2f Seconds    \n" $qo16 
# printf  "       OpenTSDB8           |       %-4.2f Seconds    \n" $qo8 
# printf  "       OpenTSDB1           |       %-4.2f Seconds    \n" $qo1 
# echo    "------------------------------------------------------"
# echo
