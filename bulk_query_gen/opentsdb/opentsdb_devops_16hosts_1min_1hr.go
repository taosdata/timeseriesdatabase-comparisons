package opentsdb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// OpenTSDBDevops16Hosts produces OpenTSDB-specific queries for the devops groupby case.
type OpenTSDBDevops16Hosts struct {
	OpenTSDBDevops
}

func NewOpenTSDBDevops16Hosts(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newOpenTSDBDevopsCommon(queriesFullRange, queryInterval, scaleVar).(*OpenTSDBDevops)
	return &OpenTSDBDevops16Hosts{
		OpenTSDBDevops: *underlying,
	}
}

func (d *OpenTSDBDevops16Hosts) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteSixteenHosts(q)
	return q
}
