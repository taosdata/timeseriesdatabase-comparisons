package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops8Hosts12Hr produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops8Hosts12Hr struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops8Hosts12Hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops8Hosts12Hr{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops8Hosts12Hr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteEightHosts12Hr(q)
	return q
}
