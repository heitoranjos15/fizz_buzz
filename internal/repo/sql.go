package repo

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Repo struct {
	tableName string
	conn      *gorm.DB
}

func NewRepo(tableName string, db *gorm.DB) *Repo {
	return &Repo{tableName: tableName, conn: db}
}

func (r *Repo) SaveMessage(words []string, multiples []int, limit int) error {
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

func (r *Repo) GetStatsByKey() (map[string]any, error) {
	return map[string]any{}, errors.New("not implemented")
}

func (r *Repo) GetTotalRequests() (int, error) {
	var count int64
	if err := r.conn.Table(r.tableName).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *Repo) GetStatsParameters() ([]map[string]any, error) {
	type Result struct {
		Words     string
		Multiples string
		Limit     int
		Count     int64
	}

	var results []Result
	if err := r.conn.Table(r.tableName).
		Select("words, multiples, limit, COUNT(*) as count").
		Group("words, multiples, limit").
		Order("count DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	var stats []map[string]any
	for _, res := range results {
		stats = append(stats, map[string]any{
			"words":     res.Words,
			"multiples": res.Multiples,
			"limit":     res.Limit,
			"count":     res.Count,
		})
	}

	return stats, nil
}
