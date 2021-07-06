package tdengine

import (
	"time"

	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
)

// tdengineDevops8Hosts produces tdengine-specific queries for the devops groupby case.
type tdengineDevops8HostsNoInterval struct {
	tdengineDevops
}

func NewtdengineDevops8HostsNoInterval(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newtdengineDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*tdengineDevops)
	return &tdengineDevops8HostsNoInterval{
		tdengineDevops: *underlying,
	}
}

func (d *tdengineDevops8HostsNoInterval) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteEightHostsNoInterval(q)
	return q
}
