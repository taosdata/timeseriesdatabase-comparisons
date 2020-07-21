// bulk_load_tdengine loads a TDengine daemon with data from stdin.
//
// The caller is responsible for assuring that the database is empty before
// bulk load.
package main

/*
#cgo CFLAGS : -I/usr/include
#cgo LDFLAGS: -L/usr/lib -ltaos
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <taos.h>
*/
import "C"

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"sync"
	"time"
	"unsafe"

	"github.com/liu0x54/timeseriesdatabase-comparisons/bulk_data_gen/common"
	//_ "github.com/taosdata/driver-go/taosSql"

	//	"github.com/caict-benchmark/BDC-TS/util/report"
	"strconv"
	"strings"
	//	"bytes"
)

// Program option vars:
var (
	daemonUrl      string
	workers        int
	batchSize      int
	slaveSource    bool
	doLoad         bool
	reportDatabase string
	reportHost     string
	reportUser     string
	reportPassword string
	reportTagsCSV  string
	useCase        string
	loadfile       string
	fileoutput     bool
	httpAPI        bool
)

// Global vars
var (
	bufPool     sync.Pool
	batchChans  []chan string
	sqlCmdChans []chan string
	statistics  []int64
	//batchChan      chan []string
	inputDone      chan struct{}
	workersGroup   sync.WaitGroup
	reportTags     [][2]string
	reportHostname string
	taosDriverName string = "taosSql"
	tablesqlname   string = "data/tables.sql"
	tablesqlfile   *os.File
)

