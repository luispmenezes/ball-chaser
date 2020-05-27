package model

type PlayerStat struct {
	Name     string
	Plaform  string
	OnlineId uint
	IsBot    bool
	Team     uint
	Assists  uint
	Saves    uint
	Shots    uint
	Goals    uint
	Score    uint
}
