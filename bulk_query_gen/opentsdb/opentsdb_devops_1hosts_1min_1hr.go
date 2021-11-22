package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops1Hosts produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops1Hosts struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops1Hosts(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops1Hosts{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops1Hosts) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteOneHost(q)
	return q
}
