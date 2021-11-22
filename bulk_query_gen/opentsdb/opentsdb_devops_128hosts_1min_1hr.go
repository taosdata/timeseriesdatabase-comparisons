package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops128Hosts produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops128Hosts struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops128Hosts(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops128Hosts{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops128Hosts) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinute128Hosts(q)
	return q
}
