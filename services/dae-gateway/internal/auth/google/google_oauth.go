package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type GoogleOAuth struct {
	conf     oauth2.Config
	verifier *oidc.IDTokenVerifier
}

type startState struct {
	State string `json:"state"`
	Nonce string `json:"nonce"`
	Exp   int64  `json:"timestamp"`
}

func NewGoogleOIDC(ctx context.Context, cfg GoogleOAuthConfig) (*GoogleOAuth, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verify := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})
	return &GoogleOAuth{
		conf:     oauth2Config,
		verifier: verify,
	}, nil
}

func (g *GoogleOAuth) Start(w http.ResponseWriter, r *http.Request) {
	state := randB64(32)
	nonce := randB64(32)

	ss := startState{
		State: state,
		Nonce: nonce,
		Exp:   time.Now().Add(10 * time.Minute).Unix(),
	}
	setJSONCookie(w, "g_oidc", ss)

	url := g.conf.AuthCodeURL(state, oidc.Nonce(nonce))

	http.Redirect(w, r, url, http.StatusFound)
}

func randB64(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func setJSONCookie(w http.ResponseWriter, name string, value any) {
	raw, _ := json.Marshal(value)
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    base64.RawURLEncoding.EncodeToString(raw),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600,
	})

}

func readJSONCookie(r *http.Request, name string, out any) error {
	cookie, err := r.Cookie(name)
	if err != nil {
		return err
	}
	raw, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, out)
}

func (g *GoogleOAuth) Callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	state := q.Get("state")
	if code == "" || state == "" {
		http.Error(w, "missing code or state", http.StatusBadRequest)
		return
	}

	var ss startState
	if err := readJSONCookie(r, "g_oidc", &ss); err != nil {
		http.Error(w, "invalid state cookie", http.StatusBadRequest)
		return
	}
	if time.Now().Unix() > ss.Exp || ss.State != state {
		http.Error(w, "invalid or expired state", http.StatusBadRequest)
		return
	}

	tokens, err := g.conf.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawID, ok := tokens.Extra("id_token").(string)
	if !ok || rawID == "" {
		http.Error(w, "no id_token in token response", http.StatusInternalServerError)
		return
	}

	idToken, err := g.verifier.Verify(r.Context(), rawID)
	if err != nil {
		http.Error(w, "failed to verify id_token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if idToken.Nonce != ss.Nonce {
		http.Error(w, "invalid nonce in id_token", http.StatusBadRequest)
		return
	}

	var claims struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "failed to parse id_token claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(claims)
}
