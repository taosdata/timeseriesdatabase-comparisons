package tdengine

import "time"
import bulkQuerygen "github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query_gen"

// tdengineDevopsSingleHost12hr produces tdengine-specific queries for the devops single-host case over a 12hr period.
type tdengineDevopsSingleHost12hr struct {
	tdengineDevops
}

func NewtdengineDevopsSingleHost12hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newtdengineDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*tdengineDevops)
	return &tdengineDevopsSingleHost12hr{
		tdengineDevops: *underlying,
	}
}

func (d *tdengineDevopsSingleHost12hr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsage12HoursByMinuteOneHost(q)
	return q
}
