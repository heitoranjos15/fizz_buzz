package repo

import (
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
