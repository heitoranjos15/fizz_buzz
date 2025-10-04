package controller

import (
	"github.com/gin-gonic/gin"
)

type FizzBuzzRequest struct {
	Multiples []int    `json:"multiples"`
	Words     []string `json:"words"`
	Limit     int      `json:"limit"`
}

type FizzBuzzController struct {
	core Core
}

type Core interface {
	ParseMessage(numbers []int, values []string, limit int) (string, error)
}

func NewFizzBuzzController(core Core) *FizzBuzzController {
	return &FizzBuzzController{
		core: core,
	}
}

func (f FizzBuzzController) FizzBuzz(ctx *gin.Context) {
	var req FizzBuzzRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "params multiples, words, limit are required and must be valid"})
		return
	}

	if len(req.Multiples) != len(req.Words) {
		ctx.JSON(400, gin.H{"error": "numbers and value arrays must have the same length"})
		return
	}

	result, err := f.core.ParseMessage(req.Multiples, req.Words, req.Limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, result)
}
