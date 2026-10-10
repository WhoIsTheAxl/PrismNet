package bedrock

import (
 "sync"
 "github.com/sandertv/gophertunnel/minecraft"
)
type Client struct {
 conn *minecraft.Conn
 once sync.Once
 session *Session
}
func newClient(conn *minecraft.Conn) *Client {
 c := &Client{conn:conn}
 c.session = newSession(conn)
 return c
}
func (c *Client) Close() error {
 if c == nil { return nil }
 var err error
 c.once.Do(func(){ c.session.stop(); err = c.conn.Close() })
 return err
}
func (c *Client) Session() *Session { if c == nil { return nil }; return c.session }
func (c *Client) Conn() *minecraft.Conn { if c == nil { return nil }; return c.conn }
