package bulk_load

import (
	"github.com/liu0x54/timeseriesdatabase-comparisons/util/report"
	"sync"
)

type BatchProcessor interface {
	PrepareProcess(i int)
	RunProcess(i int, waitGroup *sync.WaitGroup, telemetryPoints chan *report.Point, reportTags [][2]string) error
	AfterRunProcess(i int)
	EmptyBatchChanel()
}
