package timescaledb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// TimescaleDevops8Hosts produces Timescale-specific queries for the devops single-host case.
type TimescaleDevops8Hosts struct {
	TimescaleDevops
}

func NewTimescaleDevops8Hosts(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newTimescaleDevopsCommon(dbConfig, queriesFullRange, queryInterval, scaleVar).(*TimescaleDevops)
	return &TimescaleDevops8Hosts{
		TimescaleDevops: *underlying,
	}
}

func (d *TimescaleDevops8Hosts) Dispatch(i int) bulkQuerygen.Query {
	q := NewSQLQuery() // from pool
	d.MaxCPUUsageAllHour8Hosts(q)
	return q
}
