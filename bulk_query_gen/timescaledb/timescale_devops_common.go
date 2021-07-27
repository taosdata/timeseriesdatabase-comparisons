package timescaledb

import (
	"fmt"
	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
	"math/rand"
	"strings"
	"time"
)

// TimescaleDevops produces Timescale-specific queries for all the devops query types.
type TimescaleDevops struct {
	bulkQuerygen.CommonParams
	DatabaseName string
}

// newTimescaleDevopsCommon makes an TimescaleDevops object ready to generate Queries.
func newTimescaleDevopsCommon(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	if _, ok := dbConfig[bulkQuerygen.DatabaseName]; !ok {
		panic("need timescale database name")
	}

	return &TimescaleDevops{
		CommonParams: *bulkQuerygen.NewCommonParams(interval, scaleVar),
		DatabaseName: dbConfig[bulkQuerygen.DatabaseName],
	}
}

// Dispatch fulfills the QueryGenerator interface.
func (d *TimescaleDevops) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	bulkQuerygen.DevopsDispatchAll(d, i, q, d.ScaleVar)
	return q
}

func (d *TimescaleDevops) MaxCPUUsageHourByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q, 1, time.Hour)
}

func (d *TimescaleDevops) MaxCPUUsageHourByMinuteTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q, 2, time.Hour)
}

func (d *TimescaleDevops) MaxCPUUsageHourByMinuteFourHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q, 4, time.Hour)
}

func (d *TimescaleDevops) MaxCPUUsageHourByMinuteEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q, 8, time.Hour)
}

func (d *TimescaleDevops) MaxCPUUsageHourByMinuteSixteenHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q, 16, time.Hour)
}

func (d *TimescaleDevops) MaxCPUUsageHourByMinuteThirtyTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q, 32, time.Hour)
}

func (d *TimescaleDevops) MaxCPUUsage12HoursByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q, 1, 12*time.Hour)
}

// 测试用例1，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据的最大值。
func (d *TimescaleDevops) MaxCPUUsageAllHour8Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageAllHourNHosts(q, 8)
}

// 测试用例2，查询所有数据中，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1小时为粒度，查询每1小时的最大值。
func (d *TimescaleDevops) MaxCPUUsageByHour8Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageByHourNHosts(q, 8, 12*time.Hour)
}

// 测试用例3，测试用例3，随机查询12个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以10分钟为粒度，查询每10分钟的最大值
func (d *TimescaleDevops) MaxCPUUsage12HourBy10min8Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageBy10minNHosts(q, 8, 12*time.Hour)
}

// 测试用例4，随机查询1个小时的数据，用8个hostname标签进行匹配，匹配出这8个hostname对应的模拟服务器CPU数据中的usage_user这个监控数据，以1分钟为粒度，查询每1分钟的最大值
func (d *TimescaleDevops) MaxCPUUsage1HourBy1min8Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageBy1minNHosts(q, 8, time.Hour)
}

// MaxCPUUsageHourByMinuteThirtyTwoHosts populates a Query with a query that looks like:
// select time_bucket(60000000000,time) as time1min,max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >=$HOUR_START and time < $HOUR_END group by time1min order by time1min;
func (d *TimescaleDevops) maxCPUUsageHourByMinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	humanLabel := fmt.Sprintf("Timescale max cpu, rand %4d hosts, rand %s by 1m", nhosts, timeRange)

	q := qi.(*SQLQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.QuerySQL = []byte(fmt.Sprintf("select time_bucket(60000000000,time) as time1min,max(usage_user) from cpu where (%s) and time >=%d and time < %d group by time1min  ", combinedHostnameClause, interval.StartUnixNano(), interval.EndUnixNano()))
}

