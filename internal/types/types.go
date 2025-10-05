package types

type StatsWordsResult struct {
	Word  string `bson:"word" json:"word"`
	Total int    `bson:"count" json:"total"`
}

type StatsParameters struct {
	Words         []string `bson:"words" json:"words"`
	Multiples     []int    `bson:"multiples" json:"multiples"`
	Limit         int      `bson:"limit" json:"limit"`
	TotalRequests int      `bson:"count" json:"total_requests"`
}
