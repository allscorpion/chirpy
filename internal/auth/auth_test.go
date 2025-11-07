package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestHashing(t *testing.T) {
	testStr := "testing"
    hash, err := HashPassword(testStr);

	if err != nil {
		t.Error(err);
	}

	if hash == testStr {
		t.Errorf("invalid hasing %v == %v", testStr, hash)
	}

	found, err := CheckPasswordHash(testStr, hash);

	if err != nil {
		t.Error(err)
	}

	if !found {
		t.Errorf("hash does not match password")
	}
}

func TestJwt(t *testing.T) {
	tokenSecret := "mysecretstring";
	userId := uuid.New()
	jwt, err := MakeJWT(userId, tokenSecret);

	if err != nil {
		t.Fatal(err);
	}

	if jwt == "" {
		t.Fatal("jwt is empty")
	}
	
	jwtUserId, err := ValidateJWT(jwt, tokenSecret);

	if err != nil {
		t.Fatal(err);
	}

	if userId != jwtUserId {
		t.Fatal("userIds do not match")
	}
}

func TestGetBearerTokenValid(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer anExampleToken");
	token, err := GetBearerToken(headers);

	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("no token found")
	}
}

func TestGetBearerTokenInvalid(t *testing.T) {
	headers := http.Header{}
	token, err := GetBearerToken(headers);

	if err == nil {
		t.Fatal("no error found")
	}

	if token != "" {
		t.Fatal("token found")
	}
}