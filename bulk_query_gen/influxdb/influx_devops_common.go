package influxdb

import (
	"fmt"
	bulkQuerygen "github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query_gen"
	"math/rand"
	"strings"
	"time"
)

// InfluxDevops produces Influx-specific queries for all the devops query types.
type InfluxDevops struct {
	InfluxCommon
}

// NewInfluxDevops makes an InfluxDevops object ready to generate Queries.
func newInfluxDevopsCommon(lang Language, dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {

	if _, ok := dbConfig[bulkQuerygen.DatabaseName]; !ok {
		panic("need influx database name")
	}

	return &InfluxDevops{
		InfluxCommon: *newInfluxCommon(lang, dbConfig[bulkQuerygen.DatabaseName], interval, scaleVar),
	}
}

// Dispatch fulfills the QueryGenerator interface.
func (d *InfluxDevops) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	bulkQuerygen.DevopsDispatchAll(d, i, q, d.ScaleVar)
	return q
}

func (d *InfluxDevops) MaxCPUUsageHourByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 1, time.Hour)
}

func (d *InfluxDevops) MaxCPUUsageHourByMinuteTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 2, time.Hour)
}

func (d *InfluxDevops) MaxCPUUsageHourByMinuteFourHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 4, time.Hour)
}

func (d *InfluxDevops) MaxCPUUsageHourByMinuteEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, time.Hour)
}
func (d *InfluxDevops) MaxCPUUsageHourByMinuteEightHosts12HR(q bulkQuerygen.Query) {
	d.maxCPUUsageHourBy10MinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, 12*time.Hour)
}
func (d *InfluxDevops) MaxCPUUsageAllByHourEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageAllBy1HourHosts(q.(*bulkQuerygen.HTTPQuery), 8, time.Hour)
}
func (d *InfluxDevops) MaxCPUUsageHourByMinuteSixteenHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 16, time.Hour)
}
func (d *InfluxDevops) MaxCPUUsageAllEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageAllHosts(q.(*bulkQuerygen.HTTPQuery), 8, time.Hour)
}
func (d *InfluxDevops) MaxCPUUsageHourByMinuteThirtyTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 32, time.Hour)
}

func (d *InfluxDevops) MaxCPUUsage12HoursByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 1, 12*time.Hour)
}

// MaxCPUUsageHourByMinuteThirtyTwoHosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1m)
func (d *InfluxDevops) maxCPUUsageHourByMinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		if d.language == InfluxQL {
			hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
		} else {
			hostnameClauses = append(hostnameClauses, fmt.Sprintf(`r.hostname == "%s"`, s))
		}
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf("SELECT max(usage_user) from cpu where (%s) and time >= '%s' and time < '%s' group by time(1m)", combinedHostnameClause, interval.StartString(), interval.EndString())

//		query = fmt.Sprintf("SELECT mean(usage_user) from cpu where (%s)  group by time(1h)", combinedHostnameClause)
	} else { // Flux
		query = fmt.Sprintf(`from(db:"%s") `+
			`|> range(start:%s, stop:%s) `+
			`|> filter(fn:(r) => r._measurement == "cpu" and r._field == "usage_user" and (%s)) `+
			`|> keep(columns:["_start", "_stop", "_time", "_value"]) `+
			`|> window(period:1m) `+
			`|> max() `+
			`|> yield()`,
			d.DatabaseName,
			interval.StartString(), interval.EndString(),
			combinedHostnameClause)
	}

	humanLabel := fmt.Sprintf("InfluxDB (%s) max cpu, rand %4d hosts, rand %s by 1m", d.language.String(), nhosts, timeRange)

	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}
func (d *InfluxDevops) maxCPUUsageHourBy10MinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		if d.language == InfluxQL {
			hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
		} else {
			hostnameClauses = append(hostnameClauses, fmt.Sprintf(`r.hostname == "%s"`, s))
		}
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf("SELECT max(usage_user) from cpu where (%s) and time >= '%s' and time < '%s' group by time(10m)", combinedHostnameClause, interval.StartString(), interval.EndString())

