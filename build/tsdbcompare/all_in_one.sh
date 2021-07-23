#!/bin/bash

add="serv"
compareType="insert"
queryType="cgo"
fsync='false'

while getopts "a:i:q:p:f:" opt; do
    case $opt in
    a)
        add=$OPTARG
        ;;
    i)
        compareType=$OPTARG
        ;;
    q)
        queryType=$OPTARG
        ;;
    ?)
        echo "======================================================"
        echo "a | server address(default 'serv')"
        echo "------------------------------------------------------"
        echo "i | the type of comparsion you want to run"
        echo "  | (insert: run insert test | query: run query test)"
        echo "  | (all: run all insert and query test)"
        echo "  | (11: run insert test where worker = 1 batch = 1)"
        echo "  | (default: insert)"
        echo "------------------------------------------------------"
        echo "q | query method run by tdengine"
        echo "  | (cgo: run in cgo mode | fast: run in restful mode)"
        echo "  | (default: cgo)"
        echo "======================================================"
        exit 1
        ;;
    esac
done
echo "address: $add, compareType: $compareType, queryType: $queryType"

if [[ $fsync == "false" ]]; then
        ssh root@$add <<eeooff
        echo `sed -i "s/.*trickle_fsync.*/trickle_fsync: true/g" /etc/cassandra/cassandra.yaml`
        sed -i "s/.*walLevel.*/walLevel 2 /g" /etc/taos/taos.cfg
eeooff
elif [[ $fsync == "true" ]]; then
        ssh root@$add <<eeooff
        echo `sed -i "s/.*trickle_fsync.*/trickle_fsync: false/g" /etc/cassandra/cassandra.yaml`
        sed -i "s/.*walLevel.*/walLevel 2 /g" /etc/taos/taos.cfg
eeooff
fi

if [[ $compareType == "insert" || $compareType == "both" ]]; then

    echo "start insert test between TDengine and Cassandra"
    echo "Worker = 1, Batch = 1 is not included"
    for s in {800,1000}; do
        for i in {2000,1000,500,1}; do
            ./write_to_server.sh -b $i -w 100 -g 1 -s $s -a $add
            for j in {50,16}; do
                ./write_to_server.sh -b $i -w $j -g 0 -s $s -a $add
            done
        done
        for i in {2000,1000,500}; do
            ./write_to_server.sh -b $i -w 1 -g 0 -s $s -a $add
        done
    done
fi

if [ $compareType == "11" || $compareType == "all" ]; then
    echo "start insert test between TDengine and Cassandra"
    echo "only contain Worker = 1, Batch = 1"
    ./write_to_server.sh -b 1 -w 1 -g 1 -s 100 -a $add
    ./write_to_server.sh -b 1 -w 1 -g 1 -s 200 -a $add
    ./write_to_server.sh -b 1 -w 1 -g 1 -s 400 -a $add
    ./write_to_server.sh -b 1 -w 1 -g 1 -s 600 -a $add
    ./write_to_server.sh -b 1 -w 1 -g 1 -s 800 -a $add
    ./write_to_server.sh -b 1 -w 1 -g 1 -s 1000 -a $add
fi

if [[ $compareType == "query" || $compareType == "all" ]]; then
    if [ $queryType == "cgo" ]; then
        echo "run query test, query method is cgo"
        for s in {100,1000};do
	        ./read.sh -w 100 -g 1 -s $s -a $add
            for j in {50,16,1};do
                ./read.sh -w $j -g 0 -s $s -a $add
            done
        done
    elif [ $queryType == "all" ]; then
        echo "run query test, query method is restful"
        for s in {100,1000};do
	        ./read_2.sh -w 100 -g 1 -s $s -n 'fast' -a $add
            for j in {50,16,1};do
                ./read_2.sh -w $j -g 0 -s $s -n 'fast' -a $add
            done
        done
    fi
fi
