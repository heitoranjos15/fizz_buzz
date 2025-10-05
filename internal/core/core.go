package core

import (
	"errors"
	"fizzbuzz/internal/types"
	"fmt"
	"log"
)

type repo[T any] interface {
	SaveMessage(words []string, multiples []int, limit int) error
	GetStatsParameters() ([]types.StatsParameters, error)
	GetStatsWords() ([]types.StatsWordsResult, error)
	GetTotalRequests() (int, error)
}

type Core[T any] struct {
	db repo[T]
}

func NewCore[T any](db repo[T]) *Core[T] {
	return &Core[T]{db: db}
}

func (c *Core[T]) ProcessMessage(words []string, multiples []int, limit int) (string, error) {
	message := c.parseMessage(words, multiples, limit)

	if err := c.db.SaveMessage(words, multiples, limit); err != nil {
		log.Printf("Error saving message to database: %v", err)
		return message, errors.New("could not save message to database")
	}

	return message, nil
}

func (c *Core[T]) parseMessage(words []string, multiples []int, limit int) string {
	result := []string{}
	for i := 1; i <= limit; i++ {
		output := ""
		isMultiple := false
		for j, num := range multiples {
			if i%num == 0 {
				output += words[j]
				isMultiple = true
			}
		}
		if !isMultiple {
			output = fmt.Sprintf("%d", i)
		}
		result = append(result, output)
	}

	return fmt.Sprintf("%v", result)
}

func (c *Core[T]) GetStatsParameters() ([]types.StatsParameters, error) {
	stats, err := c.db.GetStatsParameters()
	if err != nil {
		log.Printf("Error retrieving stats from database: %v", err)
		return stats, errors.New("could not retrieve stats from database")
	}
	return stats, nil
}

func (c *Core[T]) GetStatsWords() ([]types.StatsWordsResult, error) {
	stats, err := c.db.GetStatsWords()
	if err != nil {
		log.Printf("Error retrieving words stats from database: %v", err)
		return stats, errors.New(fmt.Sprintf("could not retrieve words stats from database"))
	}
	return stats, nil
}

func (c *Core[T]) GetTotalRequests() (int, error) {
	total, err := c.db.GetTotalRequests()
	if err != nil {
		log.Printf("Error retrieving total requests from database: %v", err)
		return 0, errors.New("could not retrieve total requests from database")
	}
	return total, nil
}
