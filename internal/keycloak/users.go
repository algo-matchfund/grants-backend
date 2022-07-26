package keycloak

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

func (s *KeycloakService) GetUsers(params *GetUsersParams) ([]User, error) {
	if authenticated, err := s.authenticate(); !authenticated {
		return nil, err
	}

	req, err := s.getNewAdminRequest("GET", s.apiURL+"/users", s.accessToken.token)
	if err != nil {
		return nil, err
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	var users []User

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *KeycloakService) GetUserByID(id string) (*User, error) {
	if authenticated, err := s.authenticate(); !authenticated {
		return nil, err
	}

	req, err := s.getNewAdminRequest("GET", s.apiURL+"/users/"+url.PathEscape(id), s.accessToken.token)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	var user *User

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *KeycloakService) UpdateUser(id string, updatedUser User) error {
	if authenticated, err := s.authenticate(); !authenticated {
		return err
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(updatedUser)
	if err != nil {
		return err
	}
	req, err := s.getNewAdminUpdateRequest("PUT", s.apiURL+"/users/"+url.PathEscape(id), s.accessToken.token, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return err
}

func (s *KeycloakService) GetUserSessions(id string) ([]interface{}, error) {
	if authenticated, err := s.authenticate(); !authenticated {
		return nil, err
	}

	return nil, nil
}

func (s *KeycloakService) GetUserGroups(id string) ([]string, error) {
	if authenticated, err := s.authenticate(); !authenticated {
		return nil, err
	}

	return nil, nil
}

func (s *KeycloakService) AssignUserToGroup(id, groupId string) error {
	if authenticated, err := s.authenticate(); !authenticated {
		return err
	}

	return nil
}