// maxCPUUsageHourBy5MinuteNHosts populates a Query with a query that looks like:
// select time_bucket(60000000000,time) as time5min,max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >=$HOUR_START and time < $HOUR_END group by time5min order by time5min;
func (d *TimescaleDevops) maxCPUUsageHourBy5MinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	humanLabel := fmt.Sprintf("Timescale max cpu, rand %4d hosts, rand %s by 1m", nhosts, timeRange)

	q := qi.(*SQLQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.QuerySQL = []byte(fmt.Sprintf("select time_bucket(60000000000,time) as time1min,max(usage_user) from cpu where (%s) and time >=%d and time < %d group by time5min ", combinedHostnameClause, interval.StartUnixNano(), interval.EndUnixNano()))
}


// MeanCPUUsageDayByHourAllHosts populates a Query with a query that looks like:
// SELECT mean(usage_user) from cpu where time >= '$DAY_START' and time < '$DAY_END' group by time(1h),hostname
func (d *TimescaleDevops) MeanCPUUsageDayByHourAllHostsGroupbyHost(qi bulkQuerygen.Query) {
	interval := d.AllInterval.RandWindow(24 * time.Hour)

	humanLabel := "Timescale mean cpu, all hosts, rand 1day by 1hour"
	q := qi.(*SQLQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.QuerySQL = []byte(fmt.Sprintf("select time_bucket(3600000000000,time) as time1hour,avg(usage_user) from cpu where time >=%d and time < %d group by time1hour,hostname order by time1hour", interval.StartUnixNano(), interval.EndUnixNano()))
}


// maxCPUUsageAllHourNHosts populates a Query with a query that looks like:
// select max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') ;
func (d *TimescaleDevops) maxCPUUsageAllHourNHosts(qi bulkQuerygen.Query, nhosts int) {
	// interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	humanLabel := fmt.Sprintf("Timescale max cpu, rand %4d hosts, all time", nhosts )

	q := qi.(*SQLQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: all time", humanLabel ))

	q.QuerySQL = []byte(fmt.Sprintf("select max(usage_user) from cpu where (%s) ", combinedHostnameClause))
}

// maxCPUUsageByHourNHosts populates a Query with a query that looks like:
// select time_bucket('1 hour',time) as time1hour,max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >=$HOUR_START and time < $HOUR_END group by time1hour;
func (d *TimescaleDevops) maxCPUUsageByHourNHosts(qi bulkQuerygen.Query, nhosts int) {
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	humanLabel := fmt.Sprintf("Timescale max cpu, rand %4d hosts, all time", nhosts)

	q := qi.(*SQLQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: all time", humanLabel))

	q.QuerySQL = []byte(fmt.Sprintf("select time_bucket('1 hour',time) as time1hour,max(usage_user) from cpu where (%s) group by time1hour ", combinedHostnameClause))
}


// maxCPUUsageBy10minNHosts populates a Query with a query that looks like:
// select time_bucket('10 minutes',time) as time10min,max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >=$HOUR_START and time < $HOUR_END group by time10min ;
func (d *TimescaleDevops) maxCPUUsageBy10minNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	humanLabel := fmt.Sprintf("Timescale max cpu, rand %4d hosts, rand %s by 1m", nhosts, timeRange)

	q := qi.(*SQLQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.QuerySQL = []byte(fmt.Sprintf("select time_bucket('10 minutes',time) as time10min,max(usage_user) from cpu where (%s) time >=%d and time < %d group by time10min ", combinedHostnameClause, interval.StartUnixNano(), interval.EndUnixNano()))
}

// maxCPUUsageBy1minNHosts populates a Query with a query that looks like:
// select time_bucket('1 minute',time) as time1min,max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >=$HOUR_START and time < $HOUR_END group by time1min ;
func (d *TimescaleDevops) maxCPUUsageBy1minNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("hostname = '%s'", s))
	}

	combinedHostnameClause := strings.Join(hostnameClauses, " or ")

	humanLabel := fmt.Sprintf("Timescale max cpu, rand %4d hosts, rand %s by 1m", nhosts, timeRange)

	q := qi.(*SQLQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.QuerySQL = []byte(fmt.Sprintf("select time_bucket('1 minute',time) as time1min,max(usage_user) from cpu where (%s) time >=%d and time < %d group by time1min ", combinedHostnameClause, interval.StartUnixNano(), interval.EndUnixNano()))
}