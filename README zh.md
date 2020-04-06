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
如果遇到这种情况，就需要从本项目的源代码编译出二进制代码后再运行本测试程序。下面介绍编译的步骤

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

编译出二进制代码后，就可以参考run.sh脚本中的命令顺序，手动执行测试了。