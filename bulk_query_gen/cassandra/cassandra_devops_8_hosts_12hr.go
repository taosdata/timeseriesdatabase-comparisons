package cassandra

import (
	"time"

	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
)

// CassandraDevops8Hosts12HR produces Cassandra-specific queries for the devops groupby case.
type CassandraDevops8Hosts12HR struct {
	CassandraDevops
}

func NewCassandraDevops8Hosts12HR(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newCassandraDevopsCommon(dbConfig, queriesFullRange, queryInterval, scaleVar).(*CassandraDevops)
	return &CassandraDevops8Hosts12HR{
		CassandraDevops: *underlying,
	}
}

func (d *CassandraDevops8Hosts12HR) Dispatch(i int) bulkQuerygen.Query {
	q := NewCassandraQuery() // from pool
	d.MaxCPUUsage12HoursBy10MinuteEightHosts(q)
	return q
}
