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
add='10.2.0.5'
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



echo "---------------  Clean  -----------------"
ssh root@$add << eeooff
rm -rf /data/lib/taos/*
echo 1 > /proc/sys/vm/drop_caches
systemctl start taosd 
sleep 10
exit
eeooff

cat data/tdengine.dat |bin/bulk_load_tdengine --url $add --batch-size $batchsize   -do-load -report-tags n1 -workers $workers -fileout=false -http-api=$interface | grep loaded
