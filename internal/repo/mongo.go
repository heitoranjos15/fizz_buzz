package repo

import (
	"fizzbuzz/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo[T any] struct {
	collectionName string
	conn           *mongo.Collection
}

type FizzBuzzRecord struct {
	Words     []string `bson:"words"`
	Multiples []int    `bson:"multiples"`
	Limit     int      `bson:"limit"`
}

func NewMongoRepo[T any](collectionName string, db *mongo.Collection) *MongoRepo[T] {
	return &MongoRepo[T]{collectionName: collectionName, conn: db}
}

func (r *MongoRepo[T]) SaveMessage(words []string, multiples []int, limit int) error {
	record := FizzBuzzRecord{
		Words:     words,
		Multiples: multiples,
		Limit:     limit,
	}

	_, err := r.conn.InsertOne(nil, record)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepo[T]) GetStatsParameters() ([]types.StatsParameters, error) {
	var results []types.StatsParameters

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "words", Value: "$words"},
				{Key: "multiples", Value: "$multiples"},
				{Key: "limit", Value: "$limit"},
			}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "words", Value: "$_id.words"},
			{Key: "multiples", Value: "$_id.multiples"},
			{Key: "limit", Value: "$_id.limit"},
			{Key: "count", Value: 1},
			{Key: "_id", Value: 0},
		}}},
	}

	cursor, err := r.conn.Aggregate(nil, pipeline)
	if err != nil {
		return results, err
	}
	defer cursor.Close(nil)

	if err = cursor.All(nil, &results); err != nil {
		return results, err
	}

	return results, nil
}

func (r *MongoRepo[T]) GetStatsWords() ([]types.StatsWordsResult, error) {
	var statsList []types.StatsWordsResult

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$unwind", Value: "$words"}},
		bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: "$words"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			},
		}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "word", Value: "$_id"},
			{Key: "count", Value: 1},
			{Key: "_id", Value: 1},
		}}},
	}
	cursor, err := r.conn.Aggregate(nil, pipeline)
	if err != nil {
		return statsList, err
	}
	defer cursor.Close(nil)

	if err = cursor.All(nil, &statsList); err != nil {
		return statsList, err
	}

	return statsList, nil
}

func (r *MongoRepo[T]) GetTotalRequests() (int, error) {
	count, err := r.conn.CountDocuments(nil, bson.D{})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
