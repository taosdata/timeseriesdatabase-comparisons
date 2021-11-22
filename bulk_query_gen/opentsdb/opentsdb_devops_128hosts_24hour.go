package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops128Hosts12hour produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops128HostsAllByHr struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops128HostsAllBy1Hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops128HostsAllByHr{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops128HostsAllByHr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByHour128Hosts(q)
	return q
}
