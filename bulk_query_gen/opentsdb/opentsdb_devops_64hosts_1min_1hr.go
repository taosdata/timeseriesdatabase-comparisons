package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops64Hosts produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops64Hosts struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops64Hosts(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops64Hosts{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops64Hosts) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinute64Hosts(q)
	return q
}
