package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops32Hosts produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops32Hosts struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops32Hosts(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops32Hosts{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops32Hosts) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteThirtyTwoHosts(q)
	return q
}
