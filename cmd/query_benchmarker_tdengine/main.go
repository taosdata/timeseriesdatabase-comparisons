// query_benchmarker speed tests InfluxDB using requests from stdin.
//
// It reads encoded Query objects from stdin, and makes concurrent requests
// to the provided HTTP endpoint. This program has no knowledge of the
// internals of the endpoint.
package main

import (
	"encoding/base64"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query"
	"github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query/http"
	"github.com/liu0x54/timeseriesdatabase-comparisons/util/report"
	//_ "github.com/taosdata/driver-go/taosSql"
)

/*
#cgo CFLAGS : -I/usr/include
#cgo LDFLAGS: -L/usr/lib -ltaos
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <taos.h>
*/
import "C"

// Program option vars:
type TDengineQueryBenchmarker struct {
	csvDaemonUrls string
	daemonUrls    []string

	dialTimeout    time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	httpClientType string
	clientIndex    int
	scanFinished   bool
	queryPool      sync.Pool
	queryChan      chan []*http.Query
}

var cgo int = 0
var querier = &TDengineQueryBenchmarker{}
var taosDriverName string = "taosSql"
var taosConns []unsafe.Pointer
var workers = 0

// Parse args:
func init() {

	bulk_query.Benchmarker.Init()
	querier.Init()

	flag.Parse()

	bulk_query.Benchmarker.Validate()
	querier.Validate()

}

func (b *TDengineQueryBenchmarker) Init() {
	flag.StringVar(&b.csvDaemonUrls, "urls", "http://localhost:6020", "Daemon URLs, comma-separated. Will be used in a round-robin fashion.")
	flag.DurationVar(&b.dialTimeout, "dial-timeout", time.Second*15, "TCP dial timeout.")
	flag.DurationVar(&b.readTimeout, "write-timeout", time.Second*300, "TCP write timeout.")
	flag.DurationVar(&b.writeTimeout, "read-timeout", time.Second*300, "TCP read timeout.")
	flag.StringVar(&b.httpClientType, "http-client-type", "fast", "HTTP client type {fast, default,cgo}")
	flag.IntVar(&b.clientIndex, "client-index", 0, "Index of a client host running this tool. Used to distribute load")
	flag.IntVar(&workers, "threads", 1, "cgo threads")
}

func (b *TDengineQueryBenchmarker) Validate() {
	b.daemonUrls = strings.Split(b.csvDaemonUrls, ",")
	if len(b.daemonUrls) == 0 {
		log.Fatal("missing 'urls' flag")
	}
	fmt.Printf("daemon URLs: %v\n", b.daemonUrls)

	if b.httpClientType == "fast" || b.httpClientType == "default" {
		fmt.Printf("Using HTTP client: %v\n", b.httpClientType)
		http.UseFastHttp = b.httpClientType == "fast"
	} else if b.httpClientType == "cgo" {
		fmt.Printf("Using TDengine C connector: %v\n", b.httpClientType)
		cgo = 1
	} else {
		log.Fatalf("Unsupported HTPP client type: %v", b.httpClientType)
	}
	if cgo == 1 {
		for i := 0; i < workers; i++ {
			taosConn, _ := taosConnect(b.csvDaemonUrls, "")
			taosConns = append(taosConns, taosConn)
		}
	}

}

func (b *TDengineQueryBenchmarker) Prepare() {
	// Make pools to minimize heap usage:
	b.queryPool = sync.Pool{
		New: func() interface{} {
			return &http.Query{
				HumanLabel:       make([]byte, 0, 1024),
				HumanDescription: make([]byte, 0, 1024),
				Method:           make([]byte, 0, 1024),
				Path:             make([]byte, 0, 1024),
				Body:             make([]byte, 0, 1024),
			}
		},
	}

	// Make data and control channels:
	b.queryChan = make(chan []*http.Query, 1000000)
}