// Parse args:
func init() {
	flag.StringVar(&daemonUrl, "url", "localhost", "TDengine URL.")

	flag.IntVar(&batchSize, "batch-size", 100, "Batch size (input items).")
	flag.IntVar(&workers, "workers", 2, "Number of parallel requests to make.")
	flag.StringVar(&useCase, "use-case", common.UseCaseChoices[0], "Use case to set specific load behavior. Options: "+strings.Join(common.UseCaseChoices, ","))

	flag.BoolVar(&doLoad, "do-load", false, "Whether to read data from file or from stand input. Set this flag to true to get input from file.")

	flag.StringVar(&reportDatabase, "report-database", "vehicle", "Database name where to store result metrics")
	flag.StringVar(&reportHost, "report-host", "", "Host to send result metrics")
	flag.StringVar(&reportUser, "report-user", "", "User for host to send result metrics")
	flag.StringVar(&reportPassword, "report-password", "", "User password for Host to send result metrics")
	flag.StringVar(&reportTagsCSV, "report-tags", "node1", "Comma separated k:v tags to send  alongside result metrics")
	flag.BoolVar(&slaveSource, "slavesource", false, "if slave source, will not create database")
	flag.StringVar(&loadfile, "file", "", "Input file")
	flag.BoolVar(&fileoutput, "fileout", false, "if file out, will out put sql into file")
	flag.BoolVar(&httpAPI, "http-api", false, "if true, will out put sql through TDengine RestAPI")

	flag.Parse()

	//if reportHost != "" {
	fmt.Printf("results report destination: %v\n", reportHost)
	fmt.Printf("results report database: %v\n", reportDatabase)

	//		var err error
	//		reportHostname, err = os.Hostname()
	//		if err != nil {
	//			log.Fatalf("os.Hostname() error: %s", err.Error())
	//		}
	reportHostname = "TDengine"
	fmt.Printf("hostname for results report: %v\n", reportHostname)

	if reportTagsCSV != "" {
		/*pairs := strings.Split(reportTagsCSV, ",")
		for _, pair := range pairs {
			fields := strings.SplitN(pair, ":", 2)
			tagpair := [2]string{fields[0], fields[1]}
			reportTags = append(reportTags, tagpair)
		}*/
	}
	fmt.Printf("results report tags: %v\n", reportTagsCSV)
	cmd := exec.Command("rm", "data/*")
	cmd.Run()
	//}
	createtablesql, err := os.OpenFile(tablesqlname, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tablesqlfile = createtablesql
}

func main() {
	var start time.Time
	bufPool = sync.Pool{
		New: func() interface{} {
			return make([]string, 0, batchSize)
		},
	}
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()
	log.Println("Creating database ----")
	//db, err := sql.Open(taosDriverName, "root:taosdata@/tcp("+daemonUrl+":0)/")
	//if err != nil {
	//	log.Fatalf("Open database error: %s\n", err)
	//}
	//defer db.Close()

	if !slaveSource {
		createDatabase()
	}

	for i := 0; i < workers; i++ {
		batchChans = append(batchChans, make(chan string, 100*batchSize))
	}
	//batchChan = make(chan []string, workers)
	inputDone = make(chan struct{})
	log.Println("Starting workers ----")
	statistics = make([]int64, workers)
	for i := 0; i < workers; i++ {
		workersGroup.Add(1)
		go processBatches(i)
	}

	if fileoutput != true {
		start = time.Now()
	}
	//start := time.Now()
	itemsRead, bytesRead, valuesRead := scan(batchSize)

	<-inputDone
	//end := time.Now()

	for i := 0; i < workers; i++ {
		close(batchChans[i])
		fmt.Printf("worker %d, processed %d data \n", i, statistics[i])
	}
	//close(batchChan)
	workersGroup.Wait()
	if fileoutput == true {
		start = time.Now()
		wt := fmt.Sprintf("-T %d", workers)
		cmd := exec.Command("taos", "-D", "data", wt)
		cmd.Run()
	}

	end := time.Now()
	took := end.Sub(start)

	itemsRate := float64(itemsRead) / float64(took.Seconds())
	bytesRate := float64(bytesRead) / float64(took.Seconds())
	valuesRate := float64(valuesRead) / float64(took.Seconds())

	fmt.Printf("loaded %d items in %fsec with %d workers (mean point rate %.2f/s, mean value rate %.2f/s, %.2fMB/sec from stdin)\n", itemsRead, took.Seconds(), workers, itemsRate, valuesRate, bytesRate/(1<<20))
	//	sqlcmd := fmt.Sprintf("create database if not exists benchmarkreport")
	//	_, err = db.Exec(sqlcmd)
	//	_, err = db.Exec("use benchmarkreport")
	//	sqlcmd = fmt.Sprintf("create table if not exists bmreport (ts timestamp, starttime binary(50), endtime binary(50),itemsread double, bytesread double,valuesread double,timetook double,recordsrate double, bytesrate double,valuesrate double, workers double, batchsize double, usecase binary(50)) tags(host binary(50), proc_id binary(40))")
	//	_, err = db.Exec(sqlcmd)
	//	sqlcmd = fmt.Sprintf("insert into %s_%s using bmreport tags(\"%s\",\"%s\") values(0, \"%s\",\"%s\",%d,%d,%d,%f,%f,%f,%f,%d,%d,\"%s\")", reportHostname, reportTagsCSV, reportHostname, reportTagsCSV, start.Format(time.RFC3339), end.Format(time.RFC3339), itemsRead, bytesRead, valuesRead, took.Seconds(), itemsRate, bytesRate/(1<<20), valuesRate, workers, batchSize, useCase)
	//	_, err = db.Exec(sqlcmd)
	//	checkErr(err)
	tablesqlfile.Close()
}

func createDatabase() {
	if fileoutput == true {
		sqlcmd := fmt.Sprintf("Drop database if exists %s;\n", useCase)
		tablesqlfile.WriteString(sqlcmd)
		sqlcmd = fmt.Sprintf("create database %s; \n", useCase)
		tablesqlfile.WriteString(sqlcmd)
		sqlcmd = fmt.Sprintf("use %s;\n", useCase)
		tablesqlfile.WriteString(sqlcmd)
		return
	} else if httpAPI == true {
		client := new(http.Client)
		sqlcmd := fmt.Sprintf("Drop database if exists %s", useCase)
		httpExecSQL(sqlcmd, client)
		sqlcmd = fmt.Sprintf("create database %s ", useCase)
		httpExecSQL(sqlcmd, client)
		sqlcmd = fmt.Sprintf("use %s", useCase)
		httpExecSQL(sqlcmd, client)
		return
	}
	taosConn, err := taosConnect(daemonUrl, "")
	checkErr(err)
	sqlcmd := fmt.Sprintf("Drop database if exists %s", useCase)
	_, err = taosQuery(sqlcmd, taosConn)
	checkErr(err)
	sqlcmd = fmt.Sprintf("create database %s ", useCase)
	_, err = taosQuery(sqlcmd, taosConn)
	checkErr(err)

	return
}

func scan(itemsPerBatch int) (int64, int64, int64) {

	var vgid int
	var err error
	var itemsRead, bytesRead int64
	var totalPoints, totalValues int64
	var sourceReader *os.File
	var taosConn unsafe.Pointer

	if fileoutput == false && httpAPI != true {
		taosConn, err = taosConnect(daemonUrl, useCase)
		checkErr(err)
		defer taosClose(taosConn)
	}

	//buff := bufPool.Get().([]string)

	if loadfile != "" && doLoad {
		if f, err := os.Open(loadfile); err == nil {
			sourceReader = f
		} else {
			log.Fatalf("Error opening %s: %v\n", loadfile, err)
		}
	} else {
		sourceReader = os.Stdin
	}

	scanner := bufio.NewScanner(sourceReader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line[7:], "create") {

			hscode, err := strconv.ParseInt(strings.TrimSpace(line[0:6]), 10, 64)
			if err != nil {
				fmt.Println(err)
			}

			vgid = int(hscode) % workers

			//fmt.Printf("string %s, hascode : %d, vgid %d\n", line[0:6], hscode, vgid)
			batchChans[vgid] <- line[7:]
			statistics[vgid]++

			/*go func() {
				_, err = db.Exec(line[7:])
				if err != nil {
					fmt.Println(err)
				}
			}()*/

		} else if strings.HasPrefix(line, "create") {
			if fileoutput == true && httpAPI != true {
				tablesqlfile.WriteString(line + "\n")
			} else if httpAPI != true {
				_, err = taosQuery(line, taosConn)
				checkErr(err)
			} else {
				client := new(http.Client)
				httpExecSQL(line, client)
			}

		} else if strings.HasPrefix(line, "data") {
			totalPoints, totalValues, err = common.CheckTotalValues(line)
			if totalPoints > 0 || totalValues > 0 {
				continue
			}
			if err != nil {
				log.Fatal(err)
			}
		} else {
			itemsRead++
			bytesRead += int64(len(scanner.Bytes())) - 6
			if !doLoad {
				continue
			}

			hscode, _ := strconv.ParseInt(strings.TrimSpace(line[0:6]), 10, 64)

			vgid = int(hscode) % workers

			batchChans[vgid] <- line[6:]
			statistics[vgid]++

		}

	}

	// Closing inputDone signals to the application that we've read everything and can now shut down.
	close(inputDone)

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %s", err.Error())
	}
	if itemsRead != totalPoints {
		log.Fatalf("Incorrent number of read items: %d, expected: %d:", itemsRead, totalPoints)
	}
	return itemsRead, bytesRead, totalValues
}

