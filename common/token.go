package common

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
)

func GetIdFromString(id string) (pgtype.UUID, error) {
	decodedId, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("could not decode id: %w", err)
	}

	idUuid, err := uuid.Parse(string(decodedId[:]))
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("could not parse id: %w", err)
	}

	return pgtype.UUID{Bytes: idUuid, Valid: true}, nil
}

func ExtractToken(ctx *fiber.Ctx) string {
	auth := ctx.Get(fiber.HeaderAuthorization)
	if len(auth) <= 6 || !utils.EqualFold(auth[:7], "bearer ") {
		return ""
	}

	return strings.TrimSpace(auth[7:])
}

func ValidateToken(token string, publicKey string) (interface{}, interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("validation error: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}

	account, err := GetIdFromString(claims["sub"].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("could not extract account id: %w", err)
	}
	clientId, err := GetIdFromString(claims["client_id"].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("could not extract client instance id: %w", err)
	}
	return account, clientId, nil
}
