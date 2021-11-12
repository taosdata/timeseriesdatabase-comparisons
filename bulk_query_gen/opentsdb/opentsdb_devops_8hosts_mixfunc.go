package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops8HostsMixfunc produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops8HostsMixfunc struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops8HostsMixfunc(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops8HostsMixfunc{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops8HostsMixfunc) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage8HostsMixfunc(q)
	return q
}
