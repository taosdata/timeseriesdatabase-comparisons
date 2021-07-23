#!/bin/bash

PROG_NAME=$1
SDIR_PATH="monitor_log"
USAGE_REPORT=${PROG_NAME}_$(date "+%Y%m%d_%H%M%S").csv

if [ ! -d ${SDIR_PATH} ]; then
	mkdir ${SDIR_PATH}
fi
cd ${SDIR_PATH}
echo "log is enabled. csv_file \"${USAGE_REPORT}\""


max_cpu=0.0
min_cpu=100000.0
avg_cpu=0.0

max_mem=0.0
min_mem=100000.0
avg_mem=0.0

i=0

while [ 1 ]
do
	pid=`pidof $PROG_NAME`
	if [ "$pid" != "" ]; then
		i=`expr $i + 1`
		echo $i
		cpu_usage=`top -bn1 -p $pid | awk -F " " 'NR==8 {print $9}'`
		mem_usage=`top -bn1 -p $pid | awk -F " " 'NR==8 {print $10}'`
		echo "CPU usage: ${cpu_usage}, Mem usage: ${mem_usage}" | tee -a $USAGE_REPORT
	
		if [ $(echo "$cpu_usage > $max_cpu" | bc) -eq 1 ]; then
			max_cpu=$cpu_usage
		fi

		if [ $(echo "$cpu_usage < $min_cpu" | bc) -eq 1 ]; then
			min_cpu=$cpu_usage
		fi

		diff=`echo "$cpu_usage - $avg_cpu"|bc`
		avg=`echo "scale=2; $diff / $i"|bc`
		avg_cpu=`echo "$avg_cpu +$avg"|bc`
		echo "max cpu usage: ${max_cpu}, min cpu usage: ${min_cpu}, avg CPU usage: ${avg_cpu}" | tee -a $USAGE_REPORT
		
		if [ $(echo "$mem_usage > $max_mem" | bc) -eq 1 ]; then
                        max_mem=$mem_usage
                fi

                if [ $(echo "$mem_usage < $min_mem" | bc) -eq 1 ]; then
                        min_mem=$mem_usage
                fi

                diff1=`echo "$mem_usage - $avg_mem"|bc`
                avg1=`echo "scale=2; $diff1 / $i"|bc`
                avg_mem=`echo "$avg_mem +$avg1"|bc`
                echo "max mem usage: ${max_mem}, min mem usage: ${min_mem}, avg mem usage: ${avg_mem}" | tee -a $USAGE_REPORT

	else
		echo "process id:${pid} not found"
		exit 0
	fi

	
	sleep 1
done
