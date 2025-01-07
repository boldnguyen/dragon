package models

type Player struct {
	ID      string
	Name    string
	Level   int
	Wins    int
	Losses  int
	Dragons []Dragon
	Items   []Item
	Tokens  float64
}
