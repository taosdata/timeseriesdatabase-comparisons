package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops8Hosts2Hr produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops8Hosts2Hr struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops8Hosts2Hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops8Hosts2Hr{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops8Hosts2Hr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteEightHostsTwoHr(q)
	return q
}
