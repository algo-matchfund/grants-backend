package keycloak

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/algo-matchfund/grants-backend/internal/config"
)

type AccessToken struct {
	token     string
	expiresAt time.Time
	lock      *sync.RWMutex
}

type KeycloakService struct {
	authURL      string
	apiURL       string
	realmClient  string
	clientSecret string

	accessToken AccessToken
	client      *http.Client
}

func NewKeycloakService(conf *config.Config) *KeycloakService {
	return &KeycloakService{
		authURL:      fmt.Sprintf("%s:%d/auth/realms/%s/protocol/openid-connect/token", conf.Authentication.Host, conf.Authentication.Port, conf.Authentication.Realm),
		apiURL:       fmt.Sprintf("%s:%d/auth/admin/realms/%s", conf.Authentication.Host, conf.Authentication.Port, conf.Authentication.Realm),
		realmClient:  conf.Authentication.Client,
		clientSecret: conf.Authentication.Secret,
		accessToken: AccessToken{
			token:     "",
			expiresAt: time.Now(),
			lock:      &sync.RWMutex{},
		},
		client: &http.Client{},
	}
}

func (t *AccessToken) String() string {
	return t.token
}

func (s *KeycloakService) authenticate() (bool, error) {
	s.accessToken.lock.RLock()
	// do not update access token if we have a token, which will not expire for the next 5 seconds
	if s.accessToken.token != "" && time.Now().Before(s.accessToken.expiresAt) && time.Until(s.accessToken.expiresAt).Seconds() > 5 {
		s.accessToken.lock.RUnlock()
		return true, nil
	}
	s.accessToken.lock.RUnlock()

	s.accessToken.lock.Lock()
	defer s.accessToken.lock.Unlock()

	form := url.Values{
		"grant_type":    []string{"client_credentials"},
		"scope":         []string{"openid"},
		"client_id":     []string{s.realmClient},
		"client_secret": []string{s.clientSecret},
	}
	request, err := http.NewRequest("POST", s.authURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Accept", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return false, err
	}

	var resp TokenResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		return false, nil
	}

	if resp.TokenType != "Bearer" || resp.AccessToken == "" {
		return false, errors.New("unexpected Keycloak access token or token type")
	}

	s.accessToken.token = resp.AccessToken
	s.accessToken.expiresAt = time.Now().Add(time.Second * time.Duration(resp.ExpiresIn))

	return true, nil
}

func (s *KeycloakService) getNewAdminRequest(method, url, token string) (*http.Request, error) {
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Pragma", "no-cache")

	return r, nil
}

func (s *KeycloakService) getNewAdminUpdateRequest(method, url, token string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Pragma", "no-cache")

	return r, nil
}
