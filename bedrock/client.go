package bedrock

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/WhoIsTheAxl/Prismnet/bedrock"
)

type Client struct {
	conn *minecraft.Conn
}


func newClient(Conn *minecraft.Conn) *Client {
	return Client{
		conn: Conn
	}
}

func (c *Client) Close() error {
    return c.conn.Close()
}

func (c *Client) Session() *Session {
	return Session{}
}