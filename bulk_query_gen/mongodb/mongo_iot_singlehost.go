package mongodb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// MongoIotSingleHost produces Mongo-specific queries for the devops single-host case.
type MongoIotSingleHost struct {
	MongoIot
}

func NewMongoIotSingleHost(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := NewMongoIot(dbConfig, queriesFullRange, queryInterval, scaleVar).(*MongoIot)
	return &MongoIotSingleHost{
		MongoIot: *underlying,
	}
}

func (d *MongoIotSingleHost) Dispatch(i int) bulkQuerygen.Query {
	q := NewMongoQuery() // from pool
	d.AverageTemperatureDayByHourOneHome(q)
	return q
}
