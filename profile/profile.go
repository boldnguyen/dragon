package profile

import (
	"fmt"
)

type Profile struct {
	Name     string
	Level    int
	PlayerID string
	Avatar   string
	Wins     int
	Losses   int
	Dragons  []Dragon
	Tokens   float64
}

type Dragon struct {
	ID    string
	Name  string
	Power int
	Level int
}

func NewProfile(name string, level int) *Profile {
	return &Profile{
		Name:     name,
		Level:    level,
		PlayerID: "Player_" + name,
	}
}

func (p *Profile) SetLevel(level int) {
	p.Level = level
}

func (p *Profile) ChangeName(newName string, cost float64) error {
	if p.Tokens < cost {
		return fmt.Errorf("not enough tokens to change name")
	}
	p.Tokens -= cost
	p.Name = newName
	return nil
}

func (p *Profile) ChangeAvatar(newAvatar string, cost float64) error {
	if p.Tokens < cost {
		return fmt.Errorf("not enough tokens to change avatar")
	}
	p.Tokens -= cost
	p.Avatar = newAvatar
	return nil
}

func (p *Profile) String() string {
	return fmt.Sprintf("Name: %s, Level: %d, Tokens: %.2f, Dragons: %d", p.Name, p.Level, p.Tokens, len(p.Dragons))
}
