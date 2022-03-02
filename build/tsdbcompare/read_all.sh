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


query()
{
	echo "------------------ Compile go start !-----------------"

	cd ../../cmd/bulk_query_gen
	#pwd
	rm bulk_query_gen
	go build
	sleep 1
	ls -l bulk_query_gen
	cp bulk_query_gen  ../../build/tsdbcompare/bin
	#pwd
	#printf bulk_query_gen|md5sum
	
	cd ../../cmd/query_benchmarker_influxdb
	#pwd
	rm query_benchmarker_influxdb
	go build
	sleep 1
	ls -l query_benchmarker_influxdb
	cp query_benchmarker_influxdb  ../../build/tsdbcompare/bin
	#pwd
	#printf query_benchmarker_influxdb|md5sum

	cd ../../cmd/query_benchmarker_tdengine
	#pwd
	rm query_benchmarker_tdengine
	go build
	sleep 1
	ls -l query_benchmarker_tdengine
	cp query_benchmarker_tdengine  ../../build/tsdbcompare/bin
	#pwd
	#printf query_benchmarker_tdengine|md5sum
	
	echo "------------------Compile go over !-----------------"

	ls -l ../../build/tsdbcompare/bin
	sleep 2
	pwd
	cd ../../build/tsdbcompare/
	pwd
	
}



echo
echo "------------------Part 1 : Comparison of different random hosts -----------------"
echo


echo 
echo  "------------------Test case 1.1 : Comparison use 1 random host by cgo -----------------"
echo

#测试用例1.1，查询所有数据中，用1个hostname标签进行匹配，interface='cgo'
# 编译
pwd
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_1host ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
pwd
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_1host ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_1host.sh -a ${add}

# echo 
# echo  "------------------Test case 1.1 : Comparison use 1 random host by rest-----------------"
# echo

# #测试用例1.1，查询所有数据中，用1个hostname标签进行匹配，interface='fast'
#./read_1host_rest.sh

echo 
echo  "------------------Test case 1.2 : Comparison use 8 random hosts by cgo -----------------"
echo

#测试用例1.2，查询所有数据中，用8个hostname标签进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_8host ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_8host ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_8host.sh -a ${add}

# echo 
# echo  "------------------Test case 1.2 : Comparison use 8 random hosts by rest-----------------"
# echo

#测试用例1.2，查询所有数据中，用8个hostname标签进行匹配，interface='fast'
#./read_8host_rest.sh

echo 
echo  "------------------Test case 1.3 : Comparison use 16 random hosts by cgo -----------------"
echo

#测试用例1.3，查询所有数据中，用16个hostname标签进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_16host ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_16host ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_16host.sh -a ${add}

# echo 
# echo  "------------------Test case 1.3 : Comparison use 16 random hosts by rest-----------------"
# echo

# #测试用例1.3，查询所有数据中，用16个hostname标签进行匹配，interface='fast'
#./read_16host_rest.sh


echo 
echo  "------------------Test case 1.4 : Comparison use 32 random hosts by cgo -----------------"
echo

#测试用例1.4，查询所有数据中，用32个hostname标签进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_32host ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_32host ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_32host.sh -a ${add}

# echo 
# echo  "------------------Test case 1.4 : Comparison use 32 random hosts by rest-----------------"
# echo

# #测试用例1.4，查询所有数据中，用32个hostname标签进行匹配，interface='fast'
#./read_32host_rest.sh


echo 
echo  "------------------Test case 1.5 : Comparison use 64 random hosts by cgo -----------------"
echo

#测试用例1.5，查询所有数据中，用64个hostname标签进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_64host ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_64host ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_64host.sh -a ${add}

# echo 
# echo  "------------------Test case 1.5 : Comparison use 64 random hosts by rest-----------------"
# echo

# #测试用例1.5，查询所有数据中，用64个hostname标签进行匹配，interface='fast'
# ./read_64host_rest.sh


echo 
echo  "------------------Test case 1.6 : Comparison use 128 random hosts by cgo -----------------"
echo

#测试用例1.6，查询所有数据中，用128个hostname标签进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_128host ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_128host ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_128host.sh -a ${add}

# echo 
# echo  "------------------Test case 1.6 : Comparison use 128 random hosts by rest-----------------"
# echo

# #测试用例1.6，查询所有数据中，用128个hostname标签进行匹配，interface='fast'
# ./read_128host_rest.sh


echo 
echo  "------------------Test case 1.7 : Comparison use 256 random hosts by cgo -----------------"
echo

#测试用例1.7，查询所有数据中，用256个hostname标签进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_256host ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_256host ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_256host.sh -a ${add}

# echo 
# echo  "------------------Test case 1.7 : Comparison use 256 random hosts by rest-----------------"
# echo

