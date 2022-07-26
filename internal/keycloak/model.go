package keycloak

import (
	"net/url"
	"strconv"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint32 `json:"expires_in"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
	Scope       string `json:"scope"`
}

type User struct {
	CreatedAt     time.Time           `json:"createdTimestamp"`
	ID            string              `json:"id"`
	Username      string              `json:"username"`
	Email         string              `json:"email"`
	FirstName     string              `json:"firstName"`
	LastName      string              `json:"lastName"`
	Attributes    map[string][]string `json:"attributes,omitempty"`
	Enabled       bool                `json:"enabled"`
	EmailVerified bool                `json:"emailVerified"`
}

type GetUsersParams struct {
	// Boolean which defines whether brief representations are returned (default: false)
	BriefRepresentation bool
	// A String contained in email, or the complete email, if param "exact" is true
	Email string
	// whether the email has been verified
	EmailVerified bool
	// Boolean representing if user is enabled or not
	Enabled bool
	// Boolean which defines whether the params "last", "first", "email" and "username" must match exactly
	Exact bool
	// Pagination offset
	First int32
	// A String contained in firstName, or the complete firstName, if param "exact" is true
	FirstName string
	// The alias of an Identity Provider linked to the user
	IdpAlias string
	// The userId at an Identity Provider linked to the user
	IdpUserId string
	// A String contained in lastName, or the complete lastName, if param "exact" is true
	LastName string
	// Maximum results size (defaults to 100)
	Max int32
	// A String contained in username, first or last name, or email
	Search string
	// A String contained in username, or the complete username, if param "exact" is true}
	Username string
}

func (p *GetUsersParams) Encode() string {
	var query url.Values
	query.Set("briefRepresentation", strconv.FormatBool(p.BriefRepresentation))
	query.Set("email", p.Email)
	query.Set("emailVerified", strconv.FormatBool(p.EmailVerified))
	query.Set("enabled", strconv.FormatBool(p.Enabled))
	query.Set("exact", strconv.FormatBool(p.Exact))
	query.Set("first", strconv.FormatInt(int64(p.First), 10))
	query.Set("firstName", p.FirstName)
	query.Set("idpAlias", p.IdpAlias)
	query.Set("idpUserId", p.IdpUserId)
	query.Set("lastName", p.LastName)
	query.Set("max", strconv.FormatInt(int64(p.Max), 10))
	query.Set("search", p.Search)
	query.Set("username", p.Username)

	return query.Encode()
}