func (b *TDengineQueryBenchmarker) GetProcessor() bulk_query.Processor {
	return b
}
func (b *TDengineQueryBenchmarker) GetScanner() bulk_query.Scanner {
	return b
}

func (b *TDengineQueryBenchmarker) PrepareProcess(i int) {
}

func (b *TDengineQueryBenchmarker) RunProcess(i int, workersGroup *sync.WaitGroup, statPool sync.Pool, statChan chan *bulk_query.Stat) {
	daemonUrl := b.daemonUrls[(i+b.clientIndex)%len(b.daemonUrls)]
	w := http.NewHTTPClient(daemonUrl, bulk_query.Benchmarker.Debug(), b.dialTimeout, b.readTimeout, b.writeTimeout)
	b.processQueries(w, workersGroup, statPool, statChan, i)
}

func (b *TDengineQueryBenchmarker) IsScanFinished() bool {
	return b.scanFinished
}

func (b *TDengineQueryBenchmarker) CleanUp() {
	close(b.queryChan)

}

func (b TDengineQueryBenchmarker) UpdateReport(params *report.QueryReportParams, reportTags [][2]string, extraVals []report.ExtraVal) (updatedTags [][2]string, updatedExtraVals []report.ExtraVal) {
	params.DBType = "TDengine"
	params.DestinationUrl = b.csvDaemonUrls
	updatedTags = reportTags
	updatedExtraVals = extraVals
	return
}

func main() {
	bulk_query.Benchmarker.RunBenchmark(querier)
}

var qind int64

// scan reads encoded Queries and places them onto the workqueue.
func (b *TDengineQueryBenchmarker) RunScan(r io.Reader, closeChan chan int) {
	dec := gob.NewDecoder(r)

	batch := make([]*http.Query, 0, bulk_query.Benchmarker.BatchSize())
	fmt.Printf("batch size %d, limit %d\n", bulk_query.Benchmarker.BatchSize(), bulk_query.Benchmarker.Limit())
	i := 0
loop:
	for {
		if bulk_query.Benchmarker.Limit() >= 0 && qind >= bulk_query.Benchmarker.Limit() {
			break
		}

		q := b.queryPool.Get().(*http.Query)
		err := dec.Decode(q)
		if err == io.EOF {
			fmt.Printf("io.Eof occurs --------------\n")
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		q.ID = qind
		batch = append(batch, q)
		i++
		if i == bulk_query.Benchmarker.BatchSize() {
			b.queryChan <- batch
			//batch = batch[:0]
			batch = nil
			batch = make([]*http.Query, 0, bulk_query.Benchmarker.BatchSize())
			i = 0
		}

		qind++
		select {
		case <-closeChan:
			log.Println("Received finish request")
			break loop
		default:
		}

	}
	b.scanFinished = true
}

// processQueries reads byte buffers from queryChan and writes them to the
// target server, while tracking latency.
func (b *TDengineQueryBenchmarker) processQueries(w http.HTTPClient, workersGroup *sync.WaitGroup, statPool sync.Pool, statChan chan *bulk_query.Stat, i int) error {
	restAuthorization := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("root:taosdata")))
	opts := &http.HTTPClientDoOptions{
		Authorization:        restAuthorization,
		Debug:                bulk_query.Benchmarker.Debug(),
		PrettyPrintResponses: bulk_query.Benchmarker.PrettyPrintResponses(),
	}
	var queriesSeen int64
	for queries := range b.queryChan {
		if len(queries) == 1 {
			if err := b.processSingleQuery(w, queries[0], opts, nil, nil, statPool, statChan, i); err != nil {
				log.Fatal(err)
			}
			queriesSeen++
		} else {
			var err error
			errors := 0
			done := 0
			errCh := make(chan error)
			doneCh := make(chan int, len(queries))
			//fmt.Printf("exec query %d\n", len(queries))
			for _, q := range queries {
				b.processSingleQuery(w, q, opts, errCh, doneCh, statPool, statChan, i)
				queriesSeen++
				if bulk_query.Benchmarker.GradualWorkersIncrease() {
					time.Sleep(time.Duration(rand.Int63n(150)) * time.Millisecond) // random sleep 0-150ms
				}
			}

		loop:
			for {
				select {
				case err = <-errCh:
					errors++
				case <-doneCh:
					done++
					if done == len(queries) {
						break loop
					}
				}
			}
			close(errCh)
			close(doneCh)
			if err != nil {
				log.Fatal(err)
			}
		}
		if bulk_query.Benchmarker.WaitInterval().Seconds() > 0 {
			time.Sleep(bulk_query.Benchmarker.WaitInterval())
		}
	}
	workersGroup.Done()
	if cgo == 1 {
		taosClose(taosConns[i])
	}
	return nil
}