func processBatches(iworker int) {
	var i int
	var err error
	var datafile *os.File
	var client *http.Client
	var taosConn unsafe.Pointer

	if fileoutput != true && httpAPI != true {
		taosConn, err = taosConnect(daemonUrl, useCase)
		checkErr(err)
		defer taosClose(taosConn)
	} else if fileoutput == true {
		dfn := fmt.Sprintf("data/%d.sql", iworker)
		datafile, _ = os.OpenFile(dfn, os.O_CREATE|os.O_WRONLY, 0644)
		usedb := fmt.Sprintf("use %s;\n", useCase)
		datafile.WriteString(usedb)

	} else {
		client = new(http.Client)
	}
	sqlcmd := make([]string, batchSize+1)
	i = 0
	sqlcmd[i] = "Insert into"
	i++
	/*
		for batch := range batchChan {
			if !doLoad {
				continue
			}
			// Write the batch.
			_, err := db.Exec(strings.Join(batch, ""))
			if err != nil {
				log.Fatalf("Error writing: %s\n", err.Error())
			}
		}
	*/
	for onepoint := range batchChans[iworker] {
		if strings.HasPrefix(onepoint, "create") {
			if fileoutput != true && httpAPI != true {
				_, err := taosQuery(onepoint+";", taosConn)
				if err != nil {
					log.Fatalf("Error create table: %s; error:%s\n", onepoint, err) //err.Error())
				}
			} else if fileoutput == true {
				datafile.WriteString(onepoint + ";\n")
			} else {
				httpExecSQL(onepoint, client)
			}
		} else {
			sqlcmd[i] = onepoint
			i++
			if i > batchSize {
				i = 1
				if fileoutput != true && httpAPI != true {
					_, err = taosQuery(strings.Join(sqlcmd, ""), taosConn)
					checkErr(err)
				} else if fileoutput == true {
					datafile.WriteString(strings.Join(sqlcmd, "") + ";\n")
				} else {
					httpExecSQL(strings.Join(sqlcmd, ""), client)
				}

				if err != nil {
					log.Fatalf("Error writing: %s,  error:%s\n", strings.Join(sqlcmd, ""), err) //err.Error())
				}
			}

		}

	}
	if i > 1 {
		i = 1

		if fileoutput != true && httpAPI != true {
			_, err = taosQuery(strings.Join(sqlcmd, ""), taosConn)
			checkErr(err)
		} else if fileoutput == true {
			datafile.WriteString(strings.Join(sqlcmd, "") + ";\n")
		} else {
			httpExecSQL(strings.Join(sqlcmd, ""), client)
		}

	}
	datafile.Close()

	workersGroup.Done()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func httpExecSQL(sqlcmd string, client *http.Client) error {
	body := strings.NewReader(sqlcmd)
	req, _ := http.NewRequest("GET", "http://"+daemonUrl+":6020/rest/sql", body)
	req.SetBasicAuth("root", "taosdata")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	return nil
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
		return 0, errors.New(errStr)

	}

	// read result and save into mc struct
	numfields := int(C.taos_field_count(result))
	C.taos_free_result(result)
	return numfields, nil
}

func taosClose(taos unsafe.Pointer) {
	C.taos_close(taos)
}
