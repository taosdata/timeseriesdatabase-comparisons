package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops256HostsAll produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops256HostsAll struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops256HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops256HostsAll{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops256HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage256Hosts(q)
	return q
}
