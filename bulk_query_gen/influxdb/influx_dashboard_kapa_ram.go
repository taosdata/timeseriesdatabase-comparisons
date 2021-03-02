package influxdb

import "time"
import (
	"fmt"
	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
)

// InfluxDashboardKapaRam produces Influx-specific queries for the dashboard single-host case.
type InfluxDashboardKapaRam struct {
	InfluxDashboard
}

func NewInfluxQLDashboardKapaRam(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(InfluxQL, dbConfig, interval, duration, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardKapaRam{
		InfluxDashboard: *underlying,
	}
}

func NewFluxDashboardKapaRam(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(Flux, dbConfig, interval, duration, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardKapaRam{
		InfluxDashboard: *underlying,
	}
}

func (d *InfluxDashboardKapaRam) Dispatch(i int) bulkQuerygen.Query {
	q, interval := d.InfluxDashboard.DispatchCommon(i)

	var query string
	//SELECT "used_percent" FROM "telegraf"."autogen"."mem" WHERE time > :dashboardTime: AND "host"='kapacitor'
	query = fmt.Sprintf("SELECT \"used_percent\" FROM mem WHERE  hostname='kapacitor_1' and %s", d.GetTimeConstraint(interval))

	humanLabel := fmt.Sprintf("InfluxDB (%s) kapa mem used in %s", d.language.String(), interval.Duration())

	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
	return q
}
