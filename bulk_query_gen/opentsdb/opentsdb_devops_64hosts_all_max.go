package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops64HostsAll produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops64HostsAll struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops64HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops64HostsAll{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops64HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage64Hosts(q)
	return q
}
