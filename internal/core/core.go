package core

import (
	"fmt"
)

type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) ParseMessage(words []int, multiples []string, limit int) (string, error) {
	result := []string{}
	for i := 1; i <= limit; i++ {
		output := ""
		isMultiple := false
		for j, num := range words {
			if i%num == 0 {
				output += multiples[j]
				isMultiple = true
			}
		}
		if !isMultiple {
			output = fmt.Sprintf("%d", i)
		}
		result = append(result, output)
	}

	return fmt.Sprintf("%v", result), nil
}
