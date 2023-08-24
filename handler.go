package main

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"

	"net/http"
	"time"
)

type jwtAccessClaims struct {
	//password string `json:"password"`
	//email string `json:"email"`
	GUID string `json:"guid"`
	jwt.RegisteredClaims
}

type jwtRefreshSimpleClaims struct {
	GUID string `json:"guid"`
	jwt.RegisteredClaims
}

type handler struct{}

func (h *handler) login(c echo.Context) error {
	guid := c.QueryParam("guid")

	// Check in your db if the user exists or not, could verify password here if needed
	// Throws unauthorized error
	query := bson.D{{Key: "guid", Value: guid}}
	_, err := findOne(context.TODO(), query, "userData")
	if err != nil {
		return err
	}

	accessClaims := &jwtAccessClaims{
		guid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(loadedConfig.AccessTokenExpiresIn)),
		},
	}

	refreshClaims := &jwtRefreshSimpleClaims{
		guid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(loadedConfig.RefreshTokenExpiresIn)),
		},
	}

	newTokenPair, err := generateJWTtokenPairFromClaims(accessClaims, refreshClaims)
	if err != nil {
		return err
	}

	base64RefreshToken := newTokenPair["refresh_token"]
	hashSha256 := HashStringSha256(base64RefreshToken)
	fmt.Printf("HashSha256: %v for id %v", hashSha256, guid)

	hashedBase64RefreshToken, err := HashStringBcrypt(hashSha256)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "guid", Value: guid}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "rtokenid", Value: hashedBase64RefreshToken}}}}
	err = upsertOne(context.TODO(), "userLoggedInRefresh", filter, update)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newTokenPair)
}

func (h *handler) restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	accessClaims := user.Claims.(*jwtAccessClaims)
	guid := accessClaims.GUID
	return c.String(http.StatusOK, "Welcome Home!, your GUID: "+guid+"\n")
}
