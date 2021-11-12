package opentsdb

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"text/template"
	"time"

	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
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

func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteEightHostsTwoHr(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, 2*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteEightHostsFourHr(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, 4*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteEightHostsEightHr(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, 8*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsageHourByMinuteEightHosts12Hr(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHosts(q.(*bulkQuerygen.HTTPQuery), 8, 12*time.Hour)
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

func (d *OpenTSDBDevops) MaxCPUUsage1Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageNHosts(q.(*bulkQuerygen.HTTPQuery), 1, 24*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage16Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageNHosts(q.(*bulkQuerygen.HTTPQuery), 16, 24*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage32Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageNHosts(q.(*bulkQuerygen.HTTPQuery), 32, 24*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage64Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageNHosts(q.(*bulkQuerygen.HTTPQuery), 64, 24*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage128Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageNHosts(q.(*bulkQuerygen.HTTPQuery), 128, 24*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage256Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageNHosts(q.(*bulkQuerygen.HTTPQuery), 256, 24*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage512Hosts(q bulkQuerygen.Query) {
	d.maxCPUUsageNHosts(q.(*bulkQuerygen.HTTPQuery), 512, 24*time.Hour)
}

func (d *OpenTSDBDevops) MaxCPUUsage8HostsMixfunc(q bulkQuerygen.Query) {
	d.maxCPUUsageHourByMinuteNHostsMixfunc(q.(*bulkQuerygen.HTTPQuery), 8, time.Hour)
}



// MaxCPUUsageHourByMinute8Hosts populates a Query with a query that looks like:
// SELECT max(usage_user), count(usage_user), first(usage_user), last(usage_user) from cpu where (hostname = '$HOSTNAME_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1m)
func (d *OpenTSDBDevops) maxCPUUsageHourByMinuteNHostsMixfunc(qi bulkQuerygen.Query, nhosts int, timeRange time.Duration) {
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
		"start": {{.StartTimestamp}},
		"end": {{.EndTimestamp}},
		"timezone": "UTC",
		"queries": [
		{
			"aggregator":"max",
			"metric": "cpu.usage_user",
			"downsampler":"1m-max",
			"filters": [
			{
				"type": "literal_or",
				"tagk": "hostname",
				"filter": "{{.CombinedHostnameClause}}",
				"groupBy": false
			}
			],
			"expressions":[],
			"outputs":[
			{
				"alias":"output"
			}
			]
		},
		{
			"aggregator":"count",
			"metric": "cpu.usage_user",
			"downsampler":"1m-count",
			"filters": [
			{
				"type": "literal_or",
				"tagk": "hostname",
				"filter": "{{.CombinedHostnameClause}}",
				"groupBy": false
			}
			],
			"expressions":[],
			"outputs":[
			{
				"alias":"output"
			}
			]
		},
		{
			"aggregator":"first",
			"metric": "cpu.usage_user",
			"downsampler":"1m-first",
			"filters": [
			{
				"type": "literal_or",
				"tagk": "hostname",
				"filter": "{{.CombinedHostnameClause}}",
				"groupBy": false
			}
			],
			"expressions":[],
			"outputs":[
			{
				"alias":"output"
			}
			]
		},
		{
			"aggregator":"last",
			"metric": "cpu.usage_user",
			"downsampler":"1m-last",
			"filters": [
			{
				"type": "literal_or",
				"tagk": "hostname",
				"filter": "{{.CombinedHostnameClause}}",
				"groupBy": false
			}
			],
			"expressions":[],
			"outputs":[
			{
				"alias":"output"
			}
			]
		}
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
	q.Path = []byte("api/v1/query")
	q.Body = bodyWriter.Bytes()
	q.StartTimestamp = interval.StartUnixNano()
	q.EndTimestamp = interval.EndUnixNano()
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
		"start": {{.StartTimestamp}},
		"end": {{.EndTimestamp}},
		"timezone": "UTC",
		"queries": [
		{
			"aggregator":"max",
			"metric": "cpu.usage_user",
			"downsampler":"1m-max",
			"filters": [
			{
				"type": "literal_or",
				"tagk": "hostname",
				"filter": "{{.CombinedHostnameClause}}",
				"groupBy": false
			}
			],
			"expressions":[],
			"outputs":[
			{
				"alias":"output"
			}
			]
		}
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
	q.Path = []byte("api/v1/query")
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
		"start": {{.StartTimestamp}},
		"end": {{.EndTimestamp}},
		"timezone": "UTC",
		"queries": [
		{
			"aggregator":"max",
			"metric": "cpu.usage_user",
			"downsample":"1d-max",
			"filters": [
			{
				"type": "literal_or",
				"tagk": "hostname",
				"filter": "{{.CombinedHostnameClause}}",
				"groupBy": false
			}
			],
			"expressions":[],
			"outputs":[
			{
				"alias":"output"
			}
			]
		}
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
	q.Path = []byte("api/v1/query")
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
		"start": {{.StartTimestamp}},
		"end": {{.EndTimestamp}},
		"timezone": "UTC",
		"queries": [
		{
			"aggregator":"max",
			"metric": "cpu.usage_user",
			"downsampler":"10m-max",
			"filters": [
			{
				"type": "literal_or",
				"tagk": "hostname",
				"filter": "{{.CombinedHostnameClause}}",
				"groupBy": false
			}
			],
			"expressions":[],
			"outputs":[
			{
				"alias":"output"
			}
			]
		}
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
	q.Path = []byte("api/v1/query")
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
		"start": {{.StartTimestamp}},
		"end": {{.EndTimestamp}},
		"timezone": "UTC",
		"queries": [
		{
			"aggregator":"max",
			"metric": "cpu.usage_user",
			"downsampler":"1h-max",
			"filters": [
			{
				"type": "literal_or",
				"tagk": "hostname",
				"filter": "{{.CombinedHostnameClause}}",
				"groupBy": false
			}
			],
			"expressions":[],
			"outputs":[
			{
				"alias":"output"
			}
			]
		}
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
	q.Path = []byte("api/v1/query")
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
