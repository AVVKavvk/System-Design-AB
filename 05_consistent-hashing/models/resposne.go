package models

type ResponseModel struct {
	Users      []User `json:"users"`
	ServerName string `json:"serverName"`
	Hash       string `json:"hash"`
}
