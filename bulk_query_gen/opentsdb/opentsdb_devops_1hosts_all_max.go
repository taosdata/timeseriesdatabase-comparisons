package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops1HostsAll produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops1HostsAll struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops1HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops1HostsAll{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops1HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage1Hosts(q)
	return q
}
