package tdengine

import "time"
import bulkQuerygen "github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query_gen"

// tdengineDevopsSingleHost produces tdengine-specific queries for the devops single-host case.
type tdengineDevopsSingleHost struct {
	tdengineDevops
}

func NewtdengineDevopsSingleHost(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newtdengineDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*tdengineDevops)
	return &tdengineDevopsSingleHost{
		tdengineDevops: *underlying,
	}
}

func (d *tdengineDevopsSingleHost) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteOneHost(q)
	return q
}
