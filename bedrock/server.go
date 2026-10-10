package bedrock

import (
 "context"
 "github.com/sandertv/gophertunnel/minecraft"
 "golang.org/x/oauth2"
)
func serverConnect(ctx context.Context, address string, source oauth2.TokenSource) (*minecraft.Conn, error) {
 dialer := minecraft.Dialer{TokenSource: source}
 conn, err := dialer.DialContext(ctx, "raknet", address)
 if err != nil { return nil, err }
 if err := conn.DoSpawn(); err != nil { _ = conn.Close(); return nil, err }
 return conn, nil
}
