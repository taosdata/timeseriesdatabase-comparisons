# InflxuDB和TDengine的性能对比测试工具

### 前言
[TDengine开源项目](https://github.com/taosdata/TDengine)里已经包含了性能对比测试的工具源代码。[https://github.com/taosdata/TDengine/tests/comparisonTest](https://github.com/taosdata/TDengine/tree/develop/tests/comparisonTest)，并基于这个开源的测试工具开展了[TDengine和InfluxDB对比测试](https://www.taosdata.com/blog/2019/07/19/tdengine与influxdb对比测试/)，[TDengine和OpenTSDB对比测试](https://www.taosdata.com/blog/2019/08/21/tdengine与opentsdb对比测试/),[TDengine和Cassandra对比测试](https://www.taosdata.com/blog/2019/08/14/tdengine与cassandra对比测试/)等一系列性能对比测试。为了更客观的对比TDengine和其他时序数据库的性能差异，本项目采用由InfluxDB团队开源的性能对比测试工具来进行对比测试，相同的数据产生器，相同的测试用例，相同的测试方法，以保证测试的客观公平。

### 简介
本项目是基于InfluxDB发布的一个[性能对比测试项目](https://github.com/influxdata/influxdb-comparisons)的基础上开发的。数据产生模块可以模拟Devops场景下多台服务器产生大量监控数据。数据写入程序可以根据不同的数据库格式，将产生的模拟数据以不同的格式写入到不同数据库里，以测试写入性能。查询模块以相同的查询类型产生相同的查询任务，以各数据库自己的格式进行查询，并统计查询消耗的时间，来测试查询性能。

为了让测试过程更简单，本测试采用Docker容器方式来测试，所有被测的数据库都以容器的方式，从Dockerhub拉取下来，并设定固定的IP地址运行，便于脚本执行。容器镜像都是公开发布的，能保证测试的公平公正。

本测试项目目前支持以下时序数据库的对比测试
+ InfluxDB
+ TDengine

本项目的Github链接：https://github.com/liu0x54/timeseriesdatabase-comparisons
## 前提条件
为了开展测试，需要准备以下条件
- 一台linux服务器，包含10GB的空闲硬盘空间，用于存储产生的测试数据。因为测试模拟数据先生成并写入硬盘文件，由数据加载程序从文件中读取一条条的数据写入语句，写入时序数据库。这种方式能够将数据产生过程中的性能差异排除。
- root权限。测试过程需要用root权限来安装一个TDengine的客户端驱动，用于TDengine数据加载程序的调用。TDengine数据写入采用go语言调用C语言连接器的方式。

## 准备测试

先从[下载地址](http://www.taosdata.com/download/tsdbcompare.tar.gz)下载我们已经制作好的测试工具包，解压到本地。
```sh
tar -zxf tsdbcompare.tar.gz
```
解压后的目录结构如下:
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
`bin` 目录里是提前编译好的可执行文件，包括数据产生，数据加载，查询产生和查询加载。提前编译好方便大家下载即可用；如果有兴趣的同学也可以自己从源文件编译。源文件位于`cmd`下面的各个子目录里。可以自行编译后替换bin目录的文件。 

`config` 目录里是TDengine写入数据需要用到的schema配置文件，模拟数据产生的数据通过schema配置里的信息可以知道该往哪个表里存。

`data` 目录是用来存储测试过程中产生的数据文件。本测试采用先产生模拟数据，并将模拟数据按各数据库的写入格式写到文件里，再用加载程序从文件里读取按格式写好的语句往各数据库里加载的方式来开展测试。这样的方法，能够将原始数据转换成不同的格式的过程带来的差异进行屏蔽，更纯粹的对比数据库的写入性能。

`prepare.sh` 是用来准备测试环境的脚本，包含三部分，1.安装docker程序，2.安装TDengine的客户端，3.拉取influxDB和TDengine的Docker镜像。
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
如果目标系统上已经安装了docker程序，就不用执行这个`prepare.sh`脚本了，可以直接按脚本里的第二、三部分去安装TDengine Client和拉取对应的Docker镜像。

在上面的步骤都执行完成，并确认成功后，可以开展测试工作了。

## 开展测试

在整个测试过程中，建议另开一个终端，运行top，查看系统的CPU和内存占用情况
```sh
top
```

### 写入测试
本测试包提供了一个`run.sh`脚本，自动执行将docker容器按指定IP地址运行起来，然后产生数据，写入数据文件，并写入时序数据库。
数据产生和写入由以下两条命令完成
```sh
#产生模拟数据并写入数据文件
bin/bulk_data_gen -seed 123 -format influx-bulk -sampling-interval 1s -scale-var 10 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z" >data/influx.dat
bin/bulk_data_gen -seed 123 -format tdengine -sampling-interval 1s -tdschema-file config/TDengineSchema.toml -scale-var 10 -use-case devops -timestamp-start "2018-01-01T00:00:00Z" -timestamp-end "2018-01-02T00:00:00Z"  > data/tdengine.dat
```
解释一下以上的命令：按influxDB/TDengine的格式，以1秒一条数据的产生频率，模拟10台设备，以devops场景产生24小时的数据，并写入influx.dat文件。
Devops模型下，一台服务器会产生9类数据，分别是cpu，disk，mem，等，因此总共会产生7776000条数据记录。
数据文件完成后，就开始数据写入测试：
```sh
#数据写入数据库
cat data/influx.dat  |bin/bulk_load_influx --batch-size=5000 --workers=20 --urls="http://172.15.1.5:8086" 

cat data/tdengine.dat |bin/bulk_load_tdengine --url 172.15.1.6:0 --batch-size 300   -do-load -report-tags n1 -workers 20 -fileout=false 
```
上面命令的含义是以每批次写入5000/300条记录，分20个线程，将数据文件读取出来后写入influxDB/TDengine中
### 查询测试
在完成写入后，就开始查询测试。
查询测试设定了四个查询用例的语句，每个查询语句都执行1000遍，然后统计总的查询用时:
1. 测试用例1，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据的最大值。
```sh
#TDengine
bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-all -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="http://172.15.1.6:6020" -workers 50 -print-interval 0

#InfluxDB
bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-all -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://172.15.1.5:8086"  -workers 50 -print-interval 0
```
2. 测试用例2，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1小时为粒度，查询每1小时的最大值。
```sh
#TDengine
bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-allbyhr -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="http://172.15.1.6:6020" -workers 50 -print-interval 0

#InfluxDB
bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-allbyhr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://172.15.1.5:8086"  -workers 50 -print-interval 0
```
3. 测试用例3，随机查询12个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以10分钟为粒度，查询每10分钟的最大值。
```sh
#TDengine
bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-12-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="http://172.15.1.6:6020" -workers 50 -print-interval 0

#InfluxDB
bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-12-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://172.15.1.5:8086"  -workers 50 -print-interval 0
```
4. 测试用例4，随机查询1个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1分钟为粒度，查询每1分钟的最大值。
```sh
#TDengine
bin/bulk_query_gen  -seed 123 -format tdengine -query-type 8-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_tdengine  -urls="http://172.15.1.6:6020" -workers 50 -print-interval 0

#InfluxDB
bin/bulk_query_gen  -seed 123 -format influx-http -query-type 8-host-1-hr -scale-var 10 -queries 1000 | bin/query_benchmarker_influxdb  -urls="http://172.15.1.5:8086"  -workers 50 -print-interval 0
```
查询过程结束后，将测试结果以以下格式打印出来

```
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
## 结果分析
通过本测试包产生的数据和相关的写入、查询用例测试可以看出，TDengine在性能上相比InfluxDB有较大的优势。
细致分析下来可以有以下结论：
- 写入性能：相同数据源InfluxDB写入用时约是TDengine的4倍
- 全部数据聚合计算查询：InfluxDB查询用时约为TDengine的80倍
- 全部数据聚合计算查询以小时为颗粒聚合结果：InfluxDB查询用时约为TDengine的10倍
- 随机选取12小时的数据聚合计算查询以10分钟为颗粒聚合结果：InfluxDB用时约为TDengine的6倍
- 随机选取1小时的数据聚合计算查询以1分钟为颗粒聚合结果：InfluxDB用时约为TDengine的2.5倍

通过top命令的观察，我们可以看到，测试用例执行时，InfluxDB的CPU占用率基本达到满负荷，以4核CPU的服务器为例，经常达到390%以上；而TDengine的CPU占用率则低很多。

## 从源文件编译
本项目中的二进制程序是基于1.6.4.5版本的TDengine编译的，因为要用到TDengine的客户端，而TDengine客户端版本必须和TDengine的服务端的版本必须匹配，因此直接使用本项目的二进制的程序连接其他版本的TDengine或镜像时，会遇到连接不上的问题。
如果遇到这种情况，或者对于从源代码编译二进制文件感兴趣的使用者，就需要从本项目的源代码编译出二进制代码后再运行本测试程序。下面介绍编译的步骤

#### 前提条件
- 获取目标版本的TDengine，在编译所在的机器上安装和服务端一样的TDengine服务端版本。TDengine服务端版本可以是从源代码编译出来的，也可以是从Taosdata官网获取的编译好的安装包
- 如果是从源代码编译的TDengine，注意编译出来的lib文件，以及源代码中的taos.h文件必须放到系统的目录里，这样才能正确被引用。可以参考以下代码，将lib文件和taos.h文件放到系统目录里去。如果是从安装包安装的话，可以跳过这一步：
```
cp  (TDengine项目的build和源代码目录所在的路径)/build/build/lib/libtaos.so /usr/lib/libtaos.so.1
ln -s /usr/lib/libtaos.so.1 /usr/lib/libtaos.so
cp (TDengine项目的build和源代码目录所在的路径)/TDengine/src/inc/taos.h /usr/include/
```

- 安装最新版本的golang语言程序,参考golang网站相关指导页面[下载安装golang](https://golang.org/dl/)

#### 开始编译
进入本项目的cmd目录，这个目录下就是各个可执行程序的源文件目录。
其中：
bulk_data_gen 目录中的main.go是生成测试所需的数据的程序
bulk_load_* 是针对不同数据库写入数据的程序
bulk_query_gen 目录中的main.go是生成测试所需的查询语句程序
query_benchmarker_* 是针对不同数据库的查询测试执行程序

因此，如果要测试TDengine，就编译以下目录中的main.go文件：
bulk_data_gen
bulk_load_tdengine
bulk_query_gen
query_benchmarker_tdengine

编译方法如下：
进入对应的文件路径下，执行：
```
go build
```
顺利的话就回生成一个可执行的二进制文件，名称为该文件夹的名称
```
TomdeMacBook-Pro:bulk_data_gen tom$ go build
TomdeMacBook-Pro:bulk_data_gen tom$ ls
bulk_data_gen   main.go
TomdeMacBook-Pro:bulk_data_gen tom$ 
```
如果遇到某些引用的包没有找到的话，可以使用go get对应的包，获取到本地后再编译。
比如如果找不到github.com/taosdata/TDengine/src/connector/go/src/taosSql，那么执行
```
go get github.com/taosdata/TDengine/src/connector/go/src/taosSql
```
在以上语句执行成功后就可以继续编译了。

编译出二进制代码后，就可以参考run.sh脚本中的命令顺序，手动执行测试了。

## 命令参数详解
测试程序需要输入对应的参数，必须在正确理解了参数的含义的基础上，正确的设置参数，再执行测试程序。对于每个参数的含义如下
#### bulk_data_gen程序用到的参数
```
Usage of ./bulk_data_gen:
  -config-file string
        Simulator config file in TOML format (experimental) 
        说明：这个参数暂时用不上，可以不用管
  -cpu-profile file
        Write CPU profile to file
        说明：这个参数暂时用不上，可以不用管
  -debug int
        Debug printing (choices: 0, 1, 2) (default 0).
        说明：这个参数暂时用不上，可以不用管
  -format string
        Format to emit. (choices: influx-bulk, es-bulk, es-bulk6x, cassandra, mongo, opentsdb, timescaledb-sql, timescaledb-copyFrom, graphite-line, splunk-json, tdengine) (default "influx-bulk")
        说明：这个参数设置测试数据输出的格式，可以输入choices后面的选项，如果要生成TDengine的测试数据，就带tdengine
  -interleaved-generation-group-id uint
        Group (0-indexed) to perform round-robin serialization within. Use this to scale up data generation to multiple processes.
        说明：这个参数暂时用不上，可以不用管
  -interleaved-generation-groups uint
        The number of round-robin serialization groups. Use this to scale up data generation to multiple processes. (default 1)
        说明：这个参数暂时用不上，可以不用管
  -sampling-interval duration
        Simulated sampling interval. (default 10s)
        说明：这个参数设置模拟数据产生的时间戳的间隔，缺省10秒钟。
  -scale-var int
        Scaling variable specific to the use case. (default 1)
        说明：这个参数设置数据源的数量，如果为N，则模拟N个数据源产生数据。这个数量越大，相同条件下产生的数据量越大
  -scale-var-offset int
        Scaling variable offset specific to the use case.
        说明：这个参数设置的数据源的ID的偏移量，在多个程序同时产生模拟数据时，通过这个偏移量来区别每个测试程序产生的数据中的ID信息，避免重叠
  -seed int
        PRNG seed (default, or 0, uses the current timestamp).
        说明：这个是模拟数据的伪随机数种子，测试不同数据库时，这个种子建议设置成一样的，那么在相同的配置下，产生的数据值也是一致的
  -tdschema-file string
        TDengine schema config file in TOML format (experimental)
        说明：这个参数设置TDengine需要的配置文件的路径，TDengine需要配置文件来生成对应的表结构。本项目中用到了两个配置文件，分别是TDengineSchema.toml和TDDashboardSchema.toml,
        对于use-case是devops，iot时，请指定TDengineSchema.toml作为配置文件；对于use-case是dashboard时，请指定TDDashboardSchema.toml作为配置文件
  -timestamp-end string
        Ending timestamp (RFC3339). (default "2018-01-02T00:00:00Z")
        说明：模拟的测试数据时间戳结束的时间
  -timestamp-start string
        Beginning timestamp (RFC3339). (default "2018-01-01T00:00:00Z")
        说明：模拟的测试数据时间戳开始的时间
  -use-case string
        Use case to model. (choices: devops, iot, dashboard) (default "devops")
        说明：选择测试数据模型，不同的use-case模拟了不同的数据类型。
```
#### bulk_load_tdengine程序用到的参数
```
Usage of ./bulk_load_tdengine:
-batch-size int
        Batch size (input items). (default 100)
        说明：每次写入请求中向TDengine写入的记录数量，缺省为100条记录
  -do-load
        Whether to read data from file or from stand input. Set this flag to true to get input from file.
        说明：选择测试数据来源，是从文件中读取，还是从stdin中获取。如果设为true，则从文件中读取数据
  -file string
        Input file
        说明：如果-do-load设为true，则通过这个参数设定数据文件的路径
  -fileout
        if file out, will out put sql into file (default true)
        说明：如果设为true则不直接写数据库，而是将指令写到指定文件中，便于分析
  -report-database string
        Database name where to store result metrics (default "vehicle")
        说明：设定写入目标数据库中的database名称
  -report-host string
        Host to send result metrics
        说明：报告测试结果的服务器，暂时用不上
  -report-password string
        User password for Host to send result metrics
        说明：报告测试结果的服务器，暂时用不上
  -report-tags string
        Comma separated k:v tags to send  alongside result metrics (default "node1")
         说明：不用填，用缺省值即可
  -report-user string
        User for host to send result metrics
         说明：不用填，用缺省值即可
  -slavesource
        if slave source, will not create database
         说明：不用填，用缺省值即可
  -url string
        TDengine URL. (default "127.0.0.1:0")
        说明：TDengine数据库的地址信息，以ip地址加:0作为格式，比如TDengine运行22.23.123.32这个服务器上，则这里填写 22.23.123.32:0
  -use-case string
        Use case to set specific load behavior. Options: devops,iot,dashboard (default "devops")
        说明：不用填，用缺省值即可
  -workers int
        Number of parallel requests to make. (default 2)
        说明：设定多个写入线程的数量，用缺省值，或者在写入速度较慢时可以适当调大
```

#### bulk_query_gen用到的参数
```
Usage of ./bulk_query_gen:
  -document-format string
        Document format specification. (for mongo format 'simpleArrays'; leave empty for previous behaviour)
        说明：不用填，用缺省值即可
  -format string
        Format to emit. (Choices are in the use case matrix.) (default "influx-http")
        说明：查询语句的格式，填对应的数据库，如果测TDengine，则填tdengine
  -interleaved-generation-group-id uint
        Group (0-indexed) to perform round-robin serialization within. Use this to scale up data generation to multiple processes.
        说明：不用填，用缺省值即可
  -interleaved-generation-groups uint
        The number of round-robin serialization groups. Use this to scale up data generation to multiple processes. (default 1)
        说明：不用填，用缺省值即可
  -queries int
        Number of queries to generate. (default 1000)
        说明：测试查询的次数，可以填1000次，因为每次查询的时间非常短，可以填多次，测总的时间
  -query-interval duration
        Time interval query should ask for. (default 1h0m0s)
        说明：不用填，用缺省值即可
  -query-interval-type string
        Interval type query { window - either random or shifted, last - interval is defined relative to now() } (default "window")
        说明：不用填，用缺省值即可
  -query-type string
        Query type. (Choices are in the use case matrix.)
        说明：查询TDengine的话，从下面的选项中选择format: tdengine的填入
  -scale-var int
        Scaling variable (must be the equal to the scale-var used for data generation). (default 1)
        说明：这个参数和bulk_data_gen命令的参数保持一致，表示测试数据模型是模拟的多少个数据源产生的。这个数据会对查询语句中的设备ID产生影响
  -seed int
        PRNG seed (default, or 0, uses the current timestamp).
        说明：查询语句中用到的伪随机数种子，如果测试不同的数据库，确保这个种子是一致的，这样的话查询语句中的随机值可以保证一致
  -time-window-shift duration
        Sliding time window shift. (When set to > 0s, queries option is ignored - number of queries is calculated. (default -1ns)
        说明：不用填，用缺省值即可
  -timestamp-end string
        Ending timestamp (RFC3339). (default "2018-01-02T00:00:00Z")
        说明：不用填，用缺省值即可
  -timestamp-start string
        Beginning timestamp (RFC3339). (default "2018-01-01T00:00:00Z")
        说明：不用填，用缺省值即可
  -use-case string
        Use case to model. (Choices are in the use case matrix.) (default "devops")
        说明：查询TDengine的话，不用填，缺省devops
 The use case matrix of choices is:
  use case: devops, query type: groupby, format: influx-http
  use case: devops, query type: groupby, format: timescaledb
  use case: devops, query type: groupby, format: graphite
  use case: devops, query type: groupby, format: splunk
  use case: devops, query type: groupby, format: cassandra
  use case: devops, query type: groupby, format: es-http
  use case: devops, query type: groupby, format: influx-flux-http
  use case: devops, query type: 1-host-1-hr, format: cassandra
  use case: devops, query type: 1-host-1-hr, format: es-http
  use case: devops, query type: 1-host-1-hr, format: influx-http
  use case: devops, query type: 1-host-1-hr, format: mongo
  use case: devops, query type: 1-host-1-hr, format: graphite
  use case: devops, query type: 1-host-1-hr, format: influx-flux-http
  use case: devops, query type: 1-host-1-hr, format: opentsdb
  use case: devops, query type: 1-host-1-hr, format: timescaledb
  use case: devops, query type: 1-host-1-hr, format: splunk
  use case: devops, query type: 1-host-1-hr, format: tdengine
  use case: devops, query type: 1-host-12-hr, format: opentsdb
  use case: devops, query type: 1-host-12-hr, format: graphite
  use case: devops, query type: 1-host-12-hr, format: splunk
  use case: devops, query type: 1-host-12-hr, format: cassandra
  use case: devops, query type: 1-host-12-hr, format: influx-http
  use case: devops, query type: 1-host-12-hr, format: mongo
  use case: devops, query type: 1-host-12-hr, format: timescaledb
  use case: devops, query type: 1-host-12-hr, format: tdengine
  use case: devops, query type: 1-host-12-hr, format: es-http
  use case: devops, query type: 1-host-12-hr, format: influx-flux-http
  use case: devops, query type: 8-host-1-hr, format: opentsdb
  use case: devops, query type: 8-host-1-hr, format: timescaledb
  use case: devops, query type: 8-host-1-hr, format: es-http
  use case: devops, query type: 8-host-1-hr, format: influx-flux-http
  use case: devops, query type: 8-host-1-hr, format: influx-http
  use case: devops, query type: 8-host-1-hr, format: mongo
  use case: devops, query type: 8-host-1-hr, format: graphite
  use case: devops, query type: 8-host-1-hr, format: splunk
  use case: devops, query type: 8-host-1-hr, format: tdengine
  use case: devops, query type: 8-host-1-hr, format: cassandra
  use case: devops, query type: 8-host-12-hr, format: influx-http
  use case: devops, query type: 8-host-12-hr, format: tdengine
  use case: devops, query type: 8-host-allbyhr, format: influx-http
  use case: devops, query type: 8-host-allbyhr, format: tdengine
  use case: devops, query type: 8-host-all, format: influx-http
  use case: devops, query type: 8-host-all, format: tdengine
  use case: iot, query type: 1-home-12-hours, format: influx-flux-http
  use case: iot, query type: 1-home-12-hours, format: influx-http
  use case: iot, query type: 1-home-12-hours, format: timescaledb
  use case: iot, query type: 1-home-12-hours, format: cassandra
  use case: iot, query type: 1-home-12-hours, format: mongo
  use case: dashboard, query type: disk-allocated, format: influx-http
  use case: dashboard, query type: memory-utilization, format: influx-http
  use case: dashboard, query type: queue-bytes, format: influx-http
  use case: dashboard, query type: dashboard-all, format: influx-http
  use case: dashboard, query type: cpu-num, format: influx-http
  use case: dashboard, query type: disk-utilization, format: influx-http
  use case: dashboard, query type: nginx-requests, format: influx-http
  use case: dashboard, query type: throughput, format: influx-http
  use case: dashboard, query type: http-request-duration, format: influx-http
  use case: dashboard, query type: http-requests, format: influx-http
  use case: dashboard, query type: kapa-cpu, format: influx-http
  use case: dashboard, query type: redis-memory-utilization, format: influx-http
  use case: dashboard, query type: kapa-ram, format: influx-http
  use case: dashboard, query type: memory-total, format: influx-http
  use case: dashboard, query type: system-load, format: influx-http
  use case: dashboard, query type: availability, format: influx-http
  use case: dashboard, query type: cpu-utilization, format: influx-http
  use case: dashboard, query type: disk-usage, format: influx-http
  use case: dashboard, query type: kapa-load, format: influx-http       

```
#### query_benchmarker_tdengine用到的参数
```
Usage of ./query_benchmarker_tdengine:
  -batch-size int
        Number of queries in batch per worker for Dashboard use-case (default 18)
        说明：不用填，用缺省值即可
  -benchmark-duration duration
        Run querying continually for defined time interval, instead of stopping after all queries have been used
        说明：不用填，用缺省值即可
  -burn-in uint
        Number of queries to ignore before collecting statistics.
        说明：不用填，用缺省值即可
  -client-index int
        Index of a client host running this tool. Used to distribute load
        说明：不用填，用缺省值即可
  -debug int
        Whether to print debug messages.
        说明：不用填，用缺省值即可
  -dial-timeout duration
        TCP dial timeout. (default 15s)
        说明：不用填，用缺省值即可
  -file string
        Input file
        说明：不用填，用缺省值即可
  -grad-workers-inc
        Whether to gradually increase number of workers. The 'workers' params defines initial number of workers in this case.
        说明：不用填，用缺省值即可
  -grad-workers-max int
        Maximum number of workers when are added gradually. (default -1)
        说明：不用填，用缺省值即可
  -http-client-type string
        HTTP client type {fast, default} (default "fast")
        说明：不用填，用缺省值即可
  -increase-interval duration
        Interval when number of workers will increase (default 30s)
        说明：不用填，用缺省值即可
  -limit int
        Limit the number of queries to send. (default -1)
        说明：不用填，用缺省值即可
  -memprofile string
        Write a memory profile to this file.
        说明：不用填，用缺省值即可
  -moving-average-interval duration
        Interval of measuring mean response time on which moving average  is calculated. (default 30s)
        说明：不用填，用缺省值即可
  -notification-group string
        Terminate message notification siblings (comma-separated host:port list of other query benchmarkers)
        说明：不用填，用缺省值即可
  -notification-port int
        Listen port for remote notification messages. Used to remotely terminate benchmark (use -1 to disable it) (default -1)
        说明：不用填，用缺省值即可
  -notification-target string
        host:port of finish message notification receiver
        说明：不用填，用缺省值即可
  -print-interval uint
        Print timing stats to stderr after this many queries (0 to disable) (default 100)
        说明：每隔多少条查询后打印一下中间结果，建议填0或者一个比较大的数，避免打印较多的信息
  -print-responses
        Pretty print JSON response bodies (for correctness checking) (default false).
        说明：不用填，用缺省值即可
  -read-timeout duration
        TCP read timeout. (default 5m0s)
        说明：不用填，用缺省值即可
  -report-database string
        Database name where to store result metrics. (default "database_benchmarks")
        说明：不用填，用缺省值即可
  -report-host string
        Host to send result metrics.
        说明：不用填，用缺省值即可
  -report-password string
        User password for Host to send result metrics.
        说明：不用填，用缺省值即可
  -report-tags string
        Comma separated k:v tags to send  alongside result metrics.
        说明：不用填，用缺省值即可
  -report-telemetry
        Whether to report also progress info about mean, moving mean and #workers.
        说明：不用填，用缺省值即可
  -report-user string
        User for Host to send result metrics.
        说明：不用填，用缺省值即可
  -response-time-limit duration
        Query response time limit, after which will client stop.
        说明：不用填，用缺省值即可
  -rt-trend-samples int
        Number of avg response time samples used for linear regression (-1: number of samples equals increase-interval in seconds) (default -1)
        说明：不用填，用缺省值即可
  -telemetry-batch-size uint
        Telemetry batch size (lines). (default 1)
        说明：不用填，用缺省值即可
  -telemetry-stderr
        Whether to write telemetry also to stderr.
        说明：不用填，用缺省值即可
  -urls string
        Daemon URLs, comma-separated. Will be used in a round-robin fashion. (default "http://localhost:6020")
        说明：TDengine数据库的地址信息，以ip地址加tdengine的restful接口端口作为格式，比如TDengine运行22.23.123.32这个服务器上，则这里填写 http://22.23.123.32:6020
  -use-case string
        Enables use-case specific behavior. Empty for default behavior. Additional use-cases: dashboard
        说明：不用填，用缺省值即可
  -wait-interval duration
        Delay between sending batches of queries in the dashboard use-case
        说明：不用填，用缺省值即可
  -workers int
        Number of concurrent requests to make. (default 1)
        说明：查询线程的数量，这个参数决定了查询执行函数会用多少个线程并行的去查询
  -write-timeout duration
        TCP write timeout. (default 5m0s)
        说明：不用填，用缺省值即可
```
### 测试执行
基于以上对每个测试程序参数的解释，可以对照run.sh里的测试语句进行理解和尝试