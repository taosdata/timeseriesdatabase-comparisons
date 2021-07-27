package timescaledb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// TimescaleDevops8HostsAllBy1Hr produces Timescale-specific queries for the devops single-host case.
type TimescaleDevops8HostsAllBy1Hr struct {
	TimescaleDevops
}

func NewTimescaleDevops8HostsAllBy1Hr(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newTimescaleDevopsCommon(dbConfig, queriesFullRange, queryInterval, scaleVar).(*TimescaleDevops)
	return &TimescaleDevops8HostsAllBy1Hr{
		TimescaleDevops: *underlying,
	}
}

func (d *TimescaleDevops8HostsAllBy1Hr) Dispatch(i int) bulkQuerygen.Query {
	q := NewSQLQuery() // from pool
	d.MaxCPUUsageByHour8Hostss(q)
	return q
}
