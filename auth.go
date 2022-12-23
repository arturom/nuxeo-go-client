package nuxeogoclient

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

type AuthMethod interface {
	applyAuth(r *http.Request)
}

type BasicAuth struct {
	authString string
}

func (a BasicAuth) applyAuth(r *http.Request) {
	r.Header.Set("Authorization", a.authString)
}

func NewBasicAuth(username, password string) *BasicAuth {
	auth := new(BasicAuth)
	credentials := fmt.Sprintf("%s:%s", username, password)
	b64 := base64.StdEncoding.EncodeToString([]byte(credentials))
	auth.authString = fmt.Sprintf("Basic %s", b64)
	return auth
}

type TokenAuth struct {
	token string
}

func (a TokenAuth) applyAuth(r *http.Request) {
	r.Header.Set("X-Authentication-Token", a.token)
}

func NewTokenAuth(token string) *TokenAuth {
	auth := new(TokenAuth)
	auth.token = token
	return auth
}
