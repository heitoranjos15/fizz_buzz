package controller

import (
	"github.com/gin-gonic/gin"
)

type FizzBuzzRequest struct {
	Numbers []int    `json:"numbers" binding:"required,dive,gt=0"`
	Value   []string `json:"value" binding:"required,dive,notblank"`
	Limit   int      `json:"limit" binding:"required,gt=0"`
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
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(req.Numbers) != len(req.Value) {
		ctx.JSON(400, gin.H{"error": "numbers and value arrays must have the same length"})
		return
	}

	result, err := f.core.ParseMessage(req.Numbers, req.Value, req.Limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, result)
}
