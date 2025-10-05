package controller

import (
	"errors"
	"fizzbuzz/internal/types"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var errorStatsNotFound = errors.New("no stats found")

type FizzBuzzRequest struct {
	Multiples []int    `json:"multiples"`
	Words     []string `json:"words"`
	Limit     int      `json:"limit"`
}

type StatResp struct {
	TotalRequests int                      `json:"total_requests"`
	RequestStats  []types.StatsParameters  `json:"request_stats"`
	WordsStats    []types.StatsWordsResult `json:"words_stats"`
}

type FizzBuzzController struct {
	core Core
}

type Core interface {
	ProcessMessage(words []string, multiples []int, limit int) (string, error)
	GetStatsParameters() ([]types.StatsParameters, error)
	GetStatsWords() ([]types.StatsWordsResult, error)
	GetTotalRequests() (int, error)
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
		return
	}

	err = req.validate()
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := f.core.ProcessMessage(req.Words, req.Multiples, req.Limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"result": result})
}

func (f FizzBuzzController) Stats(ctx *gin.Context) {
	stats, err := f.core.GetStatsParameters()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	totalRequests, err := f.core.GetTotalRequests()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	wordsStats, err := f.core.GetStatsWords()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	statsResp := StatResp{
		TotalRequests: totalRequests,
		RequestStats:  stats,
		WordsStats:    wordsStats,
	}

	ctx.JSON(200, statsResp)
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
