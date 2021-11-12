package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops512HostsAll produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops512HostsAll struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops512HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops512HostsAll{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops512HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage512Hosts(q)
	return q
}
