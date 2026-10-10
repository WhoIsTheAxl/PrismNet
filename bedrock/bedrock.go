// Package bedrock provides client, server, and proxy abstractions for Minecraft Bedrock.
package bedrock

import (
 "context"
 "errors"
 "github.com/sandertv/gophertunnel/minecraft"
 "golang.org/x/oauth2"
)

// Realm identifies an invite-based Minecraft Realm.
type Realm struct {
 Code string
 Auth oauth2.TokenSource
}
func NewRealm(code string, source oauth2.TokenSource) *Realm { return &Realm{Code: code, Auth: source} }

// Connect dials a standard Bedrock server and completes spawning.
func Connect(ctx context.Context, addr string, source oauth2.TokenSource) (*Client, error) {
 if source == nil { return nil, errors.New("token source is required") }
 conn, err := serverConnect(ctx, addr, source)
 if err != nil { return nil, err }
 return newClient(conn), nil
}
func (r *Realm) Connect(ctx context.Context) (*Client, error) {
 if r == nil || r.Auth == nil || r.Code == "" { return nil, errors.New("realm code and auth required") }
 conn, err := realmConnect(ctx, r.Auth, r.Code)
 if err != nil { return nil, err }
 return newClient(conn), nil
}
// WrapConnection adopts an already-spawned Gophertunnel connection.
func WrapConnection(conn *minecraft.Conn) *Client { if conn == nil { return nil }; return newClient(conn) }
