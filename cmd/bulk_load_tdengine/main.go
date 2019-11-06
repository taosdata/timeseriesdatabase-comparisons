// bulk_load_tdengine loads a TDengine daemon with data from stdin.
//
// The caller is responsible for assuring that the database is empty before
// bulk load.
package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/liu0x54/timeseriesdatabase-comparisons/bulk_data_gen/common"
	_ "github.com/taosdata/TDengine/src/connector/go/src/taosSql"

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
	doLoad         bool
	reportDatabase string
	reportHost     string
	reportUser     string
	reportPassword string
	reportTagsCSV  string
	useCase        string
)

// Global vars
var (
	bufPool    sync.Pool
	batchChans []chan string
	//batchChan      chan []string
	inputDone      chan struct{}
	workersGroup   sync.WaitGroup
	reportTags     [][2]string
	reportHostname string
	taosDriverName string = "taosSql"
)

// Parse args:
func init() {
	flag.StringVar(&daemonUrl, "url", "127.0.0.1:0", "TDengine URL.")

	flag.IntVar(&batchSize, "batch-size", 100, "Batch size (input items).")
	flag.IntVar(&workers, "workers", 2, "Number of parallel requests to make.")
	flag.StringVar(&useCase, "use-case", common.UseCaseChoices[0], "Use case to set specific load behavior. Options: "+strings.Join(common.UseCaseChoices, ","))

	flag.BoolVar(&doLoad, "do-load", true, "Whether to write data. Set this flag to false to check input read speed.")

	flag.StringVar(&reportDatabase, "report-database", "vehicle", "Database name where to store result metrics")
	flag.StringVar(&reportHost, "report-host", "", "Host to send result metrics")
	flag.StringVar(&reportUser, "report-user", "", "User for host to send result metrics")
	flag.StringVar(&reportPassword, "report-password", "", "User password for Host to send result metrics")
	flag.StringVar(&reportTagsCSV, "report-tags", "", "Comma separated k:v tags to send  alongside result metrics")

	flag.Parse()

	if reportHost != "" {
		fmt.Printf("results report destination: %v\n", reportHost)
		fmt.Printf("results report database: %v\n", reportDatabase)

		var err error
		reportHostname, err = os.Hostname()
		if err != nil {
			log.Fatalf("os.Hostname() error: %s", err.Error())
		}
		fmt.Printf("hostname for results report: %v\n", reportHostname)

		if reportTagsCSV != "" {
			pairs := strings.Split(reportTagsCSV, ",")
			for _, pair := range pairs {
				fields := strings.SplitN(pair, ":", 2)
				tagpair := [2]string{fields[0], fields[1]}
				reportTags = append(reportTags, tagpair)
			}
		}
		fmt.Printf("results report tags: %v\n", reportTags)
	}
}

func main() {
	bufPool = sync.Pool{
		New: func() interface{} {
			return make([]string, 0, batchSize)
		},
	}

	if doLoad {
		log.Println("Creating database ----")
		db, err := sql.Open(taosDriverName, "root:taosdata@/tcp("+daemonUrl+")/")
		if err != nil {
			log.Fatalf("Open database error: %s\n", err)
		}
		defer db.Close()

		//fmt.Println(db)
		createDatabase(db)

		for i := 0; i < workers; i++ {
			batchChans = append(batchChans, make(chan string, batchSize))
		}
		//batchChan = make(chan []string, workers)
		inputDone = make(chan struct{})
		log.Println("Starting workers ----")
		for i := 0; i < workers; i++ {
			workersGroup.Add(1)
			go processBatches(i)
		}

		start := time.Now()
		itemsRead, bytesRead, valuesRead := scan(db, batchSize)
		
		<-inputDone
		

		

		for i := 0; i < workers; i++ {
			close(batchChans[i])
		}
		//close(batchChan)
		workersGroup.Wait()
		end := time.Now()
		took := end.Sub(start)

		itemsRate := float64(itemsRead) / float64(took.Seconds())
		bytesRate := float64(bytesRead) / float64(took.Seconds())
		valuesRate := float64(valuesRead) / float64(took.Seconds())

		fmt.Printf("loaded %d items in %fsec with %d workers (mean point rate %.2f/s, mean value rate %.2f/s, %.2fMB/sec from stdin)\n", itemsRead, took.Seconds(), workers, itemsRate, valuesRate, bytesRate/(1<<20))
	}
}

func createDatabase(db *sql.DB) {
	sqlcmd := fmt.Sprintf("Drop database if exists %s", useCase)
	_, err := db.Exec(sqlcmd)
	sqlcmd = fmt.Sprintf("create database %s", useCase)
	_, err = db.Exec(sqlcmd)
	sqlcmd = fmt.Sprintf("use %s", useCase)
	_, err = db.Exec(sqlcmd)
	checkErr(err)

	return
}

func scan(db *sql.DB, itemsPerBatch int) (int64, int64, int64) {

	var vgid int
	var err error
	var itemsRead, bytesRead int64
	var totalPoints, totalValues int64


	//buff := bufPool.Get().([]string)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "create") {
			_, err = db.Exec(line)
		} else if strings.HasPrefix(line, "data"){
			totalPoints, totalValues, err = common.CheckTotalValues(line)
			if totalPoints > 0 || totalValues > 0 {
				continue
			}
			if err != nil {
				log.Fatal(err)
			}
		}else {
			itemsRead++
			bytesRead += int64(len(scanner.Bytes()))
			if !doLoad {
				continue
			}
			//buff = append(buff, line)
			/*
				n++
				if n >= itemsPerBatch {
					batchChan <- buff
					buff = bufPool.Get().([]string)
					buff = append(buff, "Insert into")
					n = 0
				}
			*/
			hun,_:= strconv.Atoi(string(line[0]))
			ten,_:= strconv.Atoi(string(line[1]))
			vgid, _ = strconv.Atoi(string(line[2]))
			vgid = hun*100+ten*10+vgid
			vgid = vgid % workers
			
			batchChans[vgid] <- line[3:]
			
		}

	}
	/*
		if n > 0 {
			batchChan <- buff
		}
	*/

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
	db, err := sql.Open(taosDriverName, "root:taosdata@/tcp("+daemonUrl+")/"+useCase)
	checkErr(err)
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
		sqlcmd[i] = onepoint
		i++
		if i > batchSize {
			i = 1
			_, err := db.Exec(strings.Join(sqlcmd, ""))
			if err != nil {
				log.Fatalf("Error writing: %s\n", strings.Join(sqlcmd, ""))//err.Error())
			}
		}
	}
	if i > 0 {
		i = 1
		_, err := db.Exec(strings.Join(sqlcmd, ""))
		if err != nil {
			log.Fatalf("Error writing: %s\n", strings.Join(sqlcmd, ""))//err.Error())
		}
	}

	workersGroup.Done()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
