package bedrock

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
)

type Position struct {
	X float32
	Y float32
	Z float32
}

type Player struct {
	// Identity
	Name string
	UUID uuid.UUID

	// State
	Position Position

	// Internal
	session *Session
	conn    *minecraft.Conn
}



func (p *Player) Move(position Position) {
}

func (p *Player) LookAt(position Position) {
}



func (p *Player) Punch() {
}

func (p *Player) StartBreak() {
}

func (p *Player) StopBreak() {
}

// func (p *Player) Jump() {}
// func (p *Player) Sneak() {}
// func (p *Player) Sprint() {}
// func (p *Player) UseItem() {}
// func (p *Player) Interact() {}
// func (p *Player) PlaceBlock(...) {}