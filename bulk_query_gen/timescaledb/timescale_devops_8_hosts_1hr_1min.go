package timescaledb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// TimescaleDevops8Hosts1Hr1min produces Timescale-specific queries for the devops single-host case.
type TimescaleDevops8Hosts1Hr1min struct {
	TimescaleDevops
}

func NewTimescaleDevops8Hosts1Hr1min(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newTimescaleDevopsCommon(dbConfig, queriesFullRange, queryInterval, scaleVar).(*TimescaleDevops)
	return &TimescaleDevops8Hosts1Hr1min{
		TimescaleDevops: *underlying,
	}
}

func (d *TimescaleDevops8Hosts1Hr1min) Dispatch(i int) bulkQuerygen.Query {
	q := NewSQLQuery() // from pool
	d.MaxCPUUsage1HourBy1min8Hosts(q)
	return q
}
