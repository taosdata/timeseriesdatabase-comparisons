package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops256Hosts12hour produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops256HostsAllByHr struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops256HostsAllBy1Hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops256HostsAllByHr{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops256HostsAllByHr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByHour256Hosts(q)
	return q
}
