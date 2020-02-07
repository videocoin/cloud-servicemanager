// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jws"
	"golang.org/x/oauth2/jwt"
)

// ParseKey converts the binary contents of a private key file
// to an *rsa.PrivateKey. It detects whether the private key is in a
// PEM container or not. If so, it extracts the the private key
// from PEM container before conversion. It only supports PEM
// containers with no passphrase.
func ParseKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("private key should be a PEM or plain PKCS1 or PKCS8; parse error: %v", err)
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is invalid")
	}
	return parsed, nil
}

// JWTAccessTokenSourceFromJSON uses a VideoCoin service account JSON
// key file to read the credentials that authorize and authenticate the
// requests, and returns a TokenSource that does not use any OAuth2 flow but
// instead creates a JWT and sends that as the access token.
// The audience is typically a URL that specifies the scope of the credentials.
//
// Note that this is not a standard OAuth flow, but rather an
// optimization supported by a few VideoCoin services.
// Unless you know otherwise, you should use JWTConfigFromJSON instead.
func JWTAccessTokenSourceFromJSON(jsonKey []byte, audience string) (oauth2.TokenSource, error) {
	cfg, err := JWTConfigFromJSON(jsonKey)
	if err != nil {
		return nil, fmt.Errorf("videocoin: could not parse JSON key: %v", err)
	}
	pk, err := ParseKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("videocoin: could not parse key: %v", err)
	}
	ts := &jwtAccessTokenSource{
		email:    cfg.Email,
		audience: audience,
		pk:       pk,
		pkID:     cfg.PrivateKeyID,
	}
	tok, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return oauth2.ReuseTokenSource(tok, ts), nil
}

// JWTAccessTokenSourceFromConfig uses a JWT configuration to read the
// credentials that authorize and authenticate the requests, and returns a
// TokenSource that does not use any OAuth2 flow but instead creates a JWT and
// sends that as the access token.
// The audience is typically a URL that specifies the scope of the credentials.
//
// Note that this is not a standard OAuth flow, but rather an
// optimization supported by a few VideoCoin services.
// Unless you know otherwise, you should use JWTConfigFromJSON instead.
func JWTAccessTokenSourceFromConfig(cfg *jwt.Config, audience string) (oauth2.TokenSource, error) {
	pk, err := ParseKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("videocoin: could not parse key: %v", err)
	}
	ts := &jwtAccessTokenSource{
		email:    cfg.Email,
		audience: audience,
		pk:       pk,
		pkID:     cfg.PrivateKeyID,
	}
	tok, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return oauth2.ReuseTokenSource(tok, ts), nil
}

type jwtAccessTokenSource struct {
	email, audience string
	pk              *rsa.PrivateKey
	pkID            string
}

func (ts *jwtAccessTokenSource) Token() (*oauth2.Token, error) {
	iat := time.Now()
	exp := iat.Add(time.Hour)
	cs := &jws.ClaimSet{
		Iss: ts.email,
		Sub: ts.email,
		Aud: ts.audience,
		Iat: iat.Unix(),
		Exp: exp.Unix(),
	}
	hdr := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
		KeyID:     string(ts.pkID),
	}
	msg, err := jws.Encode(hdr, cs, ts.pk)
	if err != nil {
		return nil, fmt.Errorf("videocoin: could not encode JWT: %v", err)
	}
	return &oauth2.Token{AccessToken: msg, TokenType: "Bearer", Expiry: exp}, nil
}
