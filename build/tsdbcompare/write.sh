#!/bin/bash


# Color setting
RED='\033[0;31m'
GREEN='\033[1;32m'
GREEN_DARK='\033[0;32m'
GREEN_UNDERLINE='\033[4;32m'
NC='\033[0m'

batchsize=5000
workers=16
interface='false'
gene=1
add='127.0.0.1'
interval='10s'
scale=100
st='2018-01-01T00:00:00Z'
et='2018-01-02T00:00:00Z'
while getopts "b:w:n:g:a:i:s:t:e:" opt
do
    case $opt in
        b)
        echo "batchsize:$OPTARG"
        batchsize=$OPTARG
        ;;
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
        echo    "b | batchsize(1~5000)"
        echo    "------------------------------------------------------"
        echo    "w | workers"
        echo    "------------------------------------------------------"
        echo    "n | TD's interface(false:cgo,true:rest)"
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
        echo    "g | genate data(0:no,1:yes)"
        echo    "======================================================"
        exit 1;;
    esac
done
if [ $# != 18 ];then
    echo "variables not all defined,use value as follow :"
    echo "generate data: scale-var: $scale ,interval: $interval ,timestamp-start: $st ,timestamp-stop: $et"
    echo "batchsize:$batchsize ,workers:$workers ,TD's interface: $interface ,generate data: $gene , address: $add"
fi
echo "$a"

if [[ $gene == 1 ]];then 
    if [ ! -d "data" ]; then
        mkdir data
    fi
    echo
    echo "---------------Generating Data-----------------"
    echo
    echo "Prepare data for InfluxDB...."
    bin/bulk_data_gen -seed 123 -format influx-bulk -sampling-interval $interval -scale-var $scale -use-case devops -timestamp-start "$st" -timestamp-end "$et" >data/influx.dat

    echo 
    echo "Prepare data for TDengine...."
    bin/bulk_data_gen -seed 123 -format tdengine -sampling-interval $interval -tdschema-file config/TDengineSchema.toml -scale-var $scale -use-case devops -timestamp-start "$st" -timestamp-end "$et"  > data/tdengine.dat
fi
echo
echo "---------------  Clean  -----------------"
rm -rf /var/lib/taos/*
rm -rf /var/lib/influxdb/*
echo 1 > /proc/sys/vm/drop_caches
systemctl start taosd 
sleep 10

echo
echo "------------------Writing Data-----------------"

echo
echo -e "Start test TDengine, result in ${GREEN}Green line${NC}"
# taos -s 'drop database devops;'
TDENGINERES=`cat data/tdengine.dat |bin/bulk_load_tdengine --url $add --batch-size $batchsize   -do-load -report-tags n1 -workers $workers -fileout=false -http-api=$interface | grep loaded`
echo $TDENGINERES          
echo
echo -e "${GREEN}TDengine writing result:${NC}"
echo -e "${GREEN}$TDENGINERES${NC}"
DATA=`echo $TDENGINERES|awk '{print($2)}'`
TMP=`echo $TDENGINERES|awk '{print($5)}'`
TDWTM=`echo ${TMP%s*}`
systemctl stop taosd 
echo 1 > /proc/sys/vm/drop_caches
systemctl start influxd
sleep 10
echo
echo -e "Start test InfluxDB, result in ${GREEN}Green line${NC}"
# curl 'http://localhost:8086/query?q=drop%20database%20benchmark' -X POST
INFLUXRES=`cat data/influx.dat  |bin/bulk_load_influx --batch-size=$batchsize --workers=$workers --urls="http://$add:8086" | grep loaded`
echo
echo -e "${GREEN}InfluxDB writing result:${NC}"
echo -e "${GREEN}$INFLUXRES${NC}"

TMP=`echo $INFLUXRES|awk '{print($5)}'`
IFWTM=`echo ${TMP%s*}`
systemctl stop influxd
TDDISK=`du -sh /var/lib/taos/vnode | awk '{print $1}' `
IFDISK=`du -sh /var/lib/influxdb/data | awk '{print $1}' `

echo
echo
echo    "======================================================"
echo    "             tsdb performance comparision             "
printf  "       worker:%-4.2f      |       batch:%-4.2f      \n" $workers $batchsize
echo    "======================================================"
echo -e "       Writing $DATA records test takes:          "
printf  "       InfluxDB           |       %-4.2f Seconds    \n" $IFWTM 
printf  "       TDengine           |       %-4.2f Seconds    \n" $TDWTM
echo    "======================================================"
echo -e "       Writing $DATA records test disk:          "
printf  "       InfluxDB           |       %-10s     \n" $IFDISK
printf  "       TDengine           |       %-10s     \n" $TDDISK
echo    "------------------------------------------------------"



