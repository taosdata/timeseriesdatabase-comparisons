package cassandra

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gocql/gocql"
	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
)

// CassandraDevops produces Cassandra-specific queries for all the devops query types.
type CassandraDevops struct {
	bulkQuerygen.CommonParams
	KeyspaceName string
}

// NewCassandraDevops makes an CassandraDevops object ready to generate Queries.
func newCassandraDevopsCommon(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {

	return &CassandraDevops{
		CommonParams: *bulkQuerygen.NewCommonParams(queriesFullRange, scaleVar),
		KeyspaceName: dbConfig[bulkQuerygen.DatabaseName],
	}
}

// Dispatch fulfills the QueryGenerator interface.
func (d *CassandraDevops) Dispatch(i int) bulkQuerygen.Query {
	q := NewCassandraQuery() // from pool
	bulkQuerygen.DevopsDispatchAll(d, i, q, d.ScaleVar)
	return q
}

func (d *CassandraDevops) MaxCPUUsageHourByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*CassandraQuery), 1, time.Hour)
}

func (d *CassandraDevops) MaxCPUUsageHourByMinuteTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*CassandraQuery), 2, time.Hour)
}

func (d *CassandraDevops) MaxCPUUsageHourByMinuteFourHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*CassandraQuery), 4, time.Hour)
}

func (d *CassandraDevops) MaxCPUUsageHourByMinuteEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*CassandraQuery), 8, time.Hour)
}

func (d *CassandraDevops) MaxCPUUsage12HoursBy10MinuteEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourBy10MinuteNHosts(q.(*CassandraQuery), 8, 12*time.Hour)
}

func (d *CassandraDevops) MaxCPUUsageAllEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageAllNHosts(q.(*CassandraQuery), 8, 0)
}

func (d *CassandraDevops) MaxCPUUsageAllByHourEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageAllByHourNHosts(q.(*CassandraQuery), 8, 0)
}

func (d *CassandraDevops) MaxCPUUsageHourByMinuteSixteenHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*CassandraQuery), 16, time.Hour)
}

func (d *CassandraDevops) MaxCPUUsageHourByMinuteThirtyTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*CassandraQuery), 32, time.Hour)
}

func (d *CassandraDevops) MaxCPUUsage12HoursByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*CassandraQuery), 1, 12*time.Hour)
}

// MaxCPUUsageHourByMinuteThirtyTwoHosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1m)
func (d *CassandraDevops) maxCPUUsageHourByMinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf(" '%s'", s))
	}

	combinedHostnameClause := "hostname in(" + strings.Join(hostnameClauses, " , ") + ")"

	humanLabel := fmt.Sprintf("Cassandra max cpu, rand %4d hosts, rand %s by 1m", nhosts, timeRange)
	q := qi.(*CassandraQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.AggregationType = []byte("max")
	q.MeasurementName = []byte("measurements.cpu")
	q.FieldName = []byte("usage_user")

	q.TimeStart = interval.Start
	q.TimeEnd = interval.End
	q.GroupByDuration = time.Minute

	q.TagsCondition = []byte(combinedHostnameClause)
}

// MaxCPUUsageHourByMinuteThirtyTwoHosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(10m)
func (d *CassandraDevops) maxCPUUsageHourBy10MinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf(" '%s'", s))
	}

	combinedHostnameClause := "hostname in(" + strings.Join(hostnameClauses, " , ") + ")"

	humanLabel := fmt.Sprintf("Cassandra max cpu, rand %4d hosts, rand %s by 10m", nhosts, timeRange)
	q := qi.(*CassandraQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.AggregationType = []byte("max")
	q.MeasurementName = []byte("measurements.cpu")
	q.FieldName = []byte("usage_user")

	q.TimeStart = interval.Start
	q.TimeEnd = interval.End
	q.GroupByDuration = time.Minute * 10

	q.TagsCondition = []byte(combinedHostnameClause)
}

