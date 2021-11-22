package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops1Hosts12hour produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops1HostsAllByHr struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops1HostsAllBy1Hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops1HostsAllByHr{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops1HostsAllByHr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByHour1Hosts(q)
	return q
}
