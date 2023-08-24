package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"

	"net/http"
	"time"
)

func (h *handler) refreshToken(c echo.Context) error {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	tokenReq := tokenReqBody{}
	c.Bind(&tokenReq)

	//decode from full base64
	refreshToken, err := base64.RawURLEncoding.DecodeString(tokenReq.RefreshToken)
	if err != nil {
		return err
	}
	stringRefreshToken := string(refreshToken)
	fmt.Printf("Decoded refresh token: %v", stringRefreshToken) //DEBUG

	refToken, err := jwt.Parse(stringRefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(loadedConfig.RefreshTokenHSKey), nil
	})
	if err != nil {
		return err
	}

	if oldRefClaims, ok := refToken.Claims.(jwt.MapClaims); ok && refToken.Valid {

		// Get the user record from database or
		// run through your business logic to verify if the user can log in TODO
		guid := oldRefClaims["guid"].(string)
		base64RefreshToken := tokenReq.RefreshToken
		refTokenHashSha256 := HashStringSha256(base64RefreshToken)

		//Return on error or no record
		query := bson.D{{Key: "guid", Value: guid}}
		result, err := findOne(context.TODO(), query, "userLoggedInRefresh")
		if err != nil {
			return err
		}

		//Return on error or wrong hash
		bcryptHashFromDB := result["rtokenid"].(string)
		err = VerifyStringBcrypt(bcryptHashFromDB, refTokenHashSha256) //DEBUG
		if err != nil {
			return err
		}

		authClaims := &jwtAccessClaims{
			guid,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(loadedConfig.AccessTokenExpiresIn)),
			},
		}

		newRefreshClaims := &jwtRefreshSimpleClaims{
			guid,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(loadedConfig.RefreshTokenExpiresIn)),
			},
		}

		newTokenPair, err := generateJWTtokenPairFromClaims(authClaims, newRefreshClaims)
		if err != nil {
			return err
		}

		//Compare hashes
		newBase64RefreshToken := newTokenPair["refresh_token"]
		hashSha256 := HashStringSha256(newBase64RefreshToken)
		newHashedBase64RefreshToken, err := HashStringBcrypt(hashSha256)
		if err != nil {
			return err
		}

		newDocument := userLoggedIn{GUID: guid, RefreshTokenHash: newHashedBase64RefreshToken}
		filter := bson.D{{Key: "guid", Value: guid}}
		update := bson.D{{Key: "$set", Value: newDocument}}
		err = upsertOne(context.TODO(), "userLoggedInRefresh", filter, update)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, newTokenPair)

	}

	return echo.ErrUnauthorized
}
