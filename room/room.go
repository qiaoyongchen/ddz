package room

import (
	"ddz/table"
)

// Room Room
type Room struct {
	tables []table.ITable
}

// NewRoom NewRoom
func NewRoom(tables []table.ITable) *Room {
	return &Room{
		tables: tables,
	}
}

// Tables Tables
func (p *Room) Tables() []table.ITable {
	return p.tables
}