// MaxCPUUsageHourByMinuteThirtyTwoHosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1h)
func (d *CassandraDevops) maxCPUUsageAllByHourNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	q := qi.(*CassandraQuery)
	timeCluster := gocql.NewCluster("localhost:9042")
	timeCluster.Keyspace = "measurements"
	timeCluster.Consistency = gocql.Quorum
	timeCluster.Timeout = 60 * time.Second
	timeCluster.ProtoVersion = 4
	session, err := timeCluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	var timeReturn int64

	if err := session.Query(`select min(time) from measurements.cpu ;`).Consistency(gocql.One).Scan(&timeReturn); err != nil {
		log.Fatal(err)
	}

	q.TimeStart = time.Unix(timeReturn/1000000000, timeReturn%1000000000).UTC()

	if err := session.Query(`select max(time) from measurements.cpu ;`).Consistency(gocql.One).Scan(&timeReturn); err != nil {
		log.Fatal(err)
	}

	q.TimeEnd = time.Unix(timeReturn/1000000000, timeReturn%1000000000).UTC()

	interval := d.AllInterval.NewTimeInterval(q.TimeStart, q.TimeEnd)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf(" '%s'", s))
	}

	combinedHostnameClause := "hostname in(" + strings.Join(hostnameClauses, " , ") + ")"

	humanLabel := fmt.Sprintf("Cassandra max cpu, rand %4d hosts, rand %s by 1h", nhosts, timeRange)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.AggregationType = []byte("max")
	q.MeasurementName = []byte("measurements.cpu")
	q.FieldName = []byte("usage_user")

	q.GroupByDuration = time.Hour

	q.TagsCondition = []byte(combinedHostnameClause)
}

// MaxCPUUsageHourByMinuteThirtyTwoHosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1h)
func (d *CassandraDevops) maxCPUUsageAllNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	q := qi.(*CassandraQuery)
	timeCluster := gocql.NewCluster("localhost:9042")
	timeCluster.Keyspace = "measurements"
	timeCluster.Consistency = gocql.Quorum
	timeCluster.Timeout = 60 * time.Second
	timeCluster.ProtoVersion = 4
	session, err := timeCluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	//get the max time and the min from the database, and set the query interval to be the full range
	var timeReturn int64
	if err := session.Query(`select min(time) from measurements.cpu ;`).Consistency(gocql.One).Scan(&timeReturn); err != nil {
		log.Fatal(err)
	}
	q.TimeEnd = q.TimeStart
	interval := d.AllInterval.NewTimeInterval(q.TimeStart, q.TimeEnd)

	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf(" '%s'", s))
	}

	combinedHostnameClause := "hostname in(" + strings.Join(hostnameClauses, " , ") + ")"

	humanLabel := fmt.Sprintf("Cassandra max cpu, rand %4d hosts, rand %s by 1h", nhosts, timeRange)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.AggregationType = []byte("max")
	q.MeasurementName = []byte("measurements.cpu")
	q.FieldName = []byte("usage_user")
	q.GroupByDuration = interval.Duration()

	q.TagsCondition = []byte(combinedHostnameClause)
}

// MeanCPUUsageDayByHourAllHosts populates a Query with a query that looks like:
// SELECT mean(usage_user) from cpu where time >= '$DAY_START' and time < '$DAY_END' group by time(1h),hostname
func (d *CassandraDevops) MeanCPUUsageDayByHourAllHostsGroupbyHost(qi bulkQuerygen.Query) {
	interval := d.AllInterval.RandWindow(24 * time.Hour)

	humanLabel := "Cassandra mean cpu, all hosts, rand 1day by 1hour"
	q := qi.(*CassandraQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))

	q.AggregationType = []byte("avg")
	q.MeasurementName = []byte("cpu")
	q.FieldName = []byte("usage_user")

	q.TimeStart = interval.Start
	q.TimeEnd = interval.End
	q.GroupByDuration = time.Hour
}

//func (d *CassandraDevops) MeanCPUUsageDayByHourAllHostsGroupbyHost(qi Query, _ int) {
//	interval := d.AllInterval.RandWindow(24*time.Hour)
//
//	v := url.Values{}
//	v.Set("db", d.KeyspaceName)
//	v.Set("q", fmt.Sprintf("SELECT count(usage_user) from cpu where time >= '%s' and time < '%s' group by time(1h)", interval.StartString(), interval.EndString()))
//
//	humanLabel := "Cassandra mean cpu, all hosts, rand 1day by 1hour"
//	q := qi.(*CassandraQuery)
//	q.HumanLabel = []byte(humanLabel)
//	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
//	q.Method = []byte("GET")
//	q.Path = []byte(fmt.Sprintf("/query?%s", v.Encode()))
//	q.Body = nil
//}
