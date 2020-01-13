# Performance comparisons between InfluxDB and TDengine
This project is a fork of [InfluxDB comparisions project](https://github.com/influxdata/influxdb-comparisons). The testing methodology and test procedure keeps the same as origin project and we just extend the data loading/quering module to support TDengine format and add serveral query test cases. Detailed testing methodology and procedure please refer to the origin project.

Briefly, this comparision test generates devops data and writes into different format, and loads into the database accordingly, then perform the same queries, finally counts the time consumed.

To make the test procedure simple, we use the docker to run the database so that you don't need to install the InfluxDB or TDengine, just pull the docker image and `docker run` it.

Current databases supported:

+ InfluxDB
+ TDengine

## Prerequisite

- A linux server with more than 10 GB free space on the disk, which is needed for the testing data.
- root privilege. During the test the TDengine client software is needed to install on the linux server which need login as root user.

## Prepare for test

Download the test package from [download site](http://www.taosdata.com/download/tsdbcompare.tar.gz) and unpack it.
```sh
tar -zxf tsdbcompare.tar.gz
```
The unpacked directory structure looks like:
```sh
.
├── bin
│   ├── bulk_data_gen
│   ├── bulk_load_influx
│   ├── bulk_load_tdengine
│   ├── bulk_query_gen
│   ├── query_benchmarker_influxdb
│   └── query_benchmarker_tdengine
├── config
│   ├── TDDashboardSchema.toml
│   └── TDengineSchema.toml
├── data
├── prepare.sh
├── run.sh
└── TDengine-client-1.6.4.5.tar.gz
```
The `bin` directory contains the data generation, load and query binary files built from source code. You can build them by yourself and replace them. The source files locate in `cmd` directory of this repository. 

The `config` directory contains the configuration of TDengine schema files.

The `data` directory is empty and the generated data will appears in this directory.

The `prepare.sh` is the script that will install docker and TDengine client.
```sh
#!/bin/bash
set -x

#install docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh


#install tdengine client
tar -zxf TDengine-client-1.6.4.5.tar.gz
cd TDengine-client-1.6.4.5
./install_client.sh
cd ..

#pull influxdb and tdengine docker images
docker pull influxdb 
docker pull tdengine/tdengine:v1.6.4.5.c 
```
If the docker is already installed, you can skip the `prepare.sh` script and install the TDengine client by:
```sh
#install tdengine client
tar -zxf TDengine-client-1.6.4.5.tar.gz
cd TDengine-client-1.6.4.5
./install_client.sh
cd ..
```
Then pull the influxdb and TDengine.

If everything is ok, you can proceed to the test procedure.

## Run the test

It's quite simple to run the test, just run the `run.sh` and wait. After several minutes, the result will be print in the screen.

```

---------------Generating Data-----------------

Prepare data for InfluxDB....
using random seed 123
2020/01/09 16:11:13 Using sampling interval 1s
2020/01/09 16:11:27 Written 7776000 points, 87264000 values, took 14.086559 seconds

Prepare data for TDengine....
using random seed 123
2020/01/09 16:11:27 Using sampling interval 1s
2020/01/09 16:11:42 Written 7776000 points, 87264000 values, took 15.460238 seconds

------------------Writing Data-----------------


Start test TDengine, result in Green line
2020/01/09 16:11:50 Creating database ----
2020/01/09 16:11:50 Starting workers ----

TDengine writing result:
loaded 7776000 items in 25.118785sec with 20 workers (mean point rate 309569.11/s, mean value rate 3474053.33/s, 41.06MB/sec from stdin)

Start test InfluxDB, result in Green line
2020/01/09 16:12:26 Ingestion rate control is off
2020/01/09 16:12:27 First statistic report received

InfluxDB writing result:
loaded 7776000 items in 94.995382sec with 20 workers (mean point rate 81856.610822/sec, mean value rate 918613.077000/s, 35.08MB/sec from stdin)

------------------Querying Data-----------------


start query test, query max from 8 hosts group by 1 hour, TDengine

using random seed 123
TDengine max cpu, rand    8 hosts : 1000 points
2020/01/09 16:14:12 Started querying with 50 workers
2020/01/09 16:14:12 Waiting for workers to finish

TDengine query test case 1 result:
wall clock time: 1.010856sec
using random seed 123
TDengine max cpu, rand    8 hosts,  by 1 hour: 1000 points
2020/01/09 16:14:13 Started querying with 50 workers
2020/01/09 16:14:24 Waiting for workers to finish

TDengine query test case 2 result:
wall clock time: 12.019893sec
using random seed 123
TDengine max cpu, rand    8 hosts, rand 12h0m0s by 10m: 1000 points
2020/01/09 16:14:25 Started querying with 50 workers
2020/01/09 16:14:31 Waiting for workers to finish

TDengine query test case 3 result:
wall clock time: 6.716956sec
using random seed 123
TDengine max cpu, rand    8 hosts, rand 1h0m0s by 1m: 1000 points
2020/01/09 16:14:31 Started querying with 50 workers
2020/01/09 16:14:33 Waiting for workers to finish

TDengine query test case 4 result:
wall clock time: 1.978717sec

start query test, query max from 8 hosts group by 1hour, Influxdb

using random seed 123
InfluxDB (InfluxQL) max cpu, rand    8 hosts : 1000 points
2020/01/09 16:14:43 Started querying with 50 workers
2020/01/09 16:16:04 Waiting for workers to finish
InfluxDB query test case 1 result:
wall clock time: 82.508923sec
using random seed 123
InfluxDB (InfluxQL) max cpu, rand    8 hosts, by 1h: 1000 points
2020/01/09 16:16:06 Started querying with 50 workers
2020/01/09 16:17:47 Waiting for workers to finish
InfluxDB query test case 2 result:
wall clock time: 102.987983sec
using random seed 123
InfluxDB (InfluxQL) max cpu, rand    8 hosts, rand 12h0m0s by 10m: 1000 points
2020/01/09 16:17:49 Started querying with 50 workers
2020/01/09 16:18:30 Waiting for workers to finish
InfluxDB query test case 3 result:
wall clock time: 42.663898sec
using random seed 123
InfluxDB (InfluxQL) max cpu, rand    8 hosts, rand 1h0m0s by 1m: 1000 points
2020/01/09 16:18:32 Started querying with 50 workers
2020/01/09 16:18:37 Waiting for workers to finish
InfluxDB query test case 4 result:
wall clock time: 5.172598sec


======================================================
             tsdb performance comparision             
======================================================
       Writing 7776000 records test takes:          
       InfluxDB           |       95.00 Seconds    
       TDengine           |       25.12 Seconds    
------------------------------------------------------
                   Query test cases:                
 case 1: select the max(value) from all data    
 filtered out 8 hosts                                 
       Query test case 1 takes:                      
       InfluxDB           |       82.51 Seconds    
       TDengine           |       1.01 Seconds    
------------------------------------------------------
 case 2: select the max(value) from all data          
 filtered out 8 hosts with an interval of 1 hour     
 case 2 takes:                                       
       InfluxDB           |       102.99 Seconds    
       TDengine           |       12.02 Seconds    
------------------------------------------------------
 case 3: select the max(value) from random 12 hours
 data filtered out 8 hosts with an interval of 10 min         
 filtered out 8 hosts interval(1h)                   
 case 3 takes:                                       
       InfluxDB           |       42.66 Seconds    
       TDengine           |       6.72 Seconds    
------------------------------------------------------
 case 4: select the max(value) from random 1 hour data  
 data filtered out 8 hosts with an interval of 1 min         
 case 4 takes:                                        
       InfluxDB           |       5.17 Seconds    
       TDengine           |       1.98 Seconds    
------------------------------------------------------
```
