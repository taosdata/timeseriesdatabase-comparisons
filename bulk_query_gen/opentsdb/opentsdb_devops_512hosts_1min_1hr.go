package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops512Hosts produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops512Hosts struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops512Hosts(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops512Hosts{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops512Hosts) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinute512Hosts(q)
	return q
}
