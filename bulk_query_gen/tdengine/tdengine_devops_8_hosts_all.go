package tdengine

import "time"
import bulkQuerygen "github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query_gen"

// tdengineDevops8Hosts produces tdengine-specific queries for the devops groupby case.
type tdengineDevops8HostsAll struct {
	tdengineDevops
}

func NewtdengineDevops8HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newtdengineDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*tdengineDevops)
	return &tdengineDevops8HostsAll{
		tdengineDevops: *underlying,
	}
}

func (d *tdengineDevops8HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageAllEightHosts(q)
	return q
}
