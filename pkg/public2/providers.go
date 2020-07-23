// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package public

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/HuguesGuilleus/go-parsersa"
	"net/http"
	"time"
)

// The information about the providers,
type Provider struct {
	PrivKey *rsa.PrivateKey
	Pub     []byte // the PEM public key encoded
}

// Create a new app. id is the the id of the this application in the provider.
// serv is the url of this provider.
func NewProvider(keyFile string) (*Provider, error) {
	k, err := parsersa.GenPrivKey(keyFile, 2048)
	if err != nil {
		return nil, err
	}

	buff := &bytes.Buffer{}
	pem.Encode(buff, &pem.Block{
		Type:  "BEGIN PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&k.PublicKey),
	})

	return &Provider{
		PrivKey: k,
		Pub:     buff.Bytes(),
	}, nil
}

func (p *Provider) PubHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-pem-file")
	w.Write(p.Pub)
}

const (
	// The header `{"alg": "RS256","typ": "JWT"}` already encoded + dot
	jwtHead = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9."
	// The duration of a day
	day = 24 * time.Hour
)

// Create a JWT for the user u to a specific audience.
func (p *Provider) CreateJWT(u *UserInfo, audience string) (string, error) {
	j, err := json.Marshal(jwtBody{
		UserInfo:       *u,
		Audience:       audience,
		ExpirationTime: time.Now().Add(day).Unix(),
	})
	if err != nil {
		return "", err
	}

	body := jwtHead + tob64(j)
	hashed := sha256.Sum256([]byte(body))

	sig, err := rsa.SignPKCS1v15(rand.Reader, p.PrivKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return body + "." + tob64(sig), nil
}
