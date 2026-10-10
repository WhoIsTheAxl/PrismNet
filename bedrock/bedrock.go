package bedrock

import (
	"context"
	"golang.org/x/oauth2"
)

type Realm struct {
	Code string
	Auth *oauth2.Token
}

func Auth(config AuthConfig) AuthSource {
	return authenticate(config)
}
func Realm(code string, auth AuthSource) Realm {
	return Realm{Code: code, Auth: auth}
}
func Connect(ctx context.Context, addr string, auth *oauth2.Token) (*Client, error) {
	conn, err := serverConnect(ctx, addr, auth)
	if err != nil {
		return nil, err
	}
	return newClient(conn), nil
}

func (r *Realm) Connect(ctx context.Context) (*Client, error) {
	conn, err := realmConnect(ctx, r.Code, r.Auth)
	if err != nil {
		return nil, err
	}
	return newClient(conn), nil
}