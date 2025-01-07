package models

type Asset struct {
	ID    string
	Name  string
	Type  string // "token", "dragon", "item"
	Value float64
	Owner string // Player ID
}

// Token is a specialized asset
type Token struct {
	Asset
	Amount float64
}

// Dragon extends the asset for dragons
type Dragon struct {
	Asset
	Power int
	Level int
}

// Item extends the asset for items
type Item struct {
	Asset
	ItemType string // e.g., "weapon", "armor"
}
