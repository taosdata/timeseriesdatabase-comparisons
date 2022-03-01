# init path
scriptDir=$(dirname $(readlink -f $0))
workDir=${scriptDir=}/../../
echo "${workDir}"

#prepare to compile and compile program
cd ${workDir} && rm -rf go.mod go.sum 
go mod init github.com/taosdata/timeseriesdatabase-comparisons
go get github.com/golang/protobuf/proto
go get github.com/google/flatbuffers/go
go get github.com/pelletier/go-toml
go get github.com/pkg/profile
go get github.com/valyala/fasthttp

# generate program
mkdir -p build/tsdbcompare/bin
cd cmd/bulk_data_gen ;go build ;cp bulk_data_gen ../../build/tsdbcompare/bin
cd ../bulk_load_influx;go build ;cp bulk_load_influx ../../build/tsdbcompare/bin
cd ../bulk_load_tdengine;go build ; cp bulk_load_tdengine ../../build/tsdbcompare/bin