# #测试用例1.7，查询所有数据中，用256个hostname标签进行匹配，interface='fast'
# ./read_256host_rest.sh


echo 
echo  "------------------Test case 1.8 : Comparison use 512 random hosts by cgo -----------------"
echo

#测试用例1.8，查询所有数据中，用512个hostname标签进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_512host ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_512host ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_512host.sh -a ${add}

# echo 
# echo  "------------------Test case 1.8 : Comparison use 512 random hosts by rest-----------------"
# echo

# #测试用例1.8，查询所有数据中，用512个hostname标签进行匹配，interface='fast'
# ./read_512host_rest.sh



echo
echo "------------------Part 2 : Comparison of different hours -----------------"
echo


echo 
echo  "------------------Test case 2.1 : Comparison use 1 hour by cgo -----------------"
echo

#测试用例2.1，查询所有数据中，用 1 hour 时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_1hour ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_1hour ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_1hour.sh -a ${add} 



echo 
echo  "------------------Test case 2.2 : Comparison use 2 hours by cgo -----------------"
echo

#测试用例2.2，查询所有数据中，用 2 hours 时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_2hour ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_2hour ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_2hour.sh -a ${add}



echo 
echo  "------------------Test case 2.3 : Comparison use 4 hours by cgo -----------------"
echo

#测试用例2.3，查询所有数据中，用 4 hours 时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_4hour ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_4hour ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_4hour.sh -a ${add}



echo 
echo  "------------------Test case 2.4 : Comparison use 8 hours by cgo -----------------"
echo

#测试用例2.4，查询所有数据中，用 8 hours 时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_8hour ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_8hour ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_8hour.sh -a ${add}



echo 
echo  "------------------Test case 2.5 : Comparison use 12 hours by cgo -----------------"
echo

#测试用例2.5，查询所有数据中，用 12 hours 时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_12hour ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_12hour ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_12hour.sh -a ${add}




echo
echo "------------------Part 3 : Comparison of different function -----------------"
echo


echo 
echo  "------------------Test case 3.1 : Comparison max、count、first、last use 1 hour 8 host by cgo -----------------"
echo

#测试用例3.1，查询所有数据max、count、first、last中，用 1 hour 时间段 8 hosts 进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_count ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_count ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_count.sh -a ${add}



echo 
echo  "------------------Test case 3.2 : Comparison top(10) use 1 hour 8hosts by cgo -----------------"
echo

#测试用例3.2，查询所有数据top(10)中，用 1 hour 时间段 8 hosts 进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_top10 ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_top10 ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_top.sh -a ${add}



echo 
echo  "------------------Test case 3.3 :  Comparison max、count、first、last、twa use 1 hour 8 hosts by cgo  -----------------"
echo

#测试用例3.3，查询所有数据max、count、first、last、twa中，用 1 hours 时间段  8 hosts 进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_count_percentile ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_count_percentile ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_count_percentile.sh -a ${add}




echo
echo "------------------Part 4 : Comparison of different hosts -----------------"
echo


echo 
echo  "------------------Test case 4.1 : Comparison use 16 hosts by cgo -----------------"
echo

#测试用例4.1，查询所有数据中，用 16 hosts 全部时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_select16 ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_select16 ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_select16.sh -a ${add}



echo 
echo  "------------------Test case 4.2 : Comparison use 32 hosts by cgo -----------------"
echo

#测试用例4.2，查询所有数据中，用 32 hours 全部时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_select32 ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_select32 ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_select32.sh -a ${add}



echo 
echo  "------------------Test case 4.3 : Comparison use 64 hosts by cgo -----------------"
echo

#测试用例4.3，查询所有数据中，用 64 hours 全部时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_select64 ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_select64 ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_select64.sh -a ${add}



echo 
echo  "------------------Test case 4.4 : Comparison use 128 hosts by cgo -----------------"
echo

#测试用例4.4，查询所有数据中，用 128 hours 全部时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_select128 ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_select128 ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_select128.sh -a ${add}



echo 
echo  "------------------Test case 4.5 : Comparison use 256 hosts by cgo -----------------"
echo

#测试用例4.5，查询所有数据中，用 256 hours 全部时间段进行匹配，interface='cgo'
# 编译
cp ../../bulk_query_gen/tdengine/tdengine_devops_common.go_select256 ../../bulk_query_gen/tdengine/tdengine_devops_common.go
ls -l ../../bulk_query_gen/tdengine/tdengine_devops_common.go
cp ../../bulk_query_gen/influxdb/influx_devops_common.go_select256 ../../bulk_query_gen/influxdb/influx_devops_common.go
ls -l ../../bulk_query_gen/influxdb/influx_devops_common.go
query
# 执行对比程序
./read_select256.sh -a ${add}

