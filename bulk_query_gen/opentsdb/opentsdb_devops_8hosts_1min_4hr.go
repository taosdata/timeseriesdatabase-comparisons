package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops8Hosts4Hr produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops8Hosts4Hr struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops8Hosts4Hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops8Hosts4Hr{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops8Hosts4Hr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteEightHostsFourHr(q)
	return q
}
