// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package public

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Public informations, send to web page.
type UserInfo struct {
	ID     string    `json:"id"`
	Pseudo string    `json:"pseudo"`
	Email  string    `json:"email"`
	Avatar string    `json:"avatar"`
	Level  UserLevel `json:"level"`
}

// UserInfo with some information to a specific service.
type jwtBody struct {
	UserInfo

	Audience       string `json:"aud"`
	ExpirationTime int64  `json:"exp"`
}

const (
	// The header `{"alg": "RS256","typ": "JWT"}` already encoded
	jwtHead = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9."
	// The duration of a day
	day = 24 * time.Hour
)

// ToJWT create a JWT and sign it with pub for a specific audience.
func (u *UserInfo) ToJWT(priv *rsa.PrivateKey, audience string) (string, error) {
	j, err := json.Marshal(jwtBody{
		UserInfo:       *u,
		Audience:       audience,
		ExpirationTime: time.Now().Add(day).Unix(),
	})
	if err != nil {
		return "", err
	}
	body := jwtHead + base64.RawURLEncoding.EncodeToString(j)
	hashed := sha256.Sum256([]byte(body))

	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return body + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

type jwtHeaders struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

var (
	JWTWrongSyntax   = errors.New("JWT wrong syntax")
	JWTWrongHead     = errors.New("JWT wrong head")
	JWTWrongAudience = errors.New("This JWT is made for an other audience")
	JWTOutDate       = errors.New("This JWT is out date")
	JWTEmpty         = errors.New("JWT is empty")
)

// Verify the JWT and parse the JWT.
func FromJWT(j string, pub *rsa.PublicKey, audience string) (*UserInfo, error) {
	if j == "" {
		return nil, JWTEmpty
	}

	parts := strings.SplitN(j, ".", 3)
	if len(parts) != 3 {
		return nil, JWTWrongSyntax
	}

	// Check the head
	b, _ := base64.URLEncoding.DecodeString(parts[0])
	h := jwtHeaders{}
	json.Unmarshal(b, &h)
	if h.Alg != "RS256" || h.Typ != "JWT" {
		return nil, JWTWrongHead
	}

	// Check the signature
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256([]byte(strings.Join(parts[:2], ".")))
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], sig); err != nil {
		return nil, err
	}

	// Get Body
	data, _ := base64.RawURLEncoding.DecodeString(parts[1])
	var body jwtBody
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, err
	}

	// Check Body
	if body.Audience != audience {
		return nil, JWTWrongAudience
	}
	if body.ExpirationTime < time.Now().Unix() {
		return nil, JWTOutDate
	}

	return &body.UserInfo, nil
}
