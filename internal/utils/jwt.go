package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type SupabaseClaims struct {
	Subject      string                 `json:"sub"`
	Email        string                 `json:"email"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
	jwt.RegisteredClaims
}

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

func ParseSupabaseJWT(tokenString string) (*SupabaseClaims, error) {
	publicKey, err := getECPublicKeyFromJWKS()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SupabaseClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func getECPublicKeyFromJWKS() (*ecdsa.PublicKey, error) {
	jwksURL := os.Getenv("SUPABASE_JWT_URL")
	if jwksURL == "" {
		jwksURL = "https://lopfhsbiefsuauymbksy.supabase.co/auth/v1/.well-known/jwks.json"
	}

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks JWKSResponse
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, err
	}

	if len(jwks.Keys) == 0 {
		return nil, errors.New("no keys found in JWKS")
	}

	key := jwks.Keys[0]

	xBytes, err := base64.RawURLEncoding.DecodeString(key.X)
	if err != nil {
		return nil, err
	}
	yBytes, err := base64.RawURLEncoding.DecodeString(key.Y)
	if err != nil {
		return nil, err
	}

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}, nil
}

func GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := ParseSupabaseJWT(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}