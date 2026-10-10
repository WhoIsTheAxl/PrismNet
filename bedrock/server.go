package bedrock

import (
	"github.com/sandertv/gophertunnel/minecraft"
)

func serverConnect(
    ctx context.Context,
    address string,
    auth *oauth2.Token,
) (*Client, error) {

    dialer := minecraft.Dialer{
        TokenSource: auth,
    }

    conn, err := dialer.DialContext(
        ctx,
        "raknet",
        address,
    )
    if err != nil {
        return nil, err
    }

    if err := conn.DoSpawn(); err != nil {
        conn.Close()
        return nil, err
    }

    return conn, nil
}