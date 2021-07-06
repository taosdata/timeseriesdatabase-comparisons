package main

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

// A QueryPlan is a strategy used to fulfill an HLQuery.
type QueryPlan interface {
	Execute(*gocql.Session, int) ([]CQLResult, error)
	DebugQueries(int)
}

// A QueryPlanWithServerAggregation fulfills an HLQuery by performing
// aggregation on both the server and the client. This results in more
// round-trip requests, but uses the server to aggregate over large datasets.
//
// It has 1) an Aggregator, which merges data on the client, and 2) a map of
// time interval buckets to CQL queries, which are used to retrieve data
// relevant to each bucket.
type QueryPlanWithServerAggregation struct {
	AggregatorLabel    string
	BucketedCQLQueries map[TimeInterval]CQLQuery
}

// NewQueryPlanWithServerAggregation builds a QueryPlanWithServerAggregation.
// It is typically called via (*HLQuery).ToQueryPlanWithServerAggregation.
func NewQueryPlanWithServerAggregation(aggrLabel string, bucketedCQLQueries map[TimeInterval]CQLQuery) (*QueryPlanWithServerAggregation, error) {
	qp := &QueryPlanWithServerAggregation{
		AggregatorLabel:    aggrLabel,
		BucketedCQLQueries: bucketedCQLQueries,
	}
	return qp, nil
}

// Execute runs all CQLQueries in the QueryPlan and collects the results.
//
// TODO(rw): support parallel execution.
func (qp *QueryPlanWithServerAggregation) Execute(session *gocql.Session, debug int) ([]CQLResult, error) {
	// sort the time interval buckets we'll use:
	//fmt.Println("start one query")
	agg, err := GetAggregator(qp.AggregatorLabel)
	if err != nil {
		return nil, err
	}
	results := make([]CQLResult, 0, len(qp.BucketedCQLQueries))
	for k, q := range qp.BucketedCQLQueries {
		// Execute one CQLQuery and collect its result
		//
		// For server-side aggregation, this will return only
		// one row; for exclusive client-side aggregation this
		// will return a sequence.

		// if strings.Contains(q.PreparableQueryString, "max") {

		// 	fmt.Println(q.PreparableQueryString)
		// 	var usage_user float64
		// 	if err := session.Query(q.PreparableQueryString, q.Args...).Scan(&usage_user); err != nil {
		// 		log.Fatalf("Query failed: %v", err)
		// 	}
		// 	fmt.Println(usage_user)
		// 	agg.Put(usage_user)
		// } else {
		// 	//fmt.Println(q.PreparableQueryString)
		// 	iter := session.Query(q.PreparableQueryString, q.Args...).Iter()
		// 	var usage_user float64
		// 	iter.Scan(&usage_user)
		// 	fmt.Println(usage_user)
		// 	for iter.Scan(&usage_user) {
		// 		//fmt.Println(usage_user)
		// 		agg.Put(usage_user)
		// 	}
		// }

		fmt.Println(q.PreparableQueryString)
		iter := session.Query(q.PreparableQueryString, q.Args...).Iter()
		var usage_user float64
		iter.Scan(&usage_user)
		fmt.Println(usage_user)
		for iter.Scan(&usage_user) {
			//fmt.Println(usage_user)
			agg.Put(usage_user)
		}

		// if err := iter.Close(); err != nil {
		// 	return nil, err
		// }
		results = append(results, CQLResult{TimeInterval: k, Value: agg.Get()})
	}

	return results, nil
}

// DebugQueries prints debugging information.
func (qp *QueryPlanWithServerAggregation) DebugQueries(level int) {
	if level >= 1 {
		fmt.Printf("[qpsa] query with server aggregation plan has %d CQLQuery objects\n", len(qp.BucketedCQLQueries))
	}

	if level >= 2 {
		for k, q := range qp.BucketedCQLQueries {
			fmt.Printf("[qpsa] CQL:  %s, %s\n", k, q)
		}
	}
}

// A QueryPlanWithoutServerAggregation fulfills an HLQuery by performing
// table scans on the server and aggregating all data on the client. This
// results in higher bandwidth usage but fewer round-trip requests.
//
// It has 1) a map of Aggregators (one for each time bucket) which merge data
// on the client, 2) a GroupByDuration, which is used to reconstruct time
// buckets from a server response, 3) a set of TimeBuckets, which are used to
// store final aggregated items, and 4) a set of CQLQueries used to fulfill
// this plan.
type QueryPlanWithoutServerAggregation struct {
	Aggregators     map[TimeInterval]Aggregator
	GroupByDuration time.Duration
	TimeBuckets     []TimeInterval
	CQLQueries      []CQLQuery
}

// NewQueryPlanWithoutServerAggregation builds a QueryPlanWithoutServerAggregation.
// It is typically called via (*HLQuery).ToQueryPlanWithoutServerAggregation.
func NewQueryPlanWithoutServerAggregation(aggrLabel string, groupByDuration time.Duration, timeBuckets []TimeInterval, cqlQueries []CQLQuery) (*QueryPlanWithoutServerAggregation, error) {
	aggrs := make(map[TimeInterval]Aggregator, len(timeBuckets))
	for _, ti := range timeBuckets {
		aggr, err := GetAggregator(aggrLabel)
		if err != nil {
			return nil, err
		}

		aggrs[ti] = aggr
	}

	qp := &QueryPlanWithoutServerAggregation{
		Aggregators:     aggrs,
		GroupByDuration: groupByDuration,
		TimeBuckets:     timeBuckets,
		CQLQueries:      cqlQueries,
	}
	return qp, nil
}

// Execute runs all CQLQueries in the QueryPlan and collects the results.
//
// TODO(rw): support parallel execution.
func (qp *QueryPlanWithoutServerAggregation) Execute(session *gocql.Session, debug int) ([]CQLResult, error) {
	// for each query, execute it, then put each result row into the
	// client-side aggregator that matches its time bucket:
	for _, q := range qp.CQLQueries {
		cq := session.Query(q.PreparableQueryString, q.Args...)
		if debug == 1 {
			fmt.Printf("[qp] Query: %s\n", cq)
		}
		iter := cq.Iter()

		var timestamp_ns int64
		var value float64

		for iter.Scan(&timestamp_ns, &value) {
			ts := time.Unix(0, timestamp_ns).UTC()
			tsTruncated := ts.Truncate(qp.GroupByDuration)
			bucketKey := TimeInterval{
				Start: tsTruncated,
				End:   tsTruncated.Add(qp.GroupByDuration),
			}

			qp.Aggregators[bucketKey].Put(value)
		}
		if err := iter.Close(); err != nil {
			return nil, err
		}
	}

	// perform client-side aggregation across all buckets:
	results := make([]CQLResult, 0, len(qp.TimeBuckets))
	for _, ti := range qp.TimeBuckets {
		acc := qp.Aggregators[ti].Get()
		results = append(results, CQLResult{TimeInterval: ti, Value: acc})
	}

	return results, nil
}

// DebugQueries prints debugging information.
func (qp *QueryPlanWithoutServerAggregation) DebugQueries(level int) {
	if level >= 1 {
		fmt.Printf("[qpca] query with client aggregation plan has %d CQLQuery objects\n", len(qp.CQLQueries))
	}

	if level >= 2 {
		for i, q := range qp.CQLQueries {
			fmt.Printf("[qpca] CQL: %d, %s\n", i, q)
		}
	}
}
