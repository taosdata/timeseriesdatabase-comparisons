#!/bin/bash

cpuNUm=`nproc`
address='bschang1'

echo "please make sure this ubuntu does not have java or golang installed. the program is going to install java 8 and golang 1.16.6"
echo "if you want to continue, press enter, or otherwise exit the program"
read input

address=$1

ssh-copy-id root@$address

## installing java 8
## install on client device
echo "changing java version on client"
sudo apt update
sudo apt -y autoremove openjdk-*-jre-headless
sudo apt -y install openjdk-8-jdk openjdk-8-jre openjdk-8-jre-headless
echo `java -version`

##install go 
echo "updating go on client"
echo `export GOPROXY=https://goproxy.io,direct`
echo `wget https://golang.org/dl/go1.16.6.linux-amd64.tar.gz`
echo `rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.6.linux-amd64.tar.gz`
echo `export PATH=$PATH:/usr/local/go/bin`
echo `go version`

#install go and java on server
echo "updating go on client"
ssh root@$address <<eeooff
#install java
sudo apt update
sudo apt -y autoremove openjdk-*-jre-headless
sudo apt -y install openjdk-8-jdk openjdk-8-jre openjdk-8-jre-headless

#install go
echo `export GOPROXY=https://goproxy.io,direct`
echo `wget https://golang.org/dl/go1.16.6.linux-amd64.tar.gz`
echo `rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.6.linux-amd64.tar.gz`
echo `export PATH=$PATH:/usr/local/go/bin`
exit
eeooff
echo "updating go on client finish"

# ##install Cassandra
# if Cassandra version is 3.11.x, then do not update Cassandra
CassandraVersion=`cat /etc/apt/sources.list.d/cassandra.sources.list | awk '{print($3)}'`
if [ "$CassandraVersion" != "311x" ]; then
    echo "the system does not have the proper Cassandra"
    echo `sudo apt -y autoremove cassandra`
    echo `rm /etc/apt/sources.list.d/cassandra.sources.list`
    echo "deb http://downloads.apache.org/cassandra/debian 311x main" | sudo tee -a /etc/apt/sources.list.d/cassandra.sources.list
    curl https://downloads.apache.org/cassandra/KEYS | sudo apt-key add -
    echo `sudo apt-get update`
    echo `sudo apt-get install cassandra`
    echo `systemctl stop cassandra`
    echo `sudo mkdir /data /data/Cassandra`
    echo `sudo chown -R Cassandra /data/Cassandra`
    echo "Cassandra Update finished"
fi

CassandraVersion = `ssh root@$address cat /etc/apt/sources.list.d/cassandra.sources.list | awk '{print($3)}'`
if [ "$CassandraVersion" != "311x" ]; then
    ssh root@$address << eeooff
    echo "the system does not have the proper Cassandra"
    echo `sudo apt -y autoremove cassandra`
    echo `rm /etc/apt/sources.list.d/cassandra.sources.list`
    echo "deb http://downloads.apache.org/cassandra/debian 311x main" | sudo tee -a /etc/apt/sources.list.d/cassandra.sources.list
    curl https://downloads.apache.org/cassandra/KEYS | sudo apt-key add -
    echo `sudo apt-get update`
    echo `sudo apt-get install cassandra`
    echo `systemctl stop cassandra`
    echo `sudo mkdir /data /data/Cassandra`
    echo `sudo chown -R Cassandra /data/Cassandra`
    echo "Cassandra Update finished"
    echo `sudo apt-get install -y gcc cmake build-essential git autoconf`
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
else 
    ssh root@$address << eeooff
    echo `sudo apt-get install -y gcc cmake build-essential git autoconf`
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
fi

## install latest version of TDengine Community Version Master Branch
echo "installing newest tdengine"
echo `sudo apt-get install -y gcc cmake build-essential git autoconf`
echo `rm -rf TDengine`
git clone https://github.com/taosdata/TDengine.git && cd TDengine
git checkout master && git fetch && git pull
git submodule update --init --recursive
mkdir debug && cd debug
cmake .. -DJEMALLOC_ENABLED=true && cmake --build .
sudo make install
cd ../..
sudo mkdir /data /data/taos /data/taos/data


