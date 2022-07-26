package middlewares

import (
	"testing"

	"gopkg.in/square/go-jose.v2/jwt"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqKVYjrhuH2AS+qgUc1tvVJH0W9+VZwx01abTL8EhYlxlKyACdBweI5u75j2fc7jeob0ktEXHjEPolaJ9ZTsUnfbj/E5M/XkmdGIrd0bRp1IjqbXFDDebk+J0cVa1eJsCO4wzka7rSEFcX6LRT2JJBcUnC+8CgbViAwdweYt3vjYCHnqKBrNLI/qwBgqw14mbsP8xaHxlM/3RYzoOYgLejKz/SlKhmPsuyNk2amLX8qfzAABB3El9uJi1FdTkzAzrOBGw4lnTMCIWXHdQTfvkDEuCFRAE0GjICcCpIywme3Vm8qGGRAA2a4gMhB6DAaxkzbyLE1xgj7Zp1xJREHAqCQIDAQAB
-----END PUBLIC KEY-----`

const expiredJwt = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJzX0tUa3JpNGZrbHNtVEROVF95RE9Cc0VObkVxcFRfbzllX29OWG94UUxBIn0.eyJleHAiOjE2MTY0MDM2NjksImlhdCI6MTYxNjQwMzYwOSwianRpIjoiNWVhNzI0OTUtNjM0MC00OWVkLTllNjAtODFjMzIzMjQ4NTEwIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL2F1dGgvcmVhbG1zL2FuYWx5dGljcyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJkMGY4OWRhYS02ZmJiLTQxZGEtOTk3NS04NGExMTcwZjBmZjAiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhbmFseXRpY3MtZnJvbnRlbmQiLCJzZXNzaW9uX3N0YXRlIjoiNTMxMGM3MTctMWIyNS00YjJhLThlM2QtYTgzNzQ5OTgxNGIzIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwOi8vYW5hbHl0aWNzLWZyb250ZW5kOjkwMDAiLCJodHRwOi8vMTI3LjAuMC4xOjkwMDAiLCJodHRwOi8vbG9jYWxob3N0OjkwMDAiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50Iiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJlbWFpbCBwcm9maWxlIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJuYW1lIjoiS2V2aW4gQmFyYm91ciIsInByZWZlcnJlZF91c2VybmFtZSI6ImtldmluQHN1bW1pdHRvLmNvbSIsImdpdmVuX25hbWUiOiJLZXZpbiIsImZhbWlseV9uYW1lIjoiQmFyYm91ciIsImVtYWlsIjoia2V2aW5Ac3VtbWl0dG8uY29tIn0.j2kQueyFXvLkaHAP3ZFLhd-aeWavPmhoIphZrauV4zXb-JQ4MasIYSmNXVV5rs2ZjTjdJvfyPXkR1VoeFN5GGmytVMWcAH4y_Q72SQ_7GP3_UL5HeobtsJ2DQyL8Egvg3HHIJUtMoZoF-d9yXTBjzAUBGDpw5JFEOfNBtVjGOtdL-IvCkqXRJEBvLHhub85nH1S7dhbXfOqdTJ66zVOmMGHDnAI8W9He8-P0EVmIVDJK_Vysoip5TH1TuioToX4htbNfQhwoDAxhqnNAVIxVDijYnmyQlb1cSut6h4JC6MOmZNIhLStBKjzwffD73IcIaxucDbSO8UqBvP3wcoiFCg`
const goodJwt = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJzX0tUa3JpNGZrbHNtVEROVF95RE9Cc0VObkVxcFRfbzllX29OWG94UUxBIn0.eyJleHAiOjE3MDI4MDM1NDAsImlhdCI6MTYxNjQwMzU0MCwianRpIjoiNjFjN2UzYmYtMGUwOC00ZjU2LWE3OWUtNzY2M2MwY2IxZWM4IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL2F1dGgvcmVhbG1zL2FuYWx5dGljcyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJkMGY4OWRhYS02ZmJiLTQxZGEtOTk3NS04NGExMTcwZjBmZjAiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhbmFseXRpY3MtZnJvbnRlbmQiLCJzZXNzaW9uX3N0YXRlIjoiMmY3ZWYwMWItZDlkZC00ODkyLWEyZjEtODMyMTJkM2ZkNGRjIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwOi8vYW5hbHl0aWNzLWZyb250ZW5kOjkwMDAiLCJodHRwOi8vMTI3LjAuMC4xOjkwMDAiLCJodHRwOi8vbG9jYWxob3N0OjkwMDAiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50Iiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJlbWFpbCBwcm9maWxlIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJuYW1lIjoiS2V2aW4gQmFyYm91ciIsInByZWZlcnJlZF91c2VybmFtZSI6ImtldmluQHN1bW1pdHRvLmNvbSIsImdpdmVuX25hbWUiOiJLZXZpbiIsImZhbWlseV9uYW1lIjoiQmFyYm91ciIsImVtYWlsIjoia2V2aW5Ac3VtbWl0dG8uY29tIn0.WplhNKevCnyCemME402FxjUEpkx_Ou8-50I9MMjaT-WoveagTImrQUxR4FT3mc60sobOTZbz3zxFbRwcR-MdN0ueDJbH3ZlacEh_djRZHRXpu-Z8f4f0vPlNTydEvA-ADibXlVBh3ACAi3LX1PoGYPZgBus4eHhUjcytRL1NUCekn8vzpLi83hU-UIZUpldlmZ7kyyQWZ8P0U9PzHyIMkjcLNDOfrQTFKy3NIDXlfDrToNuvoeGuRbjOeR1PUPcK2pRR51IxjVOWgd9qxhmZ6-oHW-t9w82oaZlHOS6gUfQ9KNRd4GvoerVl5DfAryP4rteALvJEESm4PQaNl2F1Gw`
const badJwt = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJxUk5SRzV3QWktOF9KMUV2cWs0OXJaLVJ0eE5NXEtDbHhmNXVoZk1jeEIwIn0.eyJleHAiOjE2MTYxODczNTgsImlhdCI6MTYxNjE1MTM1OCwianRpIjoiNzQxOTE0MzgtOTA5OS00YTBiLThmNzctMjgwYTE2MWRhYjEzIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL2F1dGgvcmVhbG1zL2FuYWx5dGljcyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiI1YzJjYjhkZi1lMDJhLTQ4NjktYmMyYi02YTA4NTMzMmNjMTUiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhbmFseXRpY3MtZnJvbnRlbmQiLCJzZXNzaW9uX3N0YXRlIjoiMjNkMzVmYjctN2EyOC00M2Q4LTk5ZmItMDQ0YTIyMWI2ZmJlIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwOi8vYW5hbHl0aWNzLWZyb250ZW5kOjkwMDAiLCJodHRwOi8vMTI3LjAuMC4xOjkwMDAiLCJodHRwOi8vbG9jYWxob3N0OjkwMDAiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50Iiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJlbWFpbCBwcm9maWxlIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLZXZpbiBCIiwicHJlZmVycmVkX3VzZXJuYW1lIjoia2V2aW5Ac3VtbWl0dG8uY29tIiwiZ2l2ZW5fbmFtZSI6IktldmluIiwiZmFtaWx5X25hbWUiOiJCIiwiZW1haWwiOiJrZXZpbkBzdW1taXR0by5jb20ifQ.eI8QSlJSFSwAAqmzGyaeyQOh9yyqM0afYa52TlVEczhEP0qXWuBUt4TfyEU3dh_WDXsnOQOQfdWeD5GJgS7h-IgfY5RFXSKAUhFG08RVCheYXEKo-PjLJkAIigJhDtPjv2LwOdS4C7jS5V-lGXJDdie_IKAD4MWDSvJMgGQdldKDOxsn-WIHHAJkjGNdEZxXQE84FTZxM0RMuSRz3qj5wbqcoN5XpQtDj86wxfOibeWkla452LOA1tXaGPeGo237B4_WSgPK_QGqpyKG1LSS3dUTzhhpvTpk3tmfMw9kNBVZtBCAfzYG5C-CmvwmoGz9lhqQYZhSb8hHILU4UQxANg`
const badSigJwt = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`

func Test_VerifyExpiredJwt(t *testing.T) {
	authenticate, _ := NewAuthenticator(publicKey)

	_, err := authenticate(expiredJwt)

	if err != jwt.ErrExpired {
		t.Fatalf("Expected expiration error, got %s\n", err)
	}
}

func Test_VerifyGoodJwt(t *testing.T) {
	authenticate, _ := NewAuthenticator(publicKey)

	p, err := authenticate(goodJwt)

	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if p.Email != "kevin@matchfund.app" {
		t.Fatalf("Expected kevin@matchfund.app, got %s", p.Email)
	}

	if p.Name != "Kevin Barbour" {
		t.Fatalf("Expected Kevin B, got %s", p.Name)
	}
}

func Test_VerifyBadJwt(t *testing.T) {
	authenticate, _ := NewAuthenticator(publicKey)

	_, err := authenticate(badJwt)

	if err == nil {
		t.Fatal("Expected error...")
	}
}

func Test_CreateAuthBadKey(t *testing.T) {
	_, err := NewAuthenticator("abc")

	if err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func Test_BadSigJwt(t *testing.T) {
	authenticate, _ := NewAuthenticator(publicKey)

	_, err := authenticate(badSigJwt)

	if err == nil {
		t.Fatal("Expected error parsing validating jwt")
	}
}
