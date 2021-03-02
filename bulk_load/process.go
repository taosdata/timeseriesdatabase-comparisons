package bulk_load

import (
	"github.com/taosdata/timeseriesdatabase-comparisons/util/report"
	"sync"
)

type BatchProcessor interface {
	PrepareProcess(i int)
	RunProcess(i int, waitGroup *sync.WaitGroup, telemetryPoints chan *report.Point, reportTags [][2]string) error
	AfterRunProcess(i int)
	EmptyBatchChanel()
}
