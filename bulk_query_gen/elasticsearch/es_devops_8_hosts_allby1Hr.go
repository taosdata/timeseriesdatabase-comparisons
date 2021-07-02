package elasticsearch

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// ElasticSearchDevops8Hosts produces ElasticSearch-specific queries for the devops groupby case.
type ElasticSearchDevops8HostsAllBy1Hr struct {
	ElasticSearchDevops
}

func NewElasticSearchDevops8HostsAllBy1Hr(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := NewElasticSearchDevops(queriesFullRange, scaleVar).(*ElasticSearchDevops)
	return &ElasticSearchDevops8HostsAllBy1Hr{
		ElasticSearchDevops: *underlying,
	}
}

func (d *ElasticSearchDevops8HostsAllBy1Hr) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByHourEightHosts(q)
	return q
}