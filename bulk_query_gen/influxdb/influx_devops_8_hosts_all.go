package influxdb

import "time"
import bulkQuerygen "github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query_gen"

// InfluxDevops8Hosts produces Influx-specific queries for the devops groupby case.
type InfluxDevops8HostsAll struct {
	InfluxDevops
}

func NewInfluxQLDevops8HostsAll(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDevopsCommon(InfluxQL, dbConfig, queriesFullRange, queryInterval, scaleVar).(*InfluxDevops)
	return &InfluxDevops8HostsAll{
		InfluxDevops: *underlying,
	}
}

func NewFluxDevops8HostsAll(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDevopsCommon(Flux, dbConfig, queriesFullRange, queryInterval, scaleVar).(*InfluxDevops)
	return &InfluxDevops8HostsAll{
		InfluxDevops: *underlying,
	}
}

func (d *InfluxDevops8HostsAll) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageAllEightHosts(q)
	return q
}
