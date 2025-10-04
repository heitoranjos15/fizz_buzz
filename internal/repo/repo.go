package repo

import (
	"fmt"

	"gorm.io/gorm"
)

type Repo[T any] struct {
	tableName string
	conn      *gorm.DB
}

func NewRepo[T *gorm.DB](tableName string, db T) *Repo[T] {
	return &Repo[T]{tableName: tableName, conn: db}
}

func (r *Repo[T]) SaveMessage(words []string, multiples []int, limit int) error {
	type FizzBuzzRecord struct {
		gorm.Model
		Words     string
		Multiples string
		Limit     int
	}

	record := FizzBuzzRecord{
		Words:     fmt.Sprint(words),
		Multiples: fmt.Sprint(multiples),
		Limit:     limit,
	}

	if err := r.conn.Table(r.tableName).Create(&record).Error; err != nil {
		return err
	}
	return nil
}
