package opentsdb

import (
	"bytes"
	"fmt"
	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
	"math/rand"
	"net/url"
	"strings"
	"text/template"
	"time"
)

// OpenTSDBDevops produces OpenTSDB-specific queries for all the devops query types.
type OpenTSDBDevops struct {
	bulkQuerygen.CommonParams
}

// NewOpenTSDBDevops makes an OpenTSDBDevops object ready to generate Queries.
func newOpenTSDBDevopsCommon(interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {

	return &OpenTSDBDevops{
		CommonParams: *bulkQuerygen.NewCommonParams(interval, scaleVar),
	}
}

// Dispatch fulfills the QueryGenerator interface.
func (d *OpenTSDBDevops) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	bulkQuerygen.DevopsDispatchAll(d, i, q, d.ScaleVar)
	return q
}

func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 1, time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 2, time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteFourHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 4, time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsageHourByHourEightHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByHourNHosts(q.(*bulkQuerygen.HTTPQuery), 8, 24*time.Hour)
}


func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteSixteenHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 16, time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteThirtyTwoHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 32, time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage12HoursByMinuteOneHost(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 1, 12*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage12HourByTenMinuteNHosts(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByTenMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, 12*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage8Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageNHosts(q.(*bulkQuerygen.HTTPQuery), 8, 24*time.Hour)
}


// MaxCPUUsageHourByMinute8Hosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1m)
func (d *OpenTSDBDevops) maxCPUUsageHourByMinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	combinedHostnameClause := strings.Join(hostnames, "|")

	// opentsdb cannot handle RFC3339, nor can it handle nanoseconds,
	// so use unix epoch time in milliseconds:
	startTimestamp := interval.StartUnixNano() / 1e6
	endTimestamp := interval.EndUnixNano() / 1e6

	const tmplString = `
{
   "time": {
       "start": {{.StartTimestamp}},
       "end": {{.EndTimestamp}},
       "timezone": "UTC",
       "aggregator":"max",
       "downsampler":{"interval":"1m","aggregator":"max"}
   },
   "filters": [
       {
           "tags": [
               {
                   "type": "literal_or",
                   "tagk": "hostname",
                   "filter": "{{.CombinedHostnameClause}}",
                   "groupBy": false
               }
           ],
           "id": "f1"
       }
   ],
   "metrics": [
       {
           "id": "a",
           "metric": "cpu.usage_user",
           "filter": "f1",
           "fillPolicy":{"policy":"zero"}
       }
   ],
    "expressions":[
   ],
    "outputs":[
      {"id":"a", "alias":"output"}
    ]
}
`

	tmpl := template.Must(template.New("tmpl").Parse(tmplString))
	bodyWriter := new(bytes.Buffer)

	arg := struct {
		StartTimestamp, EndTimestamp int64
		CombinedHostnameClause       string
	}{
		startTimestamp,
		endTimestamp,
		combinedHostnameClause,
	}
	err := tmpl.Execute(bodyWriter, arg)
	if err != nil {
		panic("logic error")
	}

	humanLabel := fmt.Sprintf("OpenTSDB max cpu, rand %4d hosts, rand %s by 1m", nhosts, timeRange)
	q := qi.(*bulkQuerygen.HTTPQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")
	q.Path = []byte("/api/v1/query/exp")
	q.Body = bodyWriter.Bytes()
	q.StartTimestamp = interval.StartUnixNano()
	q.EndTimestamp = interval.EndUnixNano()
}


// MaxCPUUsage8Hosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' 
func (d *OpenTSDBDevops) maxCPUUsageNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	combinedHostnameClause := strings.Join(hostnames, "|")

	// opentsdb cannot handle RFC3339, nor can it handle nanoseconds,
	// so use unix epoch time in milliseconds:
	startTimestamp := interval.StartUnixNano() / 1e6
	endTimestamp := interval.EndUnixNano() / 1e6

	const tmplString = `
{
   "time": {
       "start": {{.StartTimestamp}},
       "end": {{.EndTimestamp}},
       "timezone": "UTC",
       "aggregator":"max",
       "downsampler":{"interval":"1d","aggregator":"max"}
   },
   "filters": [
       {
           "tags": [
               {
                   "type": "literal_or",
                   "tagk": "hostname",
                   "filter": "{{.CombinedHostnameClause}}",
                   "groupBy": false
               }
           ],
           "id": "f1"
       }
   ],
   "metrics": [
       {
           "id": "a",
           "metric": "cpu.usage_user",
           "filter": "f1",
           "fillPolicy":{"policy":"zero"}
       }
   ],
    "expressions":[
   ],
    "outputs":[
      {"id":"a", "alias":"output"}
    ]
}
`

	tmpl := template.Must(template.New("tmpl").Parse(tmplString))
	bodyWriter := new(bytes.Buffer)

	arg := struct {
		StartTimestamp, EndTimestamp int64
		CombinedHostnameClause       string
	}{
		startTimestamp,
		endTimestamp,
		combinedHostnameClause,
	}
	err := tmpl.Execute(bodyWriter, arg)
	if err != nil {
		panic("logic error")
	}

	humanLabel := fmt.Sprintf("OpenTSDB max cpu, rand %4d hosts, rand %s by alltime", nhosts, timeRange)
	q := qi.(*bulkQuerygen.HTTPQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")
	q.Path = []byte("/api/v1/query/exp")
	q.Body = bodyWriter.Bytes()
	q.StartTimestamp = interval.StartUnixNano()
	q.EndTimestamp = interval.EndUnixNano()
}

// MaxCPUUsageHourByTenMinute8Hosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(10m)
func (d *OpenTSDBDevops) maxCPUUsageHourByTenMinuteNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	combinedHostnameClause := strings.Join(hostnames, "|")

	// opentsdb cannot handle RFC3339, nor can it handle nanoseconds,
	// so use unix epoch time in milliseconds:
	startTimestamp := interval.StartUnixNano() / 1e6
	endTimestamp := interval.EndUnixNano() / 1e6

	const tmplString = `
{
   "time": {
       "start": {{.StartTimestamp}},
       "end": {{.EndTimestamp}},
       "timezone": "UTC",
       "aggregator":"max",
       "downsampler":{"interval":"10m","aggregator":"max"}
   },
   "filters": [
       {
           "tags": [
               {
                   "type": "literal_or",
                   "tagk": "hostname",
                   "filter": "{{.CombinedHostnameClause}}",
                   "groupBy": false
               }
           ],
           "id": "f1"
       }
   ],
   "metrics": [
       {
           "id": "a",
           "metric": "cpu.usage_user",
           "filter": "f1",
           "fillPolicy":{"policy":"zero"}
       }
   ],
    "expressions":[
   ],
    "outputs":[
      {"id":"a", "alias":"output"}
    ]
}
`

	tmpl := template.Must(template.New("tmpl").Parse(tmplString))
	bodyWriter := new(bytes.Buffer)

	arg := struct {
		StartTimestamp, EndTimestamp int64
		CombinedHostnameClause       string
	}{
		startTimestamp,
		endTimestamp,
		combinedHostnameClause,
	}
	err := tmpl.Execute(bodyWriter, arg)
	if err != nil {
		panic("logic error")
	}

	humanLabel := fmt.Sprintf("OpenTSDB max cpu, rand %4d hosts, rand %s by 10m", nhosts, timeRange)
	q := qi.(*bulkQuerygen.HTTPQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")
	q.Path = []byte("/api/v1/query/exp")
	q.Body = bodyWriter.Bytes()
	q.StartTimestamp = interval.StartUnixNano()
	q.EndTimestamp = interval.EndUnixNano()
}

// MaxCPUUsageHourByHourNHosts populates a Query with a query that looks like:
// SELECT max(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1h)
func (d *OpenTSDBDevops) maxCPUUsageHourByHourNHosts(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	combinedHostnameClause := strings.Join(hostnames, "|")

	// opentsdb cannot handle RFC3339, nor can it handle nanoseconds,
	// so use unix epoch time in milliseconds:
	startTimestamp := interval.StartUnixNano() / 1e6
	endTimestamp := interval.EndUnixNano() / 1e6

	const tmplString = `
{
   "time": {
       "start": {{.StartTimestamp}},
       "end": {{.EndTimestamp}},
       "timezone": "UTC",
       "aggregator":"max",
       "downsampler":{"interval":"1h","aggregator":"max"}
   },
   "filters": [
       {
           "tags": [
               {
                   "type": "literal_or",
                   "tagk": "hostname",
                   "filter": "{{.CombinedHostnameClause}}",
                   "groupBy": false
               }
           ],
           "id": "f1"
       }
   ],
   "metrics": [
       {
           "id": "a",
           "metric": "cpu.usage_user",
           "filter": "f1",
           "fillPolicy":{"policy":"zero"}
       }
   ],
    "expressions":[
   ],
    "outputs":[
      {"id":"a", "alias":"output"}
    ]
}
`

	tmpl := template.Must(template.New("tmpl").Parse(tmplString))
	bodyWriter := new(bytes.Buffer)

	arg := struct {
		StartTimestamp, EndTimestamp int64
		CombinedHostnameClause       string
	}{
		startTimestamp,
		endTimestamp,
		combinedHostnameClause,
	}
	err := tmpl.Execute(bodyWriter, arg)
	if err != nil {
		panic("logic error")
	}

	humanLabel := fmt.Sprintf("OpenTSDB max cpu, rand %4d hosts, rand %s by 1hour", nhosts, timeRange)
	q := qi.(*bulkQuerygen.HTTPQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")
	q.Path = []byte("/api/v1/query/exp")
	q.Body = bodyWriter.Bytes()
	q.StartTimestamp = interval.StartUnixNano()
	q.EndTimestamp = interval.EndUnixNano()
}

// MeanCPUUsageDayByHourAllHosts populates a Query with a query that looks like:
// SELECT mean(usage_user) from cpu where time >= '$DAY_START' and time < '$DAY_END' group by time(1h),hostname
func (d *OpenTSDBDevops) MeanCPUUsageDayByHourAllHostsGroupbyHost(qi bulkQuerygen.Query) {
	interval := d.AllInterval.RandWindow(24 * time.Hour)

	v := url.Values{}
	v.Set("q", fmt.Sprintf("SELECT mean(usage_user) from cpu where time >= '%s' and time < '%s' group by time(1h),hostname", interval.StartString(), interval.EndString()))

	humanLabel := "OpenTSDB mean cpu, all hosts, rand 1day by 1hour"
	q := qi.(*bulkQuerygen.HTTPQuery)
	q.HumanLabel = []byte(humanLabel)
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("GET")
	q.Path = []byte(fmt.Sprintf("/query?%s", v.Encode()))
	q.Body = nil
}

//func (d *OpenTSDBDevops) MeanCPUUsageDayByHourAllHostsGroupbyHost(qi Query, _ int) {
//	interval := d.AllInterval.RandWindow(24*time.Hour)
//
//	v := url.Values{}
//	v.Set("db", d.DatabaseName)
//	v.Set("q", fmt.Sprintf("SELECT count(usage_user) from cpu where time >= '%s' and time < '%s' group by time(1h)", interval.StartString(), interval.EndString()))
//
//	humanLabel := "OpenTSDB mean cpu, all hosts, rand 1day by 1hour"
//	q := qi.(*bulkQuerygen.HTTPQuery)
//	q.HumanLabel = []byte(humanLabel)
//	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
//	q.Method = []byte("GET")
//	q.Path = []byte(fmt.Sprintf("/query?%s", v.Encode()))
//	q.Body = nil
//}
