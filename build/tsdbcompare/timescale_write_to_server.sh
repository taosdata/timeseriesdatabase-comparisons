#!/bin/bash

helpdoc() {
cat <<EOF
简述：

        本shell脚本用于进行timescaleDB和TDengine的插入性能对比测试

用法:

        $0	-b <batchsize>	-w <workers>	-i <interval>	-a <address>	-t <timestamp-start> \
	-e <timestamp-end>	-g <genate data>	-c <use-case>	-s <scale>

参数:

        ===========================================================================================
        -b      batchsize,		uint,   设置batchsize大小(1~5000),默认5000;
        -------------------------------------------------------------------------------------------
        -w      workers,		uint,   设置工作客户端数量,默认16;
        -------------------------------------------------------------------------------------------
        -i      interval,		string, 设置数据插入步长,默认10s;
        -------------------------------------------------------------------------------------------
        -a      address,		string, 设置TDengine或timescaleDB服务器地址,默认127.0.0.1;
        -------------------------------------------------------------------------------------------
        -t      timestamp-start,string, 设置开始时间戳,默认2018-01-01T00:00:00Z';
        -------------------------------------------------------------------------------------------
       	-e      timestamp-end,	string, 设置结束时间戳,默认2018-01-02T00:00:00Z';
        -------------------------------------------------------------------------------------------
        -g      gene,			uint,   设置是否生成数据(0:no,1:yes),默认1;
        -------------------------------------------------------------------------------------------
        -c      use-case,		string, 设置case类型(devops/iot/cpu-only),默认devops;
        -------------------------------------------------------------------------------------------
        -s      scale,			uint,   设置模拟的设备数量，默认100.
        ===========================================================================================

EOF
}

# Color setting
RED='\033[0;31m'
GREEN='\033[1;32m'
GREEN_DARK='\033[0;32m'
GREEN_UNDERLINE='\033[4;32m'
NC='\033[0m'

batchsize=5000
workers=16
usecase='devops'
gene=1
add=10.2.0.5
interval='10s'
scale=100
st='2018-01-01T00:00:00Z'
et='2018-01-02T00:00:00Z'

echo -e "${GREEN}====================This test option: scale: ${scale}, batchsize: ${batchsize}, worker: ${workers} =====================${NC}"
echo -e "${GREEN}------------------------------------- timescaledb test start -------------------------------------${NC}"

while getopts "hb:w:g:a:i:s:t:e:c:" opt
do
    case $opt in
      h)
        helpdoc
        exit 0
        ;;
      b)
        echo "batchsize:$OPTARG"
        batchsize=$OPTARG
        ;;
      w)
        echo "workers:$OPTARG"
        workers=$OPTARG
        ;;
      g)
        echo "whether generate data:$OPTARG"
        gene=$OPTARG
        ;;
      i)
        echo "interval:$OPTARG"
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
      c)
        echo "use-case:$OPTARG"
        usecase=$OPTARG
        ;;
      ?)
        helpdoc
        exit 1
        ;;
    esac
done

