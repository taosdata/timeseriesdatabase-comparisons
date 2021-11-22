package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops512Hosts12hour produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops512Hosts12Hour10m struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops512Hosts12Hour10m(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops512Hosts12Hour10m{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops512Hosts12Hour10m) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage12HourByTenMinute512Hosts(q)
	return q
}
