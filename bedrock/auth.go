package bedrock

import (
 "encoding/json"
 "errors"
 "os"
 "path/filepath"
 "github.com/sandertv/gophertunnel/minecraft/auth"
 "golang.org/x/oauth2"
)

// AuthSource wraps an OAuth token source used by Bedrock connections.
type AuthSource struct { TokenSource oauth2.TokenSource }
type AuthConfig struct { TokenPath string }

// Auth loads a cached Microsoft token or requests a new one.
// Never commit the resulting token file to source control.
func Auth(config AuthConfig) (AuthSource, error) {
 if config.TokenPath == "" { return AuthSource{}, errors.New("token path is required") }
 var token *oauth2.Token
 data, err := os.ReadFile(config.TokenPath)
 if err == nil {
  if err = json.Unmarshal(data, &token); err != nil { return AuthSource{}, err }
 } else if !os.IsNotExist(err) { return AuthSource{}, err }
 if token == nil || !token.Valid() {
  token, err = auth.RequestLiveToken()
  if err != nil { return AuthSource{}, err }
  data, err = json.MarshalIndent(token, "", "  ")
  if err != nil { return AuthSource{}, err }
  if err = os.MkdirAll(filepath.Dir(config.TokenPath), 0700); err != nil { return AuthSource{}, err }
  if err = os.WriteFile(config.TokenPath, data, 0600); err != nil { return AuthSource{}, err }
 }
 return AuthSource{TokenSource: auth.RefreshTokenSource(token)}, nil
}
