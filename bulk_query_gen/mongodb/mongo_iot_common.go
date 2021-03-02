package mongodb

import (
	"fmt"
	bulkDataGenIot "github.com/taosdata/timeseriesdatabase-comparisons/bulk_data_gen/iot"
	bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"
	"math/rand"
	"time"
)

// MongoIot produces Mongo-specific queries for the devops use case.
type MongoIot struct {
	bulkQuerygen.CommonParams
	DatabaseName string
}

// NewMongoIot makes an MongoIot object ready to generate Queries.
func NewMongoIot(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	return &MongoIot{
		CommonParams: *bulkQuerygen.NewCommonParams(interval, scaleVar),
		DatabaseName: dbConfig[bulkQuerygen.DatabaseName],
	}
}

// Dispatch fulfills the QueryGenerator interface.
func (d *MongoIot) Dispatch(i int) bulkQuerygen.Query {
	q := NewMongoQuery() // from pool
	bulkQuerygen.IotDispatchAll(d, i, q, d.ScaleVar)
	return q
}

// AverageTemperatureDayByHourOneHome populates a Query for getting the average temperature
// for one home over the course of a half a day.
func (d *MongoIot) AverageTemperatureDayByHourOneHome(q bulkQuerygen.Query) {
	d.averageTemperatureDayByHourNHomes(q.(*MongoQuery), 1, 12*time.Hour)
}

func (d *MongoIot) averageTemperatureDayByHourNHomes(qi bulkQuerygen.Query, nHomes int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nHomes]

	homes := []string{}
	for _, n := range nn {
		homes = append(homes, fmt.Sprintf(bulkDataGenIot.SmartHomeIdFormat, n))
	}

	homeClauses := []M{}
	for _, h := range homes {
		if DocumentFormat == SimpleArraysFormat {
			homeClauses = append(homeClauses, M{"home_id": h})
		} else {
			homeClauses = append(homeClauses, M{"key": "home_id", "val": h})
		}
	}

	var fieldSpec, fieldPath string
	var fieldExpr interface{}
	if DocumentFormat == SimpleArraysFormat {
		fieldSpec = "fields.temperature"
		fieldExpr = 1
		fieldPath = "fields.temperature"
	} else {
		fieldSpec = "fields"
		fieldExpr = M{ "$filter": M{ "input": "$fields", "as": "field", "cond": M{ "$eq": []string{ "$$field.key", "temperature" } } } }
		fieldPath = "fields.val"
	}

	var bucketNano = time.Hour.Nanoseconds()
	pipelineQuery := []M{
		{
			"$match": M{
				"measurement": "air_condition_room",
				"timestamp_ns": M{
					"$gte": interval.StartUnixNano(),
					"$lt":  interval.EndUnixNano(),
				},
				"tags": M{
					"$in": homeClauses,
				},
			},
		},
		{
			"$project": M{
				"_id": 0,
				"time_bucket": M{
					"$subtract": S{
						"$timestamp_ns",
						M{"$mod": S{"$timestamp_ns", bucketNano}},
					},
				},
				fieldSpec: fieldExpr, // was value: 1
				"measurement": 1,
			},
		},
		{
			"$unwind": "$fields",
		},
		{
			"$group": M{
				"_id":       M{"time_bucket": "$time_bucket", "tags": "$tags"},
				"agg_value": M{"$avg": "$"+fieldPath}, // was: $value
			},
		},
		{
			"$sort": M{"_id.time_bucket": 1},
		},
	}

	humanLabel := []byte(fmt.Sprintf("Mongo avg temperature, rand %4d homes, rand %s by 1h", nHomes, timeRange))
	q := qi.(*MongoQuery)
	q.HumanLabel = humanLabel
	q.BsonDoc = pipelineQuery
	q.DatabaseName = []byte(d.DatabaseName)
	q.CollectionName = []byte("point_data")
	q.MeasurementName = []byte("air_condition_room")
	q.FieldName = []byte("temperature")
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s (%s, %s, %s, %s)", humanLabel, interval.StartString(), q.DatabaseName, q.CollectionName, q.MeasurementName, q.FieldName))
	q.TimeStart = interval.Start
	q.TimeEnd = interval.End
	q.GroupByDuration = time.Hour
}
