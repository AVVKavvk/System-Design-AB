package models

type ResponseWithShard struct {
	Shard int    `json:"shard"`
	Data  []User `json:"data"`
}