echo "modifying Cassandra's config"
## the following are for modifying the Cassandra configs
echo `sed -i "s/concurrent_reads:.*/concurrent_reads: `nproc` /g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/batch_size_fail_threshold_in_kb.*/batch_size_fail_threshold_in_kb: 1024/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/concurrent_writes:.*/concurrent_writes: 200/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/concurrent_counter_writes:.*/concurrent_counter_writes: 200/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/read_request_timeout_in_ms:.*/read_request_timeout_in_ms: 500000 /g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/range_request_timeout_in_ms:.*/range_request_timeout_in_ms: 1000000/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/write_request_timeout_in_ms:.*/write_request_timeout_in_ms: 200000/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/counter_write_request_timeout_in_ms:.*/counter_write_request_timeout_in_ms: 500000/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/cas_contention_timeout_in_ms:.*/cas_contention_timeout_in_ms: 100000/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/truncate_request_timeout_in_ms:.*/truncate_request_timeout_in_ms: 6000000/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/request_timeout_in_ms:.*/request_timeout_in_ms: 1000000/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/slow_query_log_timeout_in_ms:.*/slow_query_log_timeout_in_ms: 50000/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/data_file_directories.*/slow_query_log_timeout_in_ms: 50000/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/rpc_address:.*/rpc_address: $address/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/\/var\/lib\/cassandra\/data/\/data\/Cassandra/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/.*trickle_fsync.*/trickle_fsync: false/g" /etc/cassandra/cassandra.yaml`
echo `sed -i "s/.*trickle_fsync_interval_in_kb.*/trickle_fsync_interval_in_kb: 0/g" /etc/cassandra/cassandra.yaml`
scp /etc/cassandra/cassandra.yaml $address:/etc/cassandra
ssh root@$address sed -i "s/concurrent_reads:.*/concurrent_reads: `nproc` /g" /etc/cassandra/cassandra.yaml
echo "Cassandra config modification finished"

echo "modifying TDengine's config"
## the following are for modifying TDengine's configs
sed -i "s/.*dataDir.*/dataDir \/data\/lib\/taos /g" /etc/taos/taos.cfg
sed -i "s/.*walLevel.*/walLevel 1 /g" /etc/taos/taos.cfg
sed -i "s/.*fsync.*/fsync 0 /g" /etc/taos/taos.cfg
sed -i "s/.*maxSQLLength.*/maxSQLLength 1048576 /g" /etc/taos/taos.cfg
sed -i "s/.*asyncLog.*/asyncLog 1 /g" /etc/taos/taos.cfg
sed -i "s/.*debugFlag.*/debugFlag 131 /g" /etc/taos/taos.cfg
scp /etc/taos/taos.cfg $address:/etc/taos
echo "modifying TDengine's config finished"

echo "start to build the programs for comparsion"
cd ../..
rm -f go.mod go.sum
go mod init github.com/taosdata/timeseriesdatabase-comparisons
go mod tidy

mkdir -p build/tsdbcompare/bin
cd cmd/bulk_data_gen && go build 
cp bulk_data_gen ../../build/tsdbcompare/bin/
cd ../bulk_load_cassandra && go build 
cp bulk_load_cassandra ../../build/tsdbcompare/bin/
cd ../bulk_load_tdengine && go build 
cp bulk_load_tdengine ../../build/tsdbcompare/bin/
cd ../bulk_query_gen && go build 
cp bulk_query_gen ../../build/tsdbcompare/bin/
cd ../query_benchmarker_cassandra && go build 
cp query_benchmarker_cassandra ../../build/tsdbcompare/bin/
cd ../query_benchmarker_tdengine && go build 
cp query_benchmarker_tdengine ../../build/tsdbcompare/bin/

cd ../../build/tsdbcompare/

echo "building environment finish. Now the comparsion test is ready to launch"
