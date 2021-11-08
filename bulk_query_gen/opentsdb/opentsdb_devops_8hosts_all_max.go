package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops8Hosts12hour produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops8HostsAll struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops8HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops8HostsAll{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops8HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage8Hosts(q)
	return q
}
