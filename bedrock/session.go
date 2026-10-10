package bedrock

import (
 "errors"
 "sync"
 "github.com/sandertv/gophertunnel/minecraft"
 "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)
type PacketHandler func(packet.Packet)
type Session struct {
 conn *minecraft.Conn
 mu sync.RWMutex
 player *Player
 players map[string]*Player
 packetHandlers []PacketHandler
 playerHandlers []func(*Player)
 disconnectHandlers []func(error)
 started bool
 done chan struct{}
 once sync.Once
}
func newSession(conn *minecraft.Conn) *Session {
 s := &Session{conn:conn,players:make(map[string]*Player),done:make(chan struct{})}
 s.player = &Player{session:s,conn:conn}
 return s
}
func (s *Session) Player() *Player { s.mu.RLock(); defer s.mu.RUnlock(); return s.player }
func (s *Session) Players() []*Player {
 s.mu.RLock(); defer s.mu.RUnlock()
 out := make([]*Player,0,len(s.players))
 for _,p := range s.players { out=append(out,p) }
 return out
}
func (s *Session) OnPacket(h PacketHandler) { if h==nil{return};s.mu.Lock();s.packetHandlers=append(s.packetHandlers,h);s.mu.Unlock() }
func (s *Session) OnPlayer(h func(*Player)) { if h==nil{return};s.mu.Lock();s.playerHandlers=append(s.playerHandlers,h);s.mu.Unlock() }
func (s *Session) OnDisconnect(h func(error)) { if h==nil{return};s.mu.Lock();s.disconnectHandlers=append(s.disconnectHandlers,h);s.mu.Unlock() }
func (s *Session) Command() *Command { return &Command{session:s} }
func (s *Session) WritePacket(p packet.Packet) error { if s==nil||s.conn==nil{return errors.New("session closed")};return s.conn.WritePacket(p) }
func (s *Session) Done() <-chan struct{} { return s.done }
// Start begins the single packet reader. Register handlers before calling Start.
func (s *Session) Start() error {
 s.mu.Lock()
 if s.started { s.mu.Unlock(); return errors.New("session already started") }
 s.started=true
 s.mu.Unlock()
 go s.readLoop()
 return nil
}
func (s *Session) readLoop() {
 var failure error
 defer func(){s.handleDisconnect(failure);s.stop()}()
 for {
  p,err:=s.conn.ReadPacket()
  if err!=nil { failure=err;return }
  s.handlePacket(p)
 }
}
func (s *Session) stop() { s.once.Do(func(){close(s.done)}) }
func (s *Session) handlePacket(p packet.Packet) {
 s.mu.RLock(); handlers:=append([]PacketHandler(nil),s.packetHandlers...);s.mu.RUnlock()
 for _,h:=range handlers {h(p)}
}
func (s *Session) handlePlayerJoin(p *Player) {
 s.mu.Lock();s.players[p.UUID.String()]=p;handlers:=append([]func(*Player){}, s.playerHandlers...);s.mu.Unlock()
 for _,h:=range handlers {h(p)}
}
func (s *Session) handleDisconnect(err error) {
 s.mu.RLock();handlers:=append([]func(error){}, s.disconnectHandlers...);s.mu.RUnlock()
 for _,h:=range handlers {h(err)}
}
