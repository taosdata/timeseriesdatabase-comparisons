package cassandra

import (
	"time"

	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
)

// CassandraDevops8HostsAllBy1Hr produces Cassandra-specific queries for the devops groupby case.
type CassandraDevops8HostsAllBy1Hr struct {
	CassandraDevops
}

func NewCassandraDevops8HostsAllBy1Hr(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newCassandraDevopsCommon(dbConfig, queriesFullRange, queryInterval, scaleVar).(*CassandraDevops)
	return &CassandraDevops8HostsAllBy1Hr{
		CassandraDevops: *underlying,
	}
}

func (d *CassandraDevops8HostsAllBy1Hr) Dispatch(i int) bulkQuerygen.Query {
	q := NewCassandraQuery() // from pool
	d.MaxCPUUsageAllByHourEightHosts(q)
	return q
}
