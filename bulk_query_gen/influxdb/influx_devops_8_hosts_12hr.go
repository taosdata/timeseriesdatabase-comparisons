package influxdb

import "time"
import bulkQuerygen "github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query_gen"

// InfluxDevops8Hosts produces Influx-specific queries for the devops groupby case.
type InfluxDevops8Hosts12HR struct {
	InfluxDevops
}

func NewInfluxQLDevops8Hosts12HR(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDevopsCommon(InfluxQL, dbConfig, queriesFullRange, queryInterval, scaleVar).(*InfluxDevops)
	return &InfluxDevops8Hosts12HR{
		InfluxDevops: *underlying,
	}
}

func NewFluxDevops8Hosts12HR(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDevopsCommon(Flux, dbConfig, queriesFullRange, queryInterval, scaleVar).(*InfluxDevops)
	return &InfluxDevops8Hosts12HR{
		InfluxDevops: *underlying,
	}
}

func (d *InfluxDevops8Hosts12HR) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	d.MaxCPUUsageHourByMinuteEightHosts12HR(q)
	return q
}
