#!/bin/bash

batchsize=5000
workers=16
interface='false'
gene=1
add='test217'
interval='10s'
scale=100
st='2018-01-01T00:00:00Z'
et='2018-01-02T00:00:00Z'
TDPath="/var/lib/taos/"
InfPath="/var/lib/influxdb/"
while getopts "b:w:n:g:a:i:s:t:e:T:I:" opt
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
        T)
        echo "TDengine rootPath:$OPTARG"
        TDPath=$OPTARG
        ;;
        I)
        echo "Influxdb rootPath:$OPTARG"
        InfPath=$OPTARG
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
        echo    "------------------------------------------------------"
        echo    "T | TDengine rootPath (default:/var/lib/taos)"
        echo    "------------------------------------------------------"
        echo    "I | Influxdb rootPath (default:/var/lib/taos)"
        echo    "======================================================"
        exit 1;;
    esac
done

for s in {100,200,400,600,800,1000};do
    for i in {5000,2500,1000,1};do
		./write_to_server.sh -b $i -w 100 -g 1 -s $s -a "test217" -T "/var/lib/taos/"  -t '2018-01-01T00:00:00Z' -e '2018-01-01T00:10:00Z'
        for j in {50,16,1};do
            ./write_to_server.sh -b $i -w $j -g 0 -s $s -a "test217" -I "/var/lib/influxdb/"  -t '2018-01-01T00:00:00Z' -e '2018-01-01T00:10:00Z'
        done
    done
done