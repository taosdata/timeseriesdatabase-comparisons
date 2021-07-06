package cassandra

import (
	"time"

	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
)

// CassandraDevops8HostsAll produces Cassandra-specific queries for the devops groupby case.
type CassandraDevops8HostsAll struct {
	CassandraDevops
}

func NewCassandraDevops8HostsAll(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newCassandraDevopsCommon(dbConfig, queriesFullRange, queryInterval, scaleVar).(*CassandraDevops)
	return &CassandraDevops8HostsAll{
		CassandraDevops: *underlying,
	}
}

func (d *CassandraDevops8HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := NewCassandraQuery() // from pool
	d.MaxCPUUsageAllEightHosts(q)
	return q
}
