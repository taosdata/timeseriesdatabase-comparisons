// bulk_load_opentsdb loads an OpenTSDB daemon with data from stdin.
//
// The caller is responsible for assuring that the database is empty before
// bulk load.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/taosdata/timeseriesdatabase-comparisons/bulk_load"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/klauspost/compress/gzip"
	"github.com/taosdata/timeseriesdatabase-comparisons/util/report"
)

type OpenTsdbBulkLoad struct {
	// Program option vars:
	csvDaemonUrls string
	daemonUrls    []string

	// Global vars
	bufPool      sync.Pool
	batchChan    chan *bytes.Buffer
	inputDone    chan struct{}
	valuesRead   int64
	itemsRead    int64
	bytesRead    int64
	scanFinished bool
}

var load = &OpenTsdbBulkLoad{}

// Parse args:
func init() {
	bulk_load.Runner.Init(5000)
	load.Init()

	flag.Parse()

	bulk_load.Runner.Validate()
	load.Validate()

}

func main() {
	bulk_load.Runner.Run(load)
}

func (l *OpenTsdbBulkLoad) Init() {
	flag.StringVar(&l.csvDaemonUrls, "urls", "http://localhost:8086", "OpenTSDB URLs, comma-separated. Will be used in a round-robin fashion.")
}

func (l *OpenTsdbBulkLoad) Validate() {
	l.daemonUrls = strings.Split(l.csvDaemonUrls, ",")
	if len(l.daemonUrls) == 0 {
		log.Fatal("missing 'urls' flag")
	}
	fmt.Printf("daemon URLs: %v\n", l.daemonUrls)
}

func (l *OpenTsdbBulkLoad) CreateDb() {
	//不需要创建db
	return
}

func (l *OpenTsdbBulkLoad) PrepareWorkers() {
	l.bufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 4*1024*1024)) //4M
		},
	}

	l.batchChan = make(chan *bytes.Buffer, bulk_load.Runner.Workers)
	l.inputDone = make(chan struct{})

}

func (l *OpenTsdbBulkLoad) GetBatchProcessor() bulk_load.BatchProcessor {
	return l
}

func (l *OpenTsdbBulkLoad) GetScanner() bulk_load.Scanner {
	return l
}

func (l *OpenTsdbBulkLoad) SyncEnd() {
	<-l.inputDone
	close(l.batchChan)
}

func (l *OpenTsdbBulkLoad) CleanUp() {

}

func (l *OpenTsdbBulkLoad) UpdateReport(params *report.LoadReportParams) (reportTags [][2]string, extraVals []report.ExtraVal) {
	params.DBType = "OpenTSDB"
	params.DestinationUrl = l.daemonUrls[0]
	params.IsGzip = true
	return
}

func (l *OpenTsdbBulkLoad) PrepareProcess(i int) {
	//开始测试前的准备
}

func (l *OpenTsdbBulkLoad) RunProcess(i int, waitGroup *sync.WaitGroup, telemetryPoints chan *report.Point, reportTags [][2]string) error {
	//开始任务
	daemonUrl := l.daemonUrls[i%len(l.daemonUrls)]
	cfg := HTTPWriterConfig{
		Host: daemonUrl,
	}
	return l.processBatches(NewHTTPWriter(cfg), waitGroup)
}

func (l *OpenTsdbBulkLoad) AfterRunProcess(i int) {

}

func (l *OpenTsdbBulkLoad) EmptyBatchChanel() {
	for range l.batchChan {
		//read out remaining batches
	}
}

func (l *OpenTsdbBulkLoad) IsScanFinished() bool {
	return l.scanFinished
}

func (l *OpenTsdbBulkLoad) GetReadStatistics() (itemsRead, bytesRead, valuesRead int64) {
	itemsRead = l.itemsRead
	bytesRead = l.bytesRead
	valuesRead = l.valuesRead
	return
}

// scan reads one line at a time from stdin.
// When the requested number of lines per batch is met, send a batch over batchChan for the workers to write.
func (l *OpenTsdbBulkLoad) RunScanner(r io.Reader, syncChanDone chan int) {
	l.scanFinished = false
	l.itemsRead = 0
	l.bytesRead = 0
	l.valuesRead = 0
	buf := l.bufPool.Get().(*bytes.Buffer)
	zw := gzip.NewWriter(buf)

	openbracket := []byte("[")
	closebracket := []byte("]")
	commaspace := []byte(", ")
	newline := []byte("\n")

	zw.Write(openbracket)
	zw.Write(newline)

	reader := bufio.NewReaderSize(r, 4*1024*1024)
	var deadline time.Time
	if bulk_load.Runner.TimeLimit > 0 {
		deadline = time.Now().Add(bulk_load.Runner.TimeLimit)
	}
	var n = 0
	needComma := false
	for {
		if needComma {
			zw.Write(commaspace)
			zw.Write(newline)
		}
		strBytes, hasMore, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("Error reading input: %s", err.Error())
			}
		}
		zw.Write(strBytes)
		if !hasMore {
			if n >= bulk_load.Runner.BatchSize {
				zw.Write(newline)
				zw.Write(closebracket)
				zw.Close()

				l.batchChan <- buf
				buf = l.bufPool.Get().(*bytes.Buffer)
				zw = gzip.NewWriter(buf)
				zw.Write(openbracket)
				zw.Write(newline)
				n = 0
				if bulk_load.Runner.TimeLimit > 0 && time.Now().After(deadline) {
					bulk_load.Runner.SetPrematureEnd("Timeout elapsed")
					break
				}
				needComma = false
			} else {
				n += 1
				needComma = true
			}
		} else {
			needComma = false
		}
		select {
		case <-syncChanDone:
			break
		default:
		}
	}
	// Finished reading input, make sure last batch goes out.
	if n > 0 {
		zw.Write(newline)
		zw.Write(closebracket)
		zw.Close()
		l.batchChan <- buf
	}

	// Closing inputDone signals to the application that we've read everything and can now shut down.
	close(l.inputDone)
	l.scanFinished = true
}

// processBatches reads byte buffers from batchChan and writes them to the target server, while tracking stats on the write.
func (l *OpenTsdbBulkLoad) processBatches(w LineProtocolWriter, workersGroup *sync.WaitGroup) error {
	var returnErr error
	for batch := range l.batchChan {
		// Write the batch: try until backoff is not needed.
		if bulk_load.Runner.DoLoad {
			var err error
			for {
				_, err = w.WriteLineProtocol(batch.Bytes())
				if err != nil {
					break
				}
			}
			if err != nil {
				returnErr = fmt.Errorf("Error writing: %s\n", err.Error())
			}
		}
		// Return the batch buffer to the pool.
		batch.Reset()
		l.bufPool.Put(batch)
		time.Sleep(time.Second)
	}
	workersGroup.Done()
	return returnErr
}
