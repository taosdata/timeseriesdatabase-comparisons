#!/bin/bash
address='localhost'
address=$1

ssh-copy-id root@$address

echo "installing go1.16"
echo `export GOPROXY=https://goproxy.io,direct`
echo `wget https://golang.org/dl/go1.16.6.linux-amd64.tar.gz`
echo `rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.6.linux-amd64.tar.gz`
echo `export PATH=$PATH:/usr/local/go/bin`

ssh root@$address <<eeooff
#install go
echo `export GOPROXY=https://goproxy.io,direct`
echo `wget https://golang.org/dl/go1.16.6.linux-amd64.tar.gz`
echo `rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.6.linux-amd64.tar.gz`
echo `export PATH=$PATH:/usr/local/go/bin`
exit
eeooff
echo "go 1.16 installing finished"

echo "install influxDB and TDengine"
wget https://dl.influxdata.com/influxdb/releases/influxdb_1.8.7_amd64.deb
dpkg -r influxdb
dpkg -r influxdb2
dpkg -i influxdb_1.8.7_amd64.deb
echo `rm -rf TDengine`
git clone https://github.com/taosdata/TDengine.git && cd TDengine
git checkout master && git fetch && git pull
git submodule update --init --recursive
mkdir debug && cd debug
cmake .. -DJEMALLOC_ENABLED=true && cmake --build .
sudo make install
cd ../..
sudo mkdir /data /data/taos /data/taos/data

ssh root@$address <<eeooff
wget https://dl.influxdata.com/influxdb/releases/influxdb_1.8.7_amd64.deb
dpkg -r influxdb
dpkg -r influxdb2
dpkg -i influxdb_1.8.7_amd64.deb
echo `rm -rf TDengine`
git clone https://github.com/taosdata/TDengine.git && cd TDengine
git checkout master && git fetch && git pull
git submodule update --init --recursive
mkdir debug && cd debug
cmake .. -DJEMALLOC_ENABLED=true && cmake --build .
sudo make install
cd ../..
sudo mkdir /data /data/taos /data/taos/data
exit
eeooff

echo "modifying TDengine's config"
sed -i "s/.*walLevel.*/walLevel 2 /g" /etc/taos/taos.cfg
sed -i "s/.*fsync.*/fsync 0 /g" /etc/taos/taos.cfg
scp /etc/taos/taos.cfg $address:/etc/taos

echo "start to build the programs for comparsion"
cd ../..
rm -f go.mod go.sum
go mod init github.com/taosdata/timeseriesdatabase-comparisons
go mod tidy

mkdir -p build/tsdbcompare/bin
cd cmd/bulk_data_gen && go build 
cp bulk_data_gen ../../build/tsdbcompare/bin/
cd ../bulk_load_influx  && go build 
cp bulk_load_influx  ../../build/tsdbcompare/bin/
cd ../bulk_load_tdengine && go build 
cp bulk_load_tdengine ../../build/tsdbcompare/bin/
cd ../bulk_query_gen && go build 
cp bulk_query_gen ../../build/tsdbcompare/bin/
cd ../query_benchmarker_influxdb  && go build 
cp query_benchmarker_influxdb  ../../build/tsdbcompare/bin/
cd ../query_benchmarker_tdengine && go build 
cp query_benchmarker_tdengine ../../build/tsdbcompare/bin/

cd ../../build/tsdbcompare/
rm go1.16.6*
rm influxdb_1.8.7*
