package controller

import (
	"errors"
	"strconv"
	"strings"

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
	ProcessMessage(words []string, multiples []int, limit int) (string, error)
	GetStats() (string, error)
}

func NewFizzBuzzController(core Core) *FizzBuzzController {
	return &FizzBuzzController{
		core: core,
	}
}

func (f FizzBuzzController) FizzBuzz(ctx *gin.Context) {
	var req FizzBuzzRequest
	err := req.loadQueryParams(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
	}

	err = req.validate()
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
	}

	result, err := f.core.ProcessMessage(req.Words, req.Multiples, req.Limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"result": result})
}

var errorStatsNotFound = errors.New("no stats found")

func (f FizzBuzzController) Stats(ctx *gin.Context) {
	stats, err := f.core.GetStats()
	if err != nil {
		if !errors.Is(err, errorStatsNotFound) {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(200, gin.H{"stats": stats})
}

func (fbr *FizzBuzzRequest) loadQueryParams(ctx *gin.Context) error {
	multiplesParam := strings.Split(ctx.Query("multiples"), ",")
	multiples := []int{}
	for _, m := range multiplesParam {
		if m != "" {
			// convert m to int and handle Error
			multip, err := strconv.Atoi(m)
			if err != nil {
				return errors.New("invalid multiple value")
			}
			multiples = append(multiples, multip)
		}
	}

	words := strings.Split(ctx.Query("words"), ",")
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		return errors.New("invalid limit value")
	}

	fbr.Multiples = multiples
	fbr.Words = words
	fbr.Limit = limit

	return nil
}

func (fbr *FizzBuzzRequest) validate() error {
	if fbr.Limit <= 0 || len(fbr.Multiples) == 0 || len(fbr.Words) == 0 {
		return errors.New("params multiples, words and limit are required and must be valid")
	}
	if len(fbr.Multiples) != len(fbr.Words) {
		return errors.New("multiples and words arrays must have the same length")
	}
	return nil
}
