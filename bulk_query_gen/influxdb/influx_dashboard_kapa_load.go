package influxdb

import "time"
import (
	"fmt"
	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
)

// InfluxDashboardKapaLoad produces Influx-specific queries for the dashboard single-host case.
type InfluxDashboardKapaLoad struct {
	InfluxDashboard
}

func NewInfluxQLDashboardKapaLoad(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(InfluxQL, dbConfig, interval, duration, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardKapaLoad{
		InfluxDashboard: *underlying,
	}
}

func NewFluxDashboardKapaLoad(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(Flux, dbConfig, interval, duration, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardKapaLoad{
		InfluxDashboard: *underlying,
	}
}

func (d *InfluxDashboardKapaLoad) Dispatch(i int) bulkQuerygen.Query {
	q, interval := d.InfluxDashboard.DispatchCommon(i)

	var query string
	//SELECT "load5", "load15", "load1" FROM "telegraf"."autogen"."system" WHERE time > :dashboardTime: AND "host"='kapacitor'
	query = fmt.Sprintf("SELECT \"load5\", \"load15\", \"load1\" FROM system WHERE hostname='kapacitor_1' and %s", d.GetTimeConstraint(interval))

	humanLabel := fmt.Sprintf("InfluxDB (%s) kapa load 1,5,15 in %s", d.language.String(), interval.Duration())

	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
	return q
}
