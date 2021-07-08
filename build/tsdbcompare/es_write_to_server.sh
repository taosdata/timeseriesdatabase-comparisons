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
add='10.2.0.9'
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
    echo "Prepare data for Elasticsearch...."
    bin/bulk_data_gen -seed 123 -format es-bulk7x -sampling-interval $interval -scale-var $scale -use-case devops -timestamp-start "$st" -timestamp-end "$et" > data/es.dat
fi

ssh -tt root@$add << eeooff
pkill -f elasticsearch
echo 1 > /proc/sys/vm/drop_caches
rm -rf /data/esdata/*
su ubuntu
cd /home/ubuntu/elasticsearch-7.13.2/
bin/elasticsearch -d
sleep 30
exit
exit
eeooff

echo
echo -e "Start test Elasticsearch, result in ${GREEN}Green line${NC}"
pwd
echo "cat data/es.dat |bin/bulk_load_es --batch-size=$batchsize --workers=$workers --urls='http://$add:9200' > res.txt 2>&1"
cat data/es.dat |bin/bulk_load_es --batch-size=$batchsize --workers=$workers --urls="http://$add:9200" > res.txt 2>&1
ELASTICRES=`grep loaded res.txt | head -n 1`
echo
echo -e "${GREEN}scale: ${scale}, batchSize: ${batchsize}, worker: ${workers} ${NC}"
echo -e "${GREEN}Elasticsearch writing result:${NC}"
echo -e "${GREEN}$ELASTICRES${NC}"

DATA=`echo $ELASTICRES|awk '{print $4}'`
TMP=`echo $ELASTICRES|awk '{print($7)}'`
ESWTM=`echo ${TMP%s*}`


echo -e "${GREEN}Stop Elasticsearch${NC}"
ssh -tt root@$add "pkill -f elasticsearch"
echo -e "${GREEN}sleep 30 seconds${NC}"
sleep 30
ssh -tt root@$add "du -sh /data/esdata" > size.txt
TMP=`grep data size.txt`
ESDISK=`echo $TMP|awk '{print $1}'`

echo
echo
echo    "======================================================"
echo    "             tsdb performance comparision             "
printf  "       worker:%-4.2f      |       batch:%-4.2f      \n" $workers $batchsize
echo    "======================================================"
echo -e "       Writing $DATA records test takes:          "
printf  "       Elasticsearch      |       %-4.2f Seconds    \n" $ESWTM
echo    "======================================================"
echo -e "       Writing $DATA records test disk:          "
printf  "       Elasticsearch      |       %-10s             \n" $ESDISK
echo    "------------------------------------------------------"
