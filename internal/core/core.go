package core

import (
	"errors"
	"fmt"
	"log"
)

type repo[T any] interface {
	SaveMessage(words []string, multiples []int, limit int) error
	GetStats() (map[string]any, error)
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

type Stats struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type StatResp struct {
	TotalRequests int     `json:"total_requests"`
	Stats         []Stats `json:"stats"`
}

func (c *Core[T]) GetStats() (StatResp, error) {
	dbResp, err := c.db.GetStats()
	if err != nil {
		return StatResp{}, err
	}
	// check dbResp has keys "words" and "count"
	if dbResp == nil {
		return StatResp{}, errors.New("no stats found")
	}
	wordsArray, ok := dbResp["words"].(any)
	if !ok {
		return StatResp{}, errors.New("invalid stats format")
	}
	count, ok := dbResp["count"].(int32)
	if !ok {
		return StatResp{}, errors.New("invalid stats format")
	}
	stats := StatResp{
		TotalRequests: int(count),
		Stats: []Stats{
			{Word: fmt.Sprintf("%v", wordsArray), Count: int(count)},
		},
	}

	return stats, nil
}