if [ $# != 18 ];then
    echo "variables not all defined,use value as follow :"
    echo "generate data: "
    echo "  batchsize:$batchsize ,    workers:$workers ,    generate data: $gene ,"
    echo "  scale: $scale ,   interval: $interval ,   timestamp-start: $st ,    timestamp-stop: $et"
    echo "  address: $add ,   usecase:$usecase"
fi
echo "$a"

#################### 0 prepare data source ###############
datafile="data/timescaledb.gz"

if [[ $gene == 1 ]];then
    if [ ! -d "data" ]; then
        mkdir data
    fi
    echo
    echo "---------------Generating Data-----------------"

    echo
    echo "Prepare data for timescaleDB...."
    ../../bin/timescale_generate_data   --seed=123 --format="timescaledb" --log-interval=$interval \
      --scale=$scale --use-case=$usecase \
      --timestamp-start=$st --timestamp-end=$et  | gzip  > ${datafile}
fi

#################### 1 clear the old data ###############
ssh -tt root@$add >/dev/null << EOFEOF
su - postgres -c "/usr/local/bin/postgresql/bin/pg_ctl -w restart " >/dev/null 2>&1 
exit
EOFEOF

# su - postgres -c "$PGDATA/../bin/pg_ctl -w restart " >/dev/null 2>&1 

psql -U postgres -h $add<< EOF
	drop database test;
EOF

ssh -tt root@$add "du -s /usr/local/bin/postgresql/data" > size.txt
# du -s /usr/local/bin/postgresql/data > size.txt
TMP=`grep data size.txt`
TSDISKINIT=`echo $TMP|awk '{print $1}'`
echo TSDISKINIT is :$TSDISKINIT

#################### 2 load the new data ###############
echo
echo -e "Start test timescaledb, result in ${GREEN}Green line${NC}"
pwd

loadcmd="cat ${datafile} | gunzip | ../../bin/timescale_load --workers=$workers --db-name='test' --batch-size=$batchsize --host=$add > tsres.log 2>&1"
echo $loadcmd
# $($loadcmd)
cat ${datafile} | gunzip | ../../bin/timescale_load --workers=$workers --db-name='test' --batch-size=$batchsize --host=$add  > tsres.log  2>&1  

#################### 3 data processing  ###############
TIMESCALE=`grep "rows/sec" tsres.log | tail -n 1`
TIMESCALEVALUE=`grep "metrics/sec" tsres.log | tail -n 1`

echo
echo -e "${GREEN}timescaledb writing result:${NC}"
echo -e "${GREEN}$TIMESCALE${NC}"
echo -e "${GREEN}$TIMESCALEVALUE${NC}"

DATA=`echo $TIMESCALE|awk '{print $2}'`
TMP=`echo $TIMESCALE|awk '{print($5)}'`
TSWTM=`echo ${TMP%s*}`
Pointrate=`echo $TIMESCALE | awk '{print $11}'`
Valuerate=`echo $TIMESCALEVALUE | awk '{print $11}'`


#################### 4 close the server  ###############
echo -e "${GREEN}Stop Timescaledb ${NC}"

ssh -tt root@$add >/dev/null << EOFEOF
# su - postgres -c "/usr/local/bin/postgresql/bin/pg_ctl -w stop" >/dev/null 2>&1
su - postgres -c "/usr/local/bin/postgresql/bin/pg_ctl -w stop " 
exit
EOFEOF

# su - postgres -c "$PGDATA/../bin/pg_ctl -w stop " >/dev/null 2>&1 

echo -e "${GREEN}sleep 10 seconds${NC}"
sleep 10
ssh -tt root@$add "du -s /usr/local/bin/postgresql/data" > size.txt
# du -s /usr/local/bin/postgresql/data > size.txt
TMP=`grep data size.txt`
TSDISKEND=`echo $TMP|awk '{print $1}'`
echo	TSDISKEND is : $TSDISKEND

# TSDISK=$(awk 'BEGIN{printf "%.2f\n" ,( '$TSDISKEND' - '$TSDISKINIT')/1024.0/1024.0}')
TSDISK=` echo "scale=2; ($TSDISKEND-$TSDISKINIT)/1024/1024" |bc `
# echo $TSDISK1

#################### 5 printf the result ###############
echo
echo
echo    "======================================================"
echo    "             tsdb performance comparision             "
printf  "       worker:%-4.2f      |       batch:%-4.2f      \n" $workers $batchsize
echo    "------------------------------------------------------"
echo -e "       Writing $DATA records test takes:          "
printf  "       Timescaledb      |       %-4.2f Seconds    \n" $TSWTM
echo    "------------------------------------------------------"
echo -e "       Writing $DATA records point rate:          "
printf  "       Timescaledb      |       %-4.2f point/s    \n" $Pointrate
echo    "------------------------------------------------------"
echo -e "       Writing $DATA records value rate:          "
printf  "       Timescaledb      |       %-4.2f value/s    \n" $Valuerate
echo    "------------------------------------------------------"
echo -e "       Writing $DATA records test disk:          "
printf  "       Timescaledb      |       %-4.2f GB          \n" $TSDISK
echo    "======================================================"

#################### 6 restart the server ###############

rm -f tsres.log
rm -f size.txt

ssh -tt root@$add >/dev/null << EOFEOF
su - postgres -c "/usr/local/bin/postgresql/bin/pg_ctl -w start " >/dev/null 2>&1 
exit
EOFEOF
echo -e "${GREEN}------------------------------------- timescaledb test stop -------------------------------------${NC}"
