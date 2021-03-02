package mongodb

import "time"
import bulkQuerygen "github.com/taosdata/timeseriesdatabase-comparisons/bulk_query_gen"

// MongoDevopsSingleHost produces Mongo-specific queries for the devops single-host case.
type MongoDevopsSingleHost struct {
	MongoDevops
}

func NewMongoDevopsSingleHost(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := NewMongoDevops(dbConfig, queriesFullRange, queryInterval, scaleVar).(*MongoDevops)
	return &MongoDevopsSingleHost{
		MongoDevops: *underlying,
	}
}

func (d *MongoDevopsSingleHost) Dispatch(i int) bulkQuerygen.Query {
	q := NewMongoQuery() // from pool
	d.MaxCPUUsageHourByMinuteOneHost(q)
	return q
}
