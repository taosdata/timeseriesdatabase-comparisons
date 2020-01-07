package tdengine

import "time"
import bulkQuerygen "github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query_gen"

// tdengineDevops8Hosts produces tdengine-specific queries for the devops groupby case.
type tdengineDevops8Hosts12HR struct {
	tdengineDevops
}

func NewtdengineDevops8Hosts12HR(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newtdengineDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*tdengineDevops)
	return &tdengineDevops8Hosts12HR{
		tdengineDevops: *underlying,
	}
}

func (d *tdengineDevops8Hosts12HR) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteEightHosts12HR(q)
	return q
}
