package core

import (
	"fmt"
)

type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) ParseMessage(numbers []int, values []string, limit int) (string, error) {
	result := []string{}
	for i := 1; i <= limit; i++ {
		output := ""
		isMultiple := false
		for j, num := range numbers {
			if i%num == 0 {
				output += values[j]
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