//		query = fmt.Sprintf("SELECT mean(usage_user) from cpu where (%s)  group by time(1h)", combinedHostnameClause)
	} else { // Flux
		query = fmt.Sprintf(`from(db:"%s") `+
			`|> range(start:%s, stop:%s) `+
			`|> filter(fn:(r) => r._measurement == "cpu" and r._field == "usage_user" and (%s)) `+
			`|> keep(columns:["_start", "_stop", "_time", "_value"]) `+
			`|> window(period:10m) `+
			`|> max() `+
			`|> yield()`,
			d.DatabaseName,
			interval.StartString(), interval.EndString(),
			combinedHostnameClause)
	}

	humanLabel := fmt.Sprintf("InfluxDB (%s) max cpu, rand %4d hosts, rand %s by 10m", d.language.String(), nhosts, timeRange)

	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}

func (d *InfluxDevops) maxCPUUsageAllBy1HourHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		if d.language == InfluxQL {
			hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
		} else {
			hostnameClauses = append(hostnameClauses, fmt.Sprintf(`r.hostname == "%s"`, s))
		}
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf("SELECT max(usage_user) from cpu where (%s)  group by time(1h)", combinedHostnameClause)

	} else { // Flux
		query = fmt.Sprintf(`from(db:"%s") `+
			`|> filter(fn:(r) => r._measurement == "cpu" and r._field == "usage_user" and (%s)) `+
			`|> keep(columns:["_start", "_stop", "_time", "_value"]) `+
			`|> window(period:1h) `+
			`|> max() `+
			`|> yield()`,
			d.DatabaseName,
			combinedHostnameClause)
	}

	humanLabel := fmt.Sprintf("InfluxDB (%s) max cpu, rand %4d hosts, by 1h", d.language.String(), nhosts)

	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, " ", query, q)
}
func (d *InfluxDevops) maxCPUUsageAllHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		if d.language == InfluxQL {
			hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
		} else {
			hostnameClauses = append(hostnameClauses, fmt.Sprintf(`r.hostname == "%s"`, s))
		}
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf("SELECT max(usage_user) from cpu where (%s) ", combinedHostnameClause)

	} else { // Flux
		query = fmt.Sprintf(`from(db:"%s") `+
			`|> filter(fn:(r) => r._measurement == "cpu" and r._field == "usage_user" and (%s)) `+
			`|> keep(columns:["_start", "_stop", "_time", "_value"]) `+
			`|> max() `+
			`|> yield()`,
			d.DatabaseName,
			combinedHostnameClause)
	}

	humanLabel := fmt.Sprintf("InfluxDB (%s) max cpu, rand %4d hosts ", d.language.String(), nhosts)

	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, " ", query, q)
}
// MeanCPUUsageDayByHourAllHosts populates a Query with a query that looks like:
// SELECT mean(usage_user) from cpu where time >= '$DAY_START' and time < '$DAY_END' group by time(1h),hostname
func (d *InfluxDevops) MeanCPUUsageDayByHourAllHostsGroupbyHost(qi bulkQuerygen.Query) {
	interval := d.AllInterval.RandWindow(-24 * time.Hour)

	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf("SELECT mean(usage_user) from cpu where time >= '%s' and time < '%s' group by time(1h),hostname", interval.StartString(), interval.EndString())
	} else {
		query = fmt.Sprintf(`from(db:"%s") `+
			`|> range(start:%s, stop:%s) `+
			`|> filter(fn:(r) => r._measurement == "cpu" and r._field == "usage_user") `+
			`|> keep(columns:["_start", "_stop", "hostname", "_value", "_time"]) `+
			`|> window(every:1h) `+
			`|> mean() `+
			`|> group(by:["hostname"]) `+
			`|> keep(columns:["_start", "hostname", "_value"]) `+
			`|> yield()`,
			d.DatabaseName,
			interval.StartString(), interval.EndString())
	}

	humanLabel := fmt.Sprintf("InfluxDB (%s) mean cpu, all hosts, rand 1day by 1hour", d.language.String())
	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}

//func (d *InfluxDevops) MeanCPUUsageDayByHourAllHostsGroupbyHost(qi Query, _ int) {
//	interval := d.AllInterval.RandWindow(24*time.Hour)
//
//	v := url.Values{}
//	v.Set("db", d.DatabaseName)
//	v.Set("q", fmt.Sprintf("SELECT count(usage_user) from cpu where time >= '%s' and time < '%s' group by time(1h)", interval.StartString(), interval.EndString()))
//
//	humanLabel := "Influx mean cpu, all hosts, rand 1day by 1hour"
//	q := qi.(*bulkQuerygen.HTTPQuery)
//	q.HumanLabel = []byte(humanLabel)
//	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
//	q.Method = []byte("GET")
//	q.Path = []byte(fmt.Sprintf("/query?%s", v.Encode()))
//	q.Body = nil
//}
