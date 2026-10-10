package bedrock

import (
	"github.com/sandertv/gophertunnel/minecraft"
)

type Session struct {
	conn *minecraft.Conn

	// Local player
	player *Player

	// Connected players
	players map[string]*Player

	// Event handlers
	packetHandlers      []func(packet any)
	playerHandlers      []func(player *Player)
	disconnectHandlers  []func(err error)
}


func (s *Session) Player() *Player {
	return s.player
}


func (s *Session) Players() []*Player {
	return nil
}


func (s *Session) OnPacket(handler func(packet any)) {
}

func (s *Session) OnPlayer(handler func(player *Player)) {
}

func (s *Session) OnDisconnect(handler func(err error)) {
}



func (s *Session) Command() *Command {
	return Command{}
}



func (s *Session) handlePacket(packet any) {
}

func (s *Session) handlePlayerJoin(player *Player) {
}

func (s *Session) handleDisconnect(err error) {
}