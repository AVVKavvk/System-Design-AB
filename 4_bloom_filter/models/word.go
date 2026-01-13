package models

type Word struct {
	Word string `json:"word"`
}
type ResponseWordProbability struct {
	IsFound bool `json:"isFound"`
	RowIdx  int  `json:"rowIdx"`
	ColIdx  int  `json:"colIdx"`
}
