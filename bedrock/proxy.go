package bedrock

import (
 "context"
 "errors"
 "sync"
 "github.com/sandertv/gophertunnel/minecraft"
 "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Direction identifies a proxy packet's route.
type Direction uint8
const (
 ClientToServer Direction = iota
 ServerToClient
)
type ProxyPacketHandler func(Direction, packet.Packet) (packet.Packet, bool)

// ProxySession relays packets between an accepted downstream connection and
// an already established upstream connection. Connections must be spawned first.
// Hooks may replace or drop packets; hooks execute on the forwarding goroutine.
type ProxySession struct {
 Downstream *minecraft.Conn
 Upstream *minecraft.Conn
 mu sync.RWMutex
 hooks []ProxyPacketHandler
 closeOnce sync.Once
}
func NewProxySession(downstream,upstream *minecraft.Conn) (*ProxySession,error) {
 if downstream==nil||upstream==nil {return nil,errors.New("both connections required")}
 return &ProxySession{Downstream:downstream,Upstream:upstream},nil
}
func (p *ProxySession) OnPacket(h ProxyPacketHandler) {
 if h==nil{return}
 p.mu.Lock();p.hooks=append(p.hooks,h);p.mu.Unlock()
}
func (p *ProxySession) forward(from,to *minecraft.Conn, direction Direction) error {
 for {
  pkt,err:=from.ReadPacket()
  if err!=nil{return err}
  p.mu.RLock(); hooks:=append([]ProxyPacketHandler(nil),p.hooks...);p.mu.RUnlock()
  forward:=true
  for _,h:=range hooks {
   pkt,forward=h(direction,pkt)
   if !forward||pkt==nil {break}
  }
  if !forward||pkt==nil {continue}
  if err:=to.WritePacket(pkt);err!=nil{return err}
 }
}
func (p *ProxySession) Close() {
 p.closeOnce.Do(func(){_ = p.Downstream.Close();_ = p.Upstream.Close()})
}
// Run relays both directions until cancellation or either connection fails.
func (p *ProxySession) Run(ctx context.Context) error {
 if p==nil {return errors.New("nil proxy session")}
 errorsCh:=make(chan error,2)
 go func(){errorsCh<-p.forward(p.Downstream,p.Upstream,ClientToServer)}()
 go func(){errorsCh<-p.forward(p.Upstream,p.Downstream,ServerToClient)}()
 select {
 case err:=<-errorsCh:p.Close();return err
 case <-ctx.Done():p.Close();return ctx.Err()
 }
}
