package oauth2

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2/jwt"
)

const (
	// JWTTokenURL is VideoCoin's OAuth 2.0 token URL to use with the JWT flow.
	JWTTokenURL = "https://videocoin.network/oauth2/token"

	// JSON key file types.
	serviceAccountKey = "service_account"
)

// credentialsFile is the unmarshalled representation of a credentials file.
type credentialsFile struct {
	Type string `json:"type"` // serviceAccountKey

	// Service Account fields
	ClientEmail  string `json:"client_email"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	TokenURL     string `json:"token_uri"`
}

func (f *credentialsFile) jwtConfig(scopes []string) *jwt.Config {
	cfg := &jwt.Config{
		Email:        f.ClientEmail,
		PrivateKey:   []byte(f.PrivateKey),
		PrivateKeyID: f.PrivateKeyID,
		Scopes:       scopes,
		TokenURL:     f.TokenURL,
	}
	if cfg.TokenURL == "" {
		cfg.TokenURL = JWTTokenURL
	}
	return cfg
}

// JWTConfigFromJSON uses a VideoCoin service account JSON key file to read
// the credentials that authorize and authenticate the requests.
func JWTConfigFromJSON(jsonKey []byte, scope ...string) (*jwt.Config, error) {
	var f credentialsFile
	if err := json.Unmarshal(jsonKey, &f); err != nil {
		return nil, err
	}
	if f.Type != serviceAccountKey {
		return nil, fmt.Errorf("videocoin: read JWT from JSON credentials: 'type' field is %q (expected %q)", f.Type, serviceAccountKey)
	}
	scope = append([]string(nil), scope...) // copy
	return f.jwtConfig(scope), nil
}
