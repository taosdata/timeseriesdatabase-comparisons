package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops256Hosts produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops256Hosts struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops256Hosts(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops256Hosts{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops256Hosts) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinute256Hosts(q)
	return q
}
