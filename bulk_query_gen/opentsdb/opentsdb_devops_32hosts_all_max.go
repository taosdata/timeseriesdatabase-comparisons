package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops32HostsAll produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops32HostsAll struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops32HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops32HostsAll{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops32HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage32Hosts(q)
	return q
}
