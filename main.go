package main

import (
	"fizzbuzz/internal/controller"
	"fizzbuzz/internal/core"

	"github.com/gin-gonic/gin"
)

type FizzBuzzRequest struct {
	Numbers []int    `json:"numbers" binding:"required,dive,gt=0"`
	Value   []string `json:"value" binding:"required,dive,notblank"`
	Limit   int      `json:"limit" binding:"required,gt=0"`
}

func main() {
	core := core.NewCore()
	controller := controller.NewFizzBuzzController(core)

	r := gin.Default()

	r.POST("/fizzbuzz", controller.FizzBuzz)

	r.Run(":8067")
}
