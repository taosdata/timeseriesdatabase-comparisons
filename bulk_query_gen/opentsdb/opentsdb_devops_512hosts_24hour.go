package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops512Hosts12hour produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops512HostsAllByHr struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops512HostsAllBy1Hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops512HostsAllByHr{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops512HostsAllByHr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByHour512Hosts(q)
	return q
}
