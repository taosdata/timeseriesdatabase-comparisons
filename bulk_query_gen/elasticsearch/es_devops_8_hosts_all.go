package elasticsearch

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// ElasticSearchDevops8Hosts produces ElasticSearch-specific queries for the devops groupby case.
type ElasticSearchDevops8HostsAll struct {
	ElasticSearchDevops
}

func NewElasticSearchDevops8HostsAll(_ bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := NewElasticSearchDevops(queriesFullRange, scaleVar).(*ElasticSearchDevops)
	return &ElasticSearchDevops8HostsAll{
		ElasticSearchDevops: *underlying,
	}
}

func (d *ElasticSearchDevops8HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteEightHosts(q)
	return q
}
