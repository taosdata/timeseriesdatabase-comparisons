package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops128HostsAll produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops128HostsAll struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops128HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops128HostsAll{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops128HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage128Hosts(q)
	return q
}
