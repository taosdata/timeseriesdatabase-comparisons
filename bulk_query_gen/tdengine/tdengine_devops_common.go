package tdengine

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	bulkQuerygen "github.com/liu0x54/timeseriesdatabase-comparisons/bulk_query_gen"
)

// tdengineDevops produces tdengine-specific queries for all the devops query types.
type tdengineDevops struct {
	bulkQuerygen.CommonParams
}

// NewtdengineDevops makes an tdengineDevops object ready to generate Queries.
func newtdengineDevopsCommon(interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {

	return &tdengineDevops{
		CommonParams: *bulkQuerygen.NewCommonParams(interval, scaleVar),
	}
}

// Dispatch fulfills the QueryGenerator interface.
func (d *tdengineDevops) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	bulkQuerygen.DevopsDispatchAll(d, i, q, d.ScaleVar)
	return q
}

func (d *tdengineDevops) MaxCPUUsageHourByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 1, time.Hour)
}

func (d *tdengineDevops) MaxCPUUsageHourByMinuteTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 2, time.Hour)
}

func (d *tdengineDevops) MaxCPUUsageHourByMinuteFourHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 4, time.Hour)
}

func (d *tdengineDevops) MaxCPUUsageHourByMinuteEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, time.Hour)
}
func (d *tdengineDevops) MaxCPUUsageHourByMinuteEightHosts12HR(q bulkQuerygen.Query) {
	d.maxCPUUsageHourBy10MinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, 12*time.Hour)
}
func (d *tdengineDevops) MaxCPUUsageAllByHourEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageAllBy1HourHosts(q.(*bulkQuerygen.HTTPQuery), 8, time.Hour)
}
func (d *tdengineDevops) MaxCPUUsageAllEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageAllHosts(q.(*bulkQuerygen.HTTPQuery), 8, time.Hour)
}
func (d *tdengineDevops) MaxCPUUsageHourByMinuteSixteenHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 16, time.Hour)
}

func (d *tdengineDevops) MaxCPUUsageHourByMinuteThirtyTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 32, time.Hour)
}

func (d *tdengineDevops) MaxCPUUsage12HoursByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 1, 12*time.Hour)
}

// MaxCPUUsageHourByMinuteThirtyTwoHosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1m)
func (d *tdengineDevops) maxCPUUsageHourByMinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("t_hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")
	//_ := strings.Join(hostnameClauses, " or ")
	query := fmt.Sprintf("select max(f_usage_user) from devops.cpu where (%s) and ts> \"%s\" and ts < \"%s\" interval(1m) ", combinedHostnameClause, interval.StartString(), interval.EndString())
	//	query := fmt.Sprintf("select avg(f_usage_user) from devops.cpu where (%s)  interval(1h)",  combinedHostnameClause)

	humanLabel := fmt.Sprintf("TDengine max cpu, rand %4d hosts, rand %s by 1m", nhosts, timeRange)

	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}

// MaxCPUUsageHourByMinuteThirtyTwoHosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1m)
func (d *tdengineDevops) maxCPUUsageHourBy10MinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("t_hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")
	//_ := strings.Join(hostnameClauses, " or ")
	query := fmt.Sprintf("select max(f_usage_user) from devops.cpu where (%s) and ts >= '%s' and ts < '%s' interval(10m)", combinedHostnameClause, interval.StartString(), interval.EndString())
	//	query := fmt.Sprintf("select avg(f_usage_user) from devops.cpu where (%s)  interval(1h)",  combinedHostnameClause)

	humanLabel := fmt.Sprintf("TDengine max cpu, rand %4d hosts, rand %s by 10m", nhosts, timeRange)

	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}

func (d *tdengineDevops) maxCPUUsageAllBy1HourHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	//interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("t_hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	query := fmt.Sprintf("select max(f_usage_user) from devops.cpu where (%s) interval(1h)", combinedHostnameClause)

	humanLabel := fmt.Sprintf("TDengine max cpu, rand %4d hosts,  by 1 hour", nhosts)

	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, " ", query, q)
}

func (d *tdengineDevops) maxCPUUsageAllHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	//interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("t_hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	query := fmt.Sprintf("select max(f_usage_user) from devops.cpu where (%s) ", combinedHostnameClause)

	humanLabel := fmt.Sprintf("TDengine max cpu, rand %4d hosts ", nhosts)

	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, " ", query, q)
}

// MeanCPUUsageDayByHourAllHosts populates a Query with a query that looks like:
// SELECT mean(usage_user) from cpu where time >= '$DAY_START' and time < '$DAY_END' group by time(1h),hostname
func (d *tdengineDevops) MeanCPUUsageDayByHourAllHostsGroupbyHost(qi bulkQuerygen.Query) {
	interval := d.AllInterval.RandWindow(24 * time.Hour)

	query := fmt.Sprintf("select avg(f_usage_user) from devops.cpu where ts >= '%s' and time < '%s' group by t_hostname interval(1h)", interval.StartString(), interval.EndString())

	humanLabel := fmt.Sprintf("TDengine avg cpu, all hosts, rand 1day by 1hour")
	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}

//func (d *tdengineDevops) MeanCPUUsageDayByHourAllHostsGroupbyHost(qi Query, _ int) {
//	interval := d.AllInterval.RandWindow(24*time.Hour)
//
//	v := url.Values{}
//	v.Set("db", d.DatabaseName)
//	v.Set("q", fmt.Sprintf("SELECT count(usage_user) from cpu where time >= '%s' and time < '%s' group by time(1h)", interval.StartString(), interval.EndString()))
//
//	humanLabel := "tdengine mean cpu, all hosts, rand 1day by 1hour"
//	q := qi.(*bulkQuerygen.HTTPQuery)
//	q.HumanLabel = []byte(humanLabel)
//	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
//	q.Method = []byte("GET")
//	q.Path = []byte(fmt.Sprintf("/query?%s", v.Encode()))
//	q.Body = nil
//}
func (d *tdengineDevops) getHttpQuery(humanLabel, intervalStart, query string, q *bulkQuerygen.HTTPQuery) {
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, intervalStart))

	q.Method = []byte("POST")
	q.Path = []byte("/rest/sql")
	q.Body = []byte(query)
}
