package main

import (
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWTtokenPairFromClaims(claimsAuthToken *jwtAccessClaims, claimsRefreshToken *jwtRefreshSimpleClaims) (map[string]string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claimsAuthToken)
	t, err := token.SignedString([]byte(loadedConfig.AccessTokenHSKey))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claimsRefreshToken)
	rt, err := refreshToken.SignedString([]byte(loadedConfig.RefreshTokenHSKey))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  t,
		"refresh_token": base64.RawURLEncoding.EncodeToString([]byte(rt)),
	}, nil

}
