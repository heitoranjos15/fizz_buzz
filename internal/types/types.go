package types

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type StatsWordsResult struct {
	Word  string `bson:"word"`
	Count int    `bson:"count"`
}

func (stats *StatsWordsResult) FromBson(m bson.M) {
	if count, ok := m["count"].(int32); ok {
		stats.Count = int(count)
	}
	log.Printf("Converted BSON to StatsByKeyResult: %+v", stats)
}

type StatsParameters struct {
	Words         []string `bson:"words" json:"words"`
	Multiples     []int    `bson:"multiples" json:"multiples"`
	Limit         int      `bson:"limit" json:"limit"`
	TotalRequests int      `bson:"count" json:"total_requests"`
}