func (b *TDengineQueryBenchmarker) processSingleQuery(w http.HTTPClient, q *http.Query, opts *http.HTTPClientDoOptions, errCh chan error, doneCh chan int, statPool sync.Pool, statChan chan *bulk_query.Stat, i int) error {
	defer func() {
		if doneCh != nil {
			doneCh <- 1
		}
	}()
	var lagMillis float64
	var err error
	if cgo == 1 {
		taosConn := taosConns[i]
		lagMillis, err = b.execSql(q, taosConn)
	} else {
		lagMillis, err = w.Do(q, opts)
	}
	stat := statPool.Get().(*bulk_query.Stat)
	stat.Init(q.HumanLabel, lagMillis)
	statChan <- stat
	b.queryPool.Put(q)
	if err != nil {
		qerr := fmt.Errorf("Error during request of query %s: %s\n", q.String(), err.Error())
		if errCh != nil {
			errCh <- qerr
			return nil
		} else {
			return qerr
		}
	}

	return nil
}

func (b *TDengineQueryBenchmarker) execSql(q *http.Query, taosConn unsafe.Pointer) (lag float64, err error) {

	sqlcmd := string(q.Body)
	start := time.Now()
	_, err = taosQuery(sqlcmd, taosConn)
	lag = float64(time.Since(start).Nanoseconds()) / 1e6 // milliseconds
	if err != nil {
		log.Fatalf("Query error: %s\n", err)
	}
	return lag, err
}

func taosConnect(ip, db string) (unsafe.Pointer, error) {
	user := "root"
	pass := "taosdata"
	port := 0
	cuser := C.CString(user)
	cpass := C.CString(pass)
	cip := C.CString(ip)
	cdb := C.CString(db)
	defer C.free(unsafe.Pointer(cip))
	defer C.free(unsafe.Pointer(cuser))
	defer C.free(unsafe.Pointer(cpass))
	defer C.free(unsafe.Pointer(cdb))

	taosObj := C.taos_connect(cip, cuser, cpass, cdb, (C.ushort)(port))
	if taosObj == nil {
		return nil, errors.New("taos_connect() fail!")
	}

	return (unsafe.Pointer)(taosObj), nil
}

func taosQuery(sqlstr string, taos unsafe.Pointer) (int, error) {
	csqlstr := C.CString(sqlstr)
	defer C.free(unsafe.Pointer(csqlstr))

	result := unsafe.Pointer(C.taos_query(taos, csqlstr))
	code := C.taos_errno(result)
	if 0 != code {

		errStr := C.GoString(C.taos_errstr(result))
		taosClose(taos)
		fmt.Println(errStr)
		return 0, errors.New(errStr)

	}

	// read result and save into mc struct
	//numfields := int(C.taos_field_count(result))
	for {
		res := C.taos_fetch_row(result)
		if res == C.TAOS_ROW(nil) {
			break
		}

	}

	C.taos_free_result(result)
	return 0, nil
}

func taosClose(taos unsafe.Pointer) {
	C.taos_close(taos)
}
