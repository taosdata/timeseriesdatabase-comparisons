package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops16Hosts12hour produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops16HostsAllByHr struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops16HostsAllBy1Hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops16HostsAllByHr{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops16HostsAllByHr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByHour16Hosts(q)
	return q
}
