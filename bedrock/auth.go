package bedrock

import (
	"encoding/json"
	"os"

	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
)

type AuthSource struct {
	TokenSource *oauth2.TokenSource
}

type AuthConfig struct {
	TokenPath string
}

func authenticate(config AuthConfig) (AuthSource, error) {
	var token *oauth2.Token

	data, err := os.ReadFile(config.TokenPath)
	if err == nil {
		if err := json.Unmarshal(data, &token); err != nil {
			return AuthSource{}, err
		}
	}
	
	src := auth.RefreshTokenSource(token)

	if token == nil || !token.Valid() {
		token, err = auth.RequestLiveToken()
		if err != nil {
			return AuthSource{}, err
		}
		
		src := auth.RefreshTokenSource(token)

		data, err := json.Marshal(token)
		if err != nil {
			return AuthSource{}, err
		}

		if err := os.WriteFile(config.TokenPath, data, 0600); err != nil {
			return AuthSource{}, err
		}
	}

	return AuthSource{Token: src}, nil
}