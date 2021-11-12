package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops16HostsAll produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops16HostsAll struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops16HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops16HostsAll{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops16HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage16Hosts(q)
	return q
}
